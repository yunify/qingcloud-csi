package block

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/kubernetes/pkg/util/mount"
	"os"
	"path"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
}

func (ns *nodeServer) NodePublishVolume(
	ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	glog.Infof("NodePublishVolume")
	// 0. Preflight
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	// set parameter
	volumeId := req.GetVolumeId()
	targetPath := req.GetTargetPath()
	fsType := req.GetVolumeCapability().GetMount().GetFsType()
	stagePath := req.GetStagingTargetPath()
	// 1. Mount
	// if volume already mounted
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(targetPath, 0750); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}
	// do mount
	glog.Infof("Mounting %s/%s format %s...", req.GetStagingTargetPath(),req.GetVolumeId(), fsType)
	diskMounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOsExec()}
	if err := diskMounter.FormatAndMount(stagePath, targetPath, fsType, []string{}); err != nil {
		return nil, err
	}
	glog.Infof("Mount %s/%s succeed", targetPath, volumeId)

	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(
	ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	glog.Infof("NodeUnpublishVolume")
	// 0. Preflight
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	// set parameter
	volumeID := req.GetVolumeId()
	targetPath := req.GetTargetPath()

	// 1. Unmount
	// check targetPath is mounted
	mounter:= mount.New("")
	notMnt, err := mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if notMnt {
		return &csi.NodeUnpublishVolumeResponse{},nil
	}
	// do unmount
	err = mounter.Unmount(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("block image: volume %s/%s has been unmounted.",  targetPath,volumeID)

	return &csi.NodeUnpublishVolumeResponse{},nil
}


func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	glog.Infof("NodeStageVolume")
	// 0. Preflight
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	// set parameter
	volumeId := req.GetVolumeId()
	instanceId := GetCurrentInstanceId()
	targetPath := req.GetStagingTargetPath()
	fsType := req.GetVolumeCapability().GetMount().GetFsType()

	// 1. Attach
	// create StorageClass
	sc, err := NewStorageClassFromMap(req.VolumeAttributes)
	if err != nil {
		glog.Error(err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// create volume provisioner object
	vp, err := newVolumeProvisioner(sc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// attach volume
	glog.Infof("Attaching volume %s to instance %s in zone %s...", volumeId, instanceId, sc.Zone)
	devicePath, err := vp.AttachVolume(volumeId, instanceId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("Attaching volume %s succeed.", volumeId, instanceId, sc.Zone )
	// save volInfo into a file.
	glog.Infof("Save volume %s info to a file...", volumeId)
	blockVol := blockVolume{}
	blockVol.Zone = sc.Zone
	blockVol.Sc = *sc
	if err := persistVolInfo(volumeId, path.Join(PluginFolder, "node"), &blockVol); err != nil {
		glog.Warningf("Failed to store volInfo with error: %v", err)
	}
	glog.Infof("Save volume %s info succeed", volumeId)

	// 2. Mount
	// if volume already mounted
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(targetPath, 0750); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if !notMnt {
		return &csi.NodeStageVolumeResponse{}, nil
	}
	// do mount
	glog.Infof("Mounting %s/%s format %s...", req.GetStagingTargetPath(),req.GetVolumeId(), fsType)
	diskMounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOsExec()}
	if err := diskMounter.FormatAndMount(devicePath, targetPath, fsType, []string{}); err != nil {
		return nil, err
	}
	glog.Infof("Mount %s/%s succeed", targetPath, volumeId)
	return &csi.NodeStageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	glog.Infof("NodeUnstageVolume")
	// 0. Preflight
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	// set parameter
	volumeID := req.GetVolumeId()
	targetPath := req.GetStagingTargetPath()

	// 1. Unmount
	// check targetPath is mounted
	mounter:= mount.New("")
	notMnt, err := mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if notMnt {
		return &csi.NodeUnstageVolumeResponse{},nil
	}
	// count mount point
	_, cnt, err := mount.GetDeviceNameFromMount(mounter, targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// do unmount
	err = mounter.Unmount(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("block image: volume %s/%s has been unmounted.",  targetPath,volumeID)
	cnt--
	glog.Infof("block image: mount count: %d", cnt)
	if cnt > 0{
		glog.Errorf("image %s still mounted in instance %s", volumeID, GetCurrentInstanceId())
		return nil, status.Error(codes.Internal, "unmount failed")
	}

	// 2. Detach
	// retrieve sc from file
	blockVol := blockVolume{}
	if err := loadVolInfo(volumeID, path.Join(PluginFolder, "node"), &blockVol); err != nil {
		return nil, err
	}
	// create volume provisioner object
	vp, err := newVolumeProvisioner(&blockVol.Sc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// do detach
	err = vp.DetachVolume(volumeID, GetCurrentInstanceId())
	if err != nil{
		glog.Errorf("failed to detach block image: %s from instance %s with error: %v",
			volumeID, GetCurrentInstanceId(), err)
		return nil, err
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
		},
	}, nil
}
