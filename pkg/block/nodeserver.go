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
	return nil, nil
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
