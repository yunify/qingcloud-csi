/*
Copyright (C) 2018 Yunify, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this work except in compliance with the License.
You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rpcserver

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/cloud"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/util/resizefs"
	"k8s.io/kubernetes/pkg/volume"
	"os"
	"strconv"
	"strings"
)

type NodeServer struct {
	driver  *driver.DiskDriver
	cloud   cloud.CloudManager
	mounter *mount.SafeFormatAndMount
	locks   *common.ResourceLocks
}

var _ csi.NodeServer = &NodeServer{}

// NewNodeServer
// Create node server
func NewNodeServer(d *driver.DiskDriver, c cloud.CloudManager, mnt *mount.SafeFormatAndMount) *NodeServer {
	return &NodeServer{
		driver:  d,
		cloud:   c,
		mounter: mnt,
		locks:   common.NewResourceLocks(),
	}
}

// This operation MUST be idempotent
// If the volume corresponding to the volume id has already been published at the specified target path,
// and is compatible with the specified volume capability and readonly flag, the plugin MUST reply 0 OK.
// csi.NodePublishVolumeRequest:	volume id			+ Required
//									target path			+ Required
//									volume capability	+ Required
//									read only			+ Required (This field is NOT provided when requesting in Kubernetes)
func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.
	NodePublishVolumeResponse, error) {
	funcName := "NodePublishVolume"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
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

	// ensure one call in-flight
	klog.Infof("Try to lock resource %s", volumeId)
	if acquired := ns.locks.TryAcquire(volumeId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volumeId)
	}
	defer ns.locks.Release(volumeId)

	// set fsType
	qc, err := driver.NewQingStorageClassFromMap(req.GetVolumeContext())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fsType := qc.GetFsType()

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
	notMnt, err := ns.mounter.IsNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			klog.Infof("Cannot find target path %s and create it.", targetPath)
			if err = os.MkdirAll(targetPath, 0750); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			klog.Infof("Succeed to create target path %s", targetPath)
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
	if err := ns.mounter.Mount(stagePath, targetPath, fsType, options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("Mount bind %s at %s succeed", stagePath, targetPath)
	return &csi.NodePublishVolumeResponse{}, nil
}

// csi.NodeUnpublishVolumeRequest:	volume id	+ Required
//									target path	+ Required
func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.
	NodeUnpublishVolumeResponse, error) {
	funcName := "NodeUnpublishVolume"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
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
	// ensure one call in-flight
	klog.Infof("Try to lock resource %s", volumeId)
	if acquired := ns.locks.TryAcquire(volumeId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volumeId)
	}
	defer ns.locks.Release(volumeId)
	// Check volume exist
	volInfo, err := ns.cloud.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return nil, status.Errorf(codes.NotFound, "Volume %s does not exist", volumeId)
	}

	// 1. Unmount
	err = mount.CleanupMountPoint(targetPath, ns.mounter.Interface, true)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unmount target path %s error: %v", targetPath, err)
	}
	klog.Infof("Unbound mount volume succeed")

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.NodeStageVolumeRequest: 	volume id			+ Required
//								stage target path	+ Required
//								volume capability	+ Required
func (ns *NodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse,
	error) {
	funcName := "NodeStageVolume"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
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
	// ensure one call in-flight
	klog.Infof("Try to lock resource %s", volumeId)
	if acquired := ns.locks.TryAcquire(volumeId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volumeId)
	}
	defer ns.locks.Release(volumeId)
	// set fsType
	qc, err := driver.NewQingStorageClassFromMap(req.GetPublishContext())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	fsType := qc.GetFsType()

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
	if err := ns.mounter.FormatAndMount(devicePath, targetPath, fsType, []string{}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("Mount %s to %s succeed", volumeId, targetPath)
	return &csi.NodeStageVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.NodeUnstageVolumeRequest:	volume id	+ Required
//									target path	+ Required
func (ns *NodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.
	NodeUnstageVolumeResponse, error) {
	funcName := "NodeUnstageVolume"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
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
	// ensure one call in-flight
	klog.Infof("Try to lock resource %s", volumeId)
	if acquired := ns.locks.TryAcquire(volumeId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volumeId)
	}
	defer ns.locks.Release(volumeId)
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
	notMnt, err := ns.mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if notMnt {
		return &csi.NodeUnstageVolumeResponse{}, nil
	}
	// count mount point
	_, cnt, err := mount.GetDeviceNameFromMount(ns.mounter, targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// do unmount
	err = ns.mounter.Unmount(targetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	klog.Infof("Disk volume %s has been unmounted.", volumeId)
	cnt--
	klog.Infof("Disk volume mount count: %d", cnt)
	if cnt > 0 {
		klog.Errorf("Volume %s still mounted in instance %s", volumeId, ns.driver.GetInstanceId())
		return nil, status.Error(codes.Internal, "unmount failed")
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}

func (ns *NodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.
	NodeGetCapabilitiesResponse, error) {
	funcName := "NodeGetCapabilities"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: ns.driver.GetNodeCapability(),
	}, nil
}

func (ns *NodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	funcName := "NodeGetInfo"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	instInfo, err := ns.cloud.FindInstance(ns.driver.GetInstanceId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if instInfo == nil {
		return nil, status.Errorf(codes.NotFound, "cannot found instance %s", ns.driver.GetInstanceId())
	}

	instanceType, ok := driver.InstanceTypeName[driver.InstanceType(*instInfo.InstanceClass)]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported instance type %d", *instInfo.InstanceClass)
	}
	top := &csi.Topology{
		Segments: map[string]string{
			ns.driver.GetTopologyInstanceTypeKey(): instanceType,
			ns.driver.GetTopologyZoneKey():         *instInfo.ZoneID,
		},
	}
	return &csi.NodeGetInfoResponse{
		NodeId:             ns.driver.GetInstanceId(),
		MaxVolumesPerNode:  ns.driver.GetMaxVolumePerNode(),
		AccessibleTopology: top,
	}, nil
}

// NodeExpandVolume will expand filesystem of volume.
// Input Parameters:
//  volume id: REQUIRED
//  volume path: REQUIRED
func (ns *NodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (
	*csi.NodeExpandVolumeResponse, error) {
	funcName := "NodeExpandVolume"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
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
	// ensure one call in-flight
	klog.Infof("Try to lock resource %s", volumeId)
	if acquired := ns.locks.TryAcquire(volumeId); !acquired {
		return nil, status.Errorf(codes.Aborted, common.OperationPendingFmt, volumeId)
	}
	defer ns.locks.Release(volumeId)
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

// NodeGetVolumeStats
// Input Arguments:
//  volume id: REQUIRED
//  volume path: REQUIRED
func (ns *NodeServer) NodeGetVolumeStats(ctx context.Context,
	req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	funcName := "NodeGetVolumeStats"
	info, hash := common.EntryFunction(funcName)
	defer klog.Info(common.ExitFunction(funcName, hash))
	klog.Info(info)
	// 0. Preflight
	// check arguments
	klog.Infof("%s: Check input arguments", hash)
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetVolumePath()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume path missing in request")
	}

	volumeId := req.GetVolumeId()
	volumePath := req.GetVolumePath()

	// Get volume info
	klog.Infof("%s: Get volume %s info", hash, volumeId)
	volInfo, err := ns.cloud.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil || volInfo.Instance == nil || volInfo.Instance.Device == nil {
		return nil, status.Errorf(codes.NotFound, "cannot find volume %s", volumeId)
	}

	// Checkout device
	klog.Infof("%s: Get device name from mount point %s", hash, volumePath)
	devicePath, _, err := mount.GetDeviceNameFromMount(ns.mounter, volumePath)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cannot get device name from mount point %s", volumePath)
	}
	klog.Infof("%s: Succeed to get device name %s", hash, devicePath)
	if devicePath == "" || *volInfo.Instance.Device != devicePath {
		return nil, status.Errorf(codes.NotFound, "device path mismatch, from mount point %s, "+
			"from cloud provider %s", devicePath, *volInfo.Instance.Device)
	}

	// Get metrics
	metricsStatFs := volume.NewMetricsStatFS(volumePath)
	metrics, err := metricsStatFs.GetMetrics()
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	klog.Infof("%s: Succeed to get metrics", hash)
	return &csi.NodeGetVolumeStatsResponse{
		Usage: []*csi.VolumeUsage{
			{
				Available: metrics.Available.Value(),
				Total:     metrics.Capacity.Value(),
				Used:      metrics.Used.Value(),
				Unit:      csi.VolumeUsage_BYTES,
			},
			{
				Available: metrics.InodesFree.Value(),
				Total:     metrics.Inodes.Value(),
				Used:      metrics.InodesUsed.Value(),
				Unit:      csi.VolumeUsage_INODES,
			},
		},
	}, nil
}

func (ns *NodeServer) getBlockSizeBytes(devicePath string) (int64, error) {
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
