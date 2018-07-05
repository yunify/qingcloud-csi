package block

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/kubernetes/pkg/util/mount"
	"os"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
}

// This operation MUST be idempotent
// If the volume corresponding to the volume id has already been published at the specified target path,
// and is compatible with the specified volume capability and readonly flag, the plugin MUST reply 0 OK.
// csi.NodePublishVolumeRequest:	volume id			+ Required
//									target path			+ Required
//									volume capability	+ Required
//									read only			+ Required (This field is NOT provided when requesting in Kubernetes)
func (ns *nodeServer) NodePublishVolume(
	ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	glog.Info("----- Start NodePublishVolume -----")
	defer glog.Info("===== End NodePublishVolume =====")
	// 0. Preflight
	// check volume id
	if len(req.GetVolumeId()) == 0{
		return nil, status.Error(codes.InvalidArgument, "Volume id missing in request")
	}
	// check target path
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	// Check volume capability
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities missing in request")
	} else if !HasSameVolumeAccessMode(ns.Driver.GetVolumeCapabilityAccessModes(),
		[]*csi.VolumeCapability{req.GetVolumeCapability()}) {
		return nil, status.Error(codes.FailedPrecondition, "Exceed capabilities")
	}
	// check stage path
	if len(req.GetStagingTargetPath()) == 0{
		return nil, status.Error(codes.FailedPrecondition, "Staging target path not set")
	}
	// set parameter
	targetPath := req.GetTargetPath()
	stagePath := req.GetStagingTargetPath()
	volumeId := req.GetVolumeId()

	// Create VolumeManager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// Check volume exist
	volInfo, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return nil, status.Errorf(codes.NotFound, "Volume %s does not exist", volumeId)
	}

	// 1. Mount
	// Make dir if dir not presents
	_, err = os.Stat(targetPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(targetPath, 0750); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// check targetPath is mounted
	mounter := mount.New("")
	notMnt, err := mounter.IsNotMountPoint(targetPath)
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
	// For idempotent:
	// If the volume corresponding to the volume id has already been published at the specified target path,
	// and is compatible with the specified volume capability and readonly flag, the plugin MUST reply 0 OK.
	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// set bind mount options
	options := []string{"bind"}
	if req.GetReadonly() == true {
		options = append(options, "ro")
	}
	glog.Infof("Bind mount %s at %s", stagePath, targetPath)
	if err := mounter.Mount(stagePath, targetPath, "", options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("Mount bind %s at %s succeed", stagePath, targetPath)
	return &csi.NodePublishVolumeResponse{}, nil
}

// csi.NodeUnpublishVolumeRequest:	volume id	+ Required
//									target path	+ Required
func (ns *nodeServer) NodeUnpublishVolume(
	ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	glog.Info("----- Start NodeUnpublishVolume -----")
	defer glog.Info("===== End NodeUnpublishVolume =====")
	// 0. Preflight
	// check arguments
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume id missing in request")
	}
	// set parameter
	volumeId := req.GetVolumeId()
	targetPath := req.GetTargetPath()

	// Create VolumeManager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// Check volume exist
	volInfo, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return nil, status.Errorf(codes.NotFound, "Volume %s does not exist", volumeId)
	}

	// 1. Unmount
	// check targetPath is mounted
	mounter := mount.New("")
	notMnt, err := mounter.IsNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if notMnt {
		glog.Warningf("Volume %s has not mount point", volumeId)
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}
	// do unmount
	glog.Infof("Unbind mountvolume %s/%s", targetPath, volumeId)
	if err = mounter.Unmount(targetPath); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("Unbound mount volume succeed")

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.NodeStageVolumeRequest: 	volume id			+ Required
//								stage target path	+ Required
//								volume capability	+ Required
func (ns *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	glog.Info("----- Start NodeStageVolume -----")
	defer glog.Info("===== End NodeStageVolume =====")
	capRsp , _ := ns.NodeGetCapabilities(context.Background(), nil)
	if flag := HasNodeServiceCapability(capRsp.GetCapabilities(), csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME);flag == false{
		glog.Errorf("driver capability %v", capRsp.GetCapabilities())
		return nil, status.Error(codes.Unimplemented, "Node has not stage capability")
	}
	// 0. Preflight
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	if req.GetVolumeCapability() == nil{
		return nil, status.Error(codes.InvalidArgument, "Volume capability missing in request")
	}
	// set parameter
	volumeId := req.GetVolumeId()
	targetPath := req.GetStagingTargetPath()
	fsType := req.GetVolumeCapability().GetMount().GetFsType()

	// Create VolumeManager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// Check volume exist
	volInfo, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return nil, status.Errorf(codes.NotFound, "Volume %s does not exist", volumeId)
	}
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
	// already mount
	if !notMnt {
		return &csi.NodeStageVolumeResponse{}, nil
	}

	// get device path
	devicePath := ""
	if volInfo.Instance != nil && volInfo.Instance.Device != nil && *volInfo.Instance.Device != "" {
		devicePath = *volInfo.Instance.Device
		glog.Infof("Find volume %s's device path is %s", volumeId, devicePath)
	} else {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Cannot find device path of volume %s", volumeId))
	}
	// do mount
	glog.Infof("Mounting %s to %s format %s...", volumeId, targetPath, fsType)
	diskMounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOsExec()}
	if err := diskMounter.FormatAndMount(devicePath, targetPath, fsType, []string{}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("Mount %s to %s succeed", volumeId, targetPath)
	return &csi.NodeStageVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.NodeUnstageVolumeRequest:	volume id	+ Required
//									target path	+ Required
func (ns *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	glog.Info("----- Start NodeUnstageVolume -----")
	defer glog.Info("===== End NodeUnstageVolume =====")
	capRsp , _ := ns.NodeGetCapabilities(context.Background(), nil)
	if flag := HasNodeServiceCapability(capRsp.GetCapabilities(), csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME);flag == false{
		glog.Errorf("driver capability %v", capRsp.GetCapabilities())
		return nil, status.Error(codes.Unimplemented, "Node has not unstage capability")
	}
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
	targetPath := req.GetStagingTargetPath()

	// Create VolumeManager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Check volume exist
	volInfo, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return nil, status.Errorf(codes.NotFound, "Volume %s does not exist", volumeId)
	}

	// 1. Unmount
	// check targetPath is mounted
	// For idempotent:
	// If the volume corresponding to the volume id is not staged to the staging target path,
	// the plugin MUST reply 0 OK.
	mounter := mount.New("")
	notMnt, err := mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if notMnt {
		return &csi.NodeUnstageVolumeResponse{}, nil
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
	glog.Infof("block image: volume %s has been unmounted.", volumeId)
	cnt--
	glog.Infof("block image: mount count: %d", cnt)
	if cnt > 0 {
		glog.Errorf("image %s still mounted in instance %s", volumeId, GetCurrentInstanceId())
		return nil, status.Error(codes.Internal, "unmount failed")
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	glog.Info("----- Start NodeGetCapabilities -----")
	defer glog.Info("===== End NodeGetCapabilities =====")
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

func (ns *nodeServer) NodeGetId(ctx context.Context, req *csi.NodeGetIdRequest) (*csi.NodeGetIdResponse, error) {
	glog.Info("----- Start NodeGetId -----")
	defer glog.Info("===== End NodeGetId =====")

	return &csi.NodeGetIdResponse{
		NodeId: GetCurrentInstanceId(),
	}, nil
}