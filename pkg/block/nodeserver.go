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
	"strings"
	"path"
)

type nodeServer struct {
	*csicommon.DefaultNodeServer
}

func (ns *nodeServer) NodePublishVolume(
	ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	glog.Infof("NodePublishVolume")
	// Resolve request parameters
	targetPath := req.GetTargetPath()
	if !strings.HasSuffix(targetPath, "/mount") {
		return nil, fmt.Errorf("malformed the value of target path: %s", targetPath)
	}

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
	sc, err := NewStorageClassFromMap(req.VolumeAttributes)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// Create volume provisioner object
	vp, err := newVolumeProvisioner(sc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Attach volume
	devicePath, err := vp.AttachVolume(req.GetVolumeId(), GetCurrentInstanceId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// Store volInfo into a persistent file.
	blockVol := blockVolume{}
	blockVol.Zone = sc.Zone
	blockVol.Sc = *sc
	if err := persistVolInfo(req.GetVolumeId(), path.Join(PluginFolder, "node"), &blockVol); err != nil {
		glog.Warningf("failed to store volInfo with error: %v", err)
	}
	// Print log
	glog.Infof("block image: %s in %s was successfully attached at instance %s\n",
		req.GetVolumeId(), sc.Zone, GetCurrentInstanceId())

	// Mount volume
	// Get fsType
	fsType := req.GetVolumeCapability().GetMount().GetFsType()
	// Get readOnly
	options := []string{}
	readOnly := req.GetReadonly()
	if readOnly {
		options = append(options, "ro")
	}
	// Mount
	diskMounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOsExec()}
	if err := diskMounter.FormatAndMount(devicePath, targetPath, fsType, options); err != nil {
		return nil, err
	}
	glog.Infof("block image: %s in %s was successfully mounted at instance %s\n",
		req.GetVolumeId(), sc.Zone, GetCurrentInstanceId())
	return &csi.NodePublishVolumeResponse{}, nil
}

func (ns *nodeServer) NodeUnpublishVolume(
	ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	glog.Infof("NodeUnpublishVolume")

	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	// Get parameter
	volumeID := req.GetVolumeId()
	targetPath := req.GetTargetPath()

	// Check targetPath is mounted
	mounter:= mount.New("")
	notMnt, err := mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if notMnt {
		return nil, status.Error(codes.NotFound, "Volume not mounted")
	}

	_, cnt, err := mount.GetDeviceNameFromMount(mounter, targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Unmount the image
	err = mounter.Unmount(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("block image: volume %s/%s has been unmounted.",  targetPath,volumeID)
	cnt--
	glog.Infof("block image: mount count: %d", cnt)
	if cnt != 0{
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}

	// Detach block image
	// Retrieve sc from file
	blockVol := blockVolume{}
	if err := loadVolInfo(volumeID, path.Join(PluginFolder, "node"), &blockVol); err != nil {
		return nil, err
	}
	// Create volume provisioner object
	vp, err := newVolumeProvisioner(&blockVol.Sc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Detach
	err = vp.DetachVolume(volumeID, GetCurrentInstanceId())
	if err != nil{
		glog.Errorf("failed to detach block image: %s from instance %s with error: %v",
			volumeID, GetCurrentInstanceId(), err)
		return nil, err
	}

	glog.Infof("success to detach block image: %s from instance %s", volumeID, GetCurrentInstanceId())
	return &csi.NodeUnpublishVolumeResponse{},nil
}


func (ns *nodeServer) NodeStageVolume(
	ctx context.Context,
	req *csi.NodeStageVolumeRequest) (
	*csi.NodeStageVolumeResponse, error) {
	glog.Infof("NodeStageVolume")
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *nodeServer) NodeUnstageVolume(
	ctx context.Context,
	req *csi.NodeUnstageVolumeRequest) (
	*csi.NodeUnstageVolumeResponse, error) {
	glog.Infof("NodeUnstageVolume")
	return nil, status.Error(codes.Unimplemented, "")
}
