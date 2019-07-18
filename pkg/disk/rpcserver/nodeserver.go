// +-------------------------------------------------------------------------
// | Copyright (C) 2018 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package rpcserver

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/cloudprovider"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/util/resizefs"
	"os"
	"strconv"
	"strings"
)

type DiskNodeServer struct {
	driver  *driver.DiskDriver
	cloud   cloudprovider.CloudManager
	mounter *mount.SafeFormatAndMount
}

// NewNodeServer
// Create node server
func NewNodeServer(d *driver.DiskDriver, c cloudprovider.CloudManager, mnt *mount.SafeFormatAndMount) *DiskNodeServer {
	return &DiskNodeServer{
		driver:  d,
		cloud:   c,
		mounter: mnt,
	}
}

// This operation MUST be idempotent
// If the volume corresponding to the volume id has already been published at the specified target path,
// and is compatible with the specified volume capability and readonly flag, the plugin MUST reply 0 OK.
// csi.NodePublishVolumeRequest:	volume id			+ Required
//									target path			+ Required
//									volume capability	+ Required
//									read only			+ Required (This field is NOT provided when requesting in Kubernetes)
func (ns *DiskNodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.
	NodePublishVolumeResponse, error) {
	klog.Info("----- Start NodePublishVolume -----")
	defer klog.Info("===== End NodePublishVolume =====")
	// 0. Preflight
	// check volume id
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume id missing in request")
	}
	// check target path
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Target path missing in request")
	}
	// Check volume capability
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities missing in request")
	} else if !ns.driver.ValidateVolumeCapability(req.GetVolumeCapability()) {
		return nil, status.Error(codes.FailedPrecondition, "Exceed capabilities")
	}
	// check stage path
	if len(req.GetStagingTargetPath()) == 0 {
		return nil, status.Error(codes.FailedPrecondition, "Staging target path not set")
	}
	// set parameter
	targetPath := req.GetTargetPath()
	stagePath := req.GetStagingTargetPath()
	volumeId := req.GetVolumeId()

	// set fsType
	qc, err := driver.NewQingStorageClassFromMap(req.GetVolumeContext())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fsType := qc.FsType

	// Check volume exist
	volInfo, err := ns.cloud.FindVolume(volumeId)
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
	klog.Infof("Bind mount %s at %s, fsType %s, options %v ...", stagePath, targetPath, fsType, options)
	if err := mounter.Mount(stagePath, targetPath, fsType, options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("Mount bind %s at %s succeed", stagePath, targetPath)
	return &csi.NodePublishVolumeResponse{}, nil
}

// csi.NodeUnpublishVolumeRequest:	volume id	+ Required
//									target path	+ Required
func (ns *DiskNodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.
	NodeUnpublishVolumeResponse, error) {
	klog.Info("----- Start NodeUnpublishVolume -----")
	defer klog.Info("===== End NodeUnpublishVolume =====")
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

	// Check volume exist
	volInfo, err := ns.cloud.FindVolume(volumeId)
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
		klog.Warningf("Volume %s has not mount point", volumeId)
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}
	// do unmount
	klog.Infof("Unbind mountvolume %s/%s", targetPath, volumeId)
	if err = mounter.Unmount(targetPath); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("Unbound mount volume succeed")

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.NodeStageVolumeRequest: 	volume id			+ Required
//								stage target path	+ Required
//								volume capability	+ Required
func (ns *DiskNodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse,
	error) {
	klog.Info("----- Start NodeStageVolume -----")
	defer klog.Info("===== End NodeStageVolume =====")
	if flag := ns.driver.ValidateNodeServiceRequest(csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME); flag == false {
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
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capability missing in request")
	}
	// set parameter
	volumeId := req.GetVolumeId()
	targetPath := req.GetStagingTargetPath()
	// set fsType
	qc, err := driver.NewQingStorageClassFromMap(req.GetPublishContext())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fsType := qc.FsType

	// Check volume exist
	volInfo, err := ns.cloud.FindVolume(volumeId)
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
		klog.Infof("Find volume %s's device path is %s", volumeId, devicePath)
	} else {
		return nil, status.Errorf(codes.Internal, "Cannot find device path of volume %s", volumeId)
	}
	// do mount
	klog.Infof("Mounting %s to %s format %s...", volumeId, targetPath, fsType)
	diskMounter := &mount.SafeFormatAndMount{Interface: mount.New(""), Exec: mount.NewOsExec()}
	if err := diskMounter.FormatAndMount(devicePath, targetPath, fsType, []string{}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("Mount %s to %s succeed", volumeId, targetPath)
	return &csi.NodeStageVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.NodeUnstageVolumeRequest:	volume id	+ Required
//									target path	+ Required
func (ns *DiskNodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.
	NodeUnstageVolumeResponse, error) {
	klog.Info("----- Start NodeUnstageVolume -----")
	defer klog.Info("===== End NodeUnstageVolume =====")
	if flag := ns.driver.ValidateNodeServiceRequest(csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME); flag == false {
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

	// Check volume exist
	volInfo, err := ns.cloud.FindVolume(volumeId)
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
	klog.Infof("disk volume %s has been unmounted.", volumeId)
	cnt--
	klog.Infof("disk volume mount count: %d", cnt)
	if cnt > 0 {
		klog.Errorf("image %s still mounted in instance %s", volumeId, ns.driver.GetInstanceId())
		return nil, status.Error(codes.Internal, "unmount failed")
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *DiskNodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.
	NodeGetCapabilitiesResponse, error) {
	klog.Info("----- Start NodeGetCapabilities -----")
	defer klog.Info("===== End NodeGetCapabilities =====")
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: ns.driver.GetNodeCapability(),
	}, nil
}

func (ns *DiskNodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.V(2).Info("----- Start NodeGetInfo -----")
	defer klog.Info("===== End NodeGetInfo =====")

	return &csi.NodeGetInfoResponse{
		NodeId:            ns.driver.GetInstanceId(),
		MaxVolumesPerNode: ns.driver.GetMaxVolumePerNode(),
	}, nil
}

// NodeExpandVolume will expand filesystem of volume.
// Input Parameters:
//  volume id: REQUIRED
//  volume path: REQUIRED
func (ns *DiskNodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (
	*csi.NodeExpandVolumeResponse, error) {
	defer common.EntryFunction("NodeExpandVolume")()
	// 0. Preflight
	// check arguments
	klog.Info("Check input arguments")
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetVolumePath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume path missing in request")
	}
	requestSizeBytes, err := common.GetRequestSizeBytes(req.GetCapacityRange())
	if err != nil {
		return nil, status.Error(codes.OutOfRange, err.Error())
	}
	// Set parameter
	volumeId := req.GetVolumeId()
	volumePath := req.GetVolumePath()

	// Check volume exist
	klog.Infof("Get volume %s info", volumeId)
	volInfo, err := ns.cloud.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return nil, status.Errorf(codes.NotFound, "Volume %s does not exist", volumeId)
	}
	// get device path
	devicePath := ""
	if volInfo.Instance != nil && volInfo.Instance.Device != nil && *volInfo.Instance.Device != "" {
		devicePath = *volInfo.Instance.Device
		klog.Infof("Find volume %s's device path is %s", volumeId, devicePath)
	} else {
		return nil, status.Errorf(codes.Internal, "Cannot find device path of volume %s", volumeId)
	}

	resizer := resizefs.NewResizeFs(ns.mounter)
	klog.Infof("Resize file system device %s, mount path %s ...", devicePath, volumePath)
	ok, err := resizer.Resize(devicePath, volumePath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if ok != true {
		return nil, status.Error(codes.Internal, "failed to expand volume filesystem")
	}
	klog.Info("Succeed to resize file system")

	//  Check the block size
	blkSizeBytes, err := ns.getBlockSizeBytes(devicePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"expand volume error when getting size of block volume at path %s: %v", devicePath, err)
	}
	klog.Infof("Block size %d Byte, request size %d Byte", blkSizeBytes, requestSizeBytes)

	if blkSizeBytes < requestSizeBytes {
		// It's possible that the somewhere the volume size was rounded up, getting more size than requested is a success
		return nil, status.Errorf(codes.Internal, "resize requested for %v but after resize volume was size %v",
			requestSizeBytes, blkSizeBytes)
	}
	return &csi.NodeExpandVolumeResponse{
		CapacityBytes: blkSizeBytes,
	}, nil
}

func (ns *DiskNodeServer) NodeGetVolumeStats(ctx context.Context,
	req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (ns *DiskNodeServer) getBlockSizeBytes(devicePath string) (int64, error) {
	output, err := ns.mounter.Exec.Run("blockdev", "--getsize64", devicePath)
	if err != nil {
		return -1, fmt.Errorf("error when getting size of block volume at path %s: output: %s, err: %v", devicePath, string(output), err)
	}
	strOut := strings.TrimSpace(string(output))
	gotSizeBytes, err := strconv.ParseInt(strOut, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse size %s into int a size", strOut)
	}
	return gotSizeBytes, nil
}
