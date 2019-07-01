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

package disk

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	timestamp2 "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/yunify/qingcloud-csi/pkg/server"
	"github.com/yunify/qingcloud-csi/pkg/server/instance"
	"github.com/yunify/qingcloud-csi/pkg/server/snapshot"
	"github.com/yunify/qingcloud-csi/pkg/server/storageclass"
	"github.com/yunify/qingcloud-csi/pkg/server/volume"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
	cloudServer *server.ServerConfig
}

// This operation MUST be idempotent
// csi.CreateVolumeRequest: name 				+Required
//							capability			+Required
func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	glog.Info("----- Start CreateVolume -----")
	defer glog.Info("===== End CreateVolume =====")
	// 0. Prepare
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.Errorf("Invalid create volume req: %v", req)
		return nil, err
	}
	// Required volume capability
	if req.VolumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities missing in request")
	} else if !server.ContainsVolumeCapabilities(cs.Driver.GetVolumeCapabilityAccessModes(), req.GetVolumeCapabilities()) {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities not match")
	}
	// Check sanity of request Name, Volume Capabilities
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume name missing in request")
	}
	volumeName := req.GetName()

	// create VolumeManager object
	vm, err := volume.NewVolumeManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// create StorageClass object
	sc, err := storageclass.NewQingStorageClassFromMap(req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// get request volume capacity range
	requiredByte := req.GetCapacityRange().GetRequiredBytes()
	requiredGib := sc.FormatVolumeSize(server.ByteCeilToGib(requiredByte), sc.VolumeStepSize)
	limitByte := req.GetCapacityRange().GetLimitBytes()
	if limitByte == 0 {
		limitByte = server.Int64Max
	}
	// check volume range
	if server.GibToByte(requiredGib) < requiredByte || server.GibToByte(requiredGib) > limitByte || requiredGib < sc.
		VolumeMinSize || requiredGib > sc.VolumeMaxSize {
		glog.Errorf("Request capacity range [%d, %d] bytes, storage class capacity range [%d, %d] GB, format required size: %d gb",
			requiredByte, limitByte, sc.VolumeMinSize, sc.VolumeMaxSize, requiredGib)
		return nil, status.Error(codes.OutOfRange, "Unsupport capacity range")
	}

	// should not fail when requesting to create a volume with already exisiting name and same capacity
	// should fail when requesting to create a volume with already exisiting name and different capacity.
	exVol, err := vm.FindVolumeByName(volumeName)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Find volume by name error: %s, %s", volumeName, err.Error()))
	}
	if exVol != nil {
		glog.Infof("Request volume name: %s, capacity range [%d,%d] bytes, type: %d, zone: %s",
			volumeName, requiredByte, limitByte, sc.VolumeType, vm.GetZone())
		glog.Infof("Exist volume name: %s, id: %s, capacity: %d bytes, type: %d, zone: %s",
			*exVol.VolumeName, *exVol.VolumeID, server.GibToByte(*exVol.Size), *exVol.VolumeType, vm.GetZone())
		if *exVol.Size >= requiredGib && int64(*exVol.Size)*server.Gib <= limitByte && *exVol.VolumeType == sc.
			VolumeType {
			// exisiting volume is compatible with new request and should be reused.
			return &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					VolumeId:      *exVol.VolumeID,
					CapacityBytes: int64(*exVol.Size) * server.Gib,
					VolumeContext: req.GetParameters(),
				},
			}, nil
		}
		return nil, status.Error(codes.AlreadyExists,
			fmt.Sprintf("Volume %s already exsit but is incompatible", volumeName))
	}

	// do create volume
	glog.Infof("Creating volume %s with %d GB in zone %s...", volumeName, requiredGib, vm.GetZone())
	volumeId, err := vm.CreateVolume(volumeName, requiredGib, *sc)
	if err != nil {
		return nil, err
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volumeId,
			CapacityBytes: int64(requiredGib) * server.Gib,
			VolumeContext: req.GetParameters(),
		},
	}, nil
}

// This operation MUST be idempotent
// volume id is REQUIRED in csi.DeleteVolumeRequest
func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	glog.Info("----- Start DeleteVolume -----")
	defer glog.Info("===== End DeleteVolume =====")
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.Errorf("invalid delete volume req: %v", req)
		return nil, err
	}
	// Check sanity of request Name, Volume Capabilities
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume id missing in request")
	}
	// For now the image get unconditionally deleted, but here retention policy can be checked
	volumeId := req.GetVolumeId()

	// Deleting disk
	glog.Infof("deleting volume %s", volumeId)
	// Create VolumeManager object
	vm, err := volume.NewVolumeManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// For idempotent:
	// MUST reply OK when volume does not exist
	volInfo, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return &csi.DeleteVolumeResponse{}, nil
	}
	// Is volume in use
	if *volInfo.Status == volume.DiskStatusInuse {
		return nil, status.Errorf(codes.FailedPrecondition, "volume is in use by another resource")
	}
	// Do delete volume
	glog.Infof("Deleting volume %s status %s in zone %s...", volumeId, *volInfo.Status, vm.GetZone())
	// When return with retry message at deleting volume, retry after several seconds.
	// Retry times is 10.
	// Retry interval is changed from 1 second to 10 seconds.
	for i := 1; i <= 10; i++ {
		err = vm.DeleteVolume(volumeId)
		if err != nil {
			glog.Infof("Failed to delete disk volume: %s in %s with error: %v", volumeId, vm.GetZone(), err)
			if strings.Contains(err.Error(), server.RetryString) {
				time.Sleep(time.Duration(i*2) * time.Second)
			} else {
				return nil, status.Error(codes.Internal, err.Error())
			}
		} else {
			return &csi.DeleteVolumeResponse{}, nil
		}
	}
	return nil, status.Error(codes.Internal, "Exceed retry times: "+err.Error())
}

// csi.ControllerPublishVolumeRequest: 	volume id 			+ Required
//										node id				+ Required
//										volume capability 	+ Required
//										readonly			+ Required (This field is NOT provided when requesting in Kubernetes)
func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	glog.Info("----- Start ControllerPublishVolume -----")
	defer glog.Info("===== End ControllerPublishVolume =====")
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME); err != nil {
		glog.Errorf("invalid publish volume req: %v", req)
		return nil, err
	}
	// 0. Preflight
	// check volume id arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	// check nodeId arguments
	if len(req.GetNodeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Node ID missing in request")
	}
	// check volume capability
	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "No volume capability is provided ")
	}

	// create volume manager object
	vm, err := volume.NewVolumeManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// create instance manager object
	im, err := instance.NewInstanceManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// if volume id not exist
	volumeId := req.GetVolumeId()
	exVol, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exVol == nil {
		return nil, status.Errorf(codes.NotFound, "Volume: %s does not exist", volumeId)
	}

	// if instance id not exist
	nodeId := req.GetNodeId()
	exIns, err := im.FindInstance(nodeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exIns == nil {
		return nil, status.Errorf(codes.NotFound, "Node: %s does not exist", nodeId)
	}

	// Volume published to another node
	if len(*exVol.Instance.InstanceID) != 0 && *exVol.Instance.InstanceID != nodeId {
		return nil, status.Error(codes.FailedPrecondition, "Volume published to another node")
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume capability missing in request")
	}
	// 1. Attach
	// attach volume
	glog.Infof("Attaching volume %s to instance %s in zone %s...", volumeId, nodeId, vm.GetZone())
	err = vm.AttachVolume(volumeId, nodeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// When return with retry message at describe volume, retry after several seconds.
	// Retry times is 3.
	// Retry interval is changed from 1 second to 3 seconds.
	for i := 1; i <= 3; i++ {
		volInfo, err := vm.FindVolume(volumeId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		// check device path
		if *volInfo.Instance.Device != "" {
			// found device path
			glog.Infof("Attaching volume %s on instance %s succeed.", volumeId, nodeId)
			return &csi.ControllerPublishVolumeResponse{}, nil
		} else {
			// cannot found device path
			glog.Infof("Cannot find device path and retry to find volume device %s", volumeId)
			time.Sleep(time.Duration(i) * time.Second)
		}
	}
	// Cannot find device path
	// Try to detach volume
	glog.Infof("Cannot find device path and going to detach volume %s", volumeId)
	if err := vm.DetachVolume(volumeId, nodeId); err != nil {
		return nil, status.Errorf(codes.Internal,
			"cannot find device path, detach volume %s failed", volumeId)
	} else {
		return nil, status.Errorf(codes.Internal,
			"cannot find device path, volume %s has been detached, please try attaching to instance %s again.",
			volumeId, nodeId)
	}
}

// This operation MUST be idempotent
// csi.ControllerUnpublishVolumeRequest: 	volume id	+Required
func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	glog.Info("----- Start ControllerUnpublishVolume -----")
	defer glog.Info("===== End ControllerUnpublishVolume =====")
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME); err != nil {
		glog.Errorf("invalid unpublish volume req: %v", req)
		return nil, err
	}
	// 0. Preflight
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	volumeId := req.GetVolumeId()
	nodeId := req.GetNodeId()

	// 1. Detach
	// create volume provisioner object
	vm, err := volume.NewVolumeManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// create instance manager object
	im, err := instance.NewInstanceManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check volume exist
	exVol, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exVol == nil {
		return nil, status.Errorf(codes.NotFound, "Volume: %s does not exist", volumeId)
	}

	// check node exist
	exIns, err := im.FindInstance(nodeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exIns == nil {
		return nil, status.Errorf(codes.NotFound, "Node: %s does not exist", nodeId)
	}

	// do detach
	glog.Infof("Detaching volume %s to instance %s in zone %s...", volumeId, nodeId, vm.GetZone())
	err = vm.DetachVolume(volumeId, nodeId)
	if err != nil {
		glog.Errorf("failed to detach disk image: %s from instance %s with error: %v",
			volumeId, nodeId, err)
		return nil, err
	}

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.ValidateVolumeCapabilitiesRequest: 	volume id 			+ Required
// 											volume capability 	+ Required
func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	glog.Info("----- Start ValidateVolumeCapabilities -----")
	defer glog.Info("===== End ValidateVolumeCapabilities =====")

	// require volume id parameter
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No volume id is provided")
	}

	// require capability parameter
	if len(req.GetVolumeCapabilities()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No volume capabilities are provided")
	}

	// check volume exist
	vm, err := volume.NewVolumeManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	volumeId := req.GetVolumeId()
	vol, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if vol == nil {
		return nil, status.Errorf(codes.NotFound, "Volume %s does not exist", volumeId)
	}

	// check capability
	for _, c := range req.GetVolumeCapabilities() {
		found := false
		for _, c1 := range cs.Driver.GetVolumeCapabilityAccessModes() {
			if c1.GetMode() == c.GetAccessMode().GetMode() {
				found = true
			}
		}
		if !found {
			return &csi.ValidateVolumeCapabilitiesResponse{
				Message: "Driver does not support mode:" + c.GetAccessMode().GetMode().String(),
			}, nil
		}
	}

	return &csi.ValidateVolumeCapabilitiesResponse{}, nil
}

// ControllerExpandVolume allows the CO to expand the size of a volume
// volume id is REQUIRED in csi.ControllerExpandVolumeRequest
// capacity range is REQUIRED in csi.ControllerExpandVolumeRequest
func (cs *controllerServer) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest,
) (*csi.ControllerExpandVolumeResponse, error) {
	defer server.EntryFunction("ControllerExpandVolume")()
	// 0. check input args
	// require volume id parameter
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "No volume id is provided")
	}

	vm, err := volume.NewVolumeManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 1. Check volume status
	// does volume exist
	volumeId := req.GetVolumeId()
	volInfo, err := vm.FindVolume(volumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil {
		return nil, status.Errorf(codes.NotFound, "Volume: %s does not exist", volumeId)
	}
	// volume in use
	if *volInfo.Status == volume.DiskStatusInuse {
		return nil, status.Errorf(codes.FailedPrecondition,
			"Volume [%s] currently published on a node but plugin only support OFFLINE expansion", volumeId)
	}

	// 2. Get capacity
	volTypeInt := *volInfo.VolumeType
	if volTypeStr, ok := server.VolumeTypeToString[volTypeInt]; ok == true {
		glog.Infof("Succeed to get volume [%s] type [%s]", volumeId, volTypeStr)
	} else {
		glog.Errorf("Unsupported volume [%s] type [%d]", volumeId, volTypeInt)
		return nil, status.Errorf(codes.Internal, "Unsupported volume [%s] type [%d]", volumeId, volTypeInt)
	}
	volTypeMinSize := server.VolumeTypeToMinSize[volTypeInt]
	volTypeMaxSize := server.VolumeTypeToMaxSize[volTypeInt]
	requiredByte := req.GetCapacityRange().GetRequiredBytes()
	requiredGib := server.FormatVolumeSize(volTypeInt, server.ByteCeilToGib(requiredByte))
	limitByte := req.GetCapacityRange().GetLimitBytes()
	if limitByte == 0 {
		limitByte = server.Int64Max
	}
	// check volume range
	if server.GibToByte(requiredGib) < requiredByte || server.GibToByte(requiredGib) > limitByte ||
		requiredGib < volTypeMinSize || requiredGib > volTypeMaxSize {
		glog.Errorf("Request capacity range [%d, %d] bytes, storage class capacity range [%d, %d] GB, format required size: %d gb",
			requiredByte, limitByte, volTypeMinSize, volTypeMaxSize, requiredGib)
		return nil, status.Error(codes.OutOfRange, "Unsupport capacity range")
	}

	// 3. Expand volume
	err = vm.ResizeVolume(volumeId, requiredGib)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         int64(requiredGib) * server.Gib,
		NodeExpansionRequired: true,
	}, nil
}

func (cs *controllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (cs *controllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// CreateSnapshot allows the CO to create a snapshot.
// This operation MUST be idempotent.
// 1. If snapshot successfully cut and ready to use, the plugin MUST reply 0 OK.
// 2. If an error occurs before a snapshot is cut, the plugin SHOULD reply a corresponding error code.
// 3. If snapshot successfully cut but still being precessed,
// the plugin SHOULD return 0 OK and ready_to_use SHOULD be set to false.
// Source volume id is REQUIRED
// Snapshot name is REQUIRED
func (cs *controllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	glog.Info("----- Start CreateSnapshot -----")
	defer glog.Info("===== End CreateSnapshot =====")
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); err != nil {
		glog.Errorf("invalid create snapshot request: %v", req)
		return nil, err
	}
	// 0. Preflight
	// Check source volume id
	if len(req.GetSourceVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "volume ID missing in request")
	}
	// Check snapshot name
	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "snapshot name missing in request")
	}

	// Create snapshot manager object
	sm, err := snapshot.NewSnapshotManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	sourceVolumeId := req.GetSourceVolumeId()
	snapshotName := req.GetName()
	var timestamp *timestamp2.Timestamp
	var isReadyToUse bool
	// For idempotent
	// If a snapshot corresponding to the specified snapshot name is successfully cut and ready to use (meaning it MAY
	// be specified as a volume_content_source in a CreateVolumeRequest), the Plugin MUST reply 0 OK with the
	// corresponding CreateSnapshotResponse.
	exSnap, err := sm.FindSnapshotByName(snapshotName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "find snapshot by name error: %s, %s", snapshotName, err.Error())
	}
	if exSnap != nil {
		glog.Infof("Exist snapshot name: %s, snapshot id %s, source volume id %s",
			*exSnap.SnapshotName, *exSnap.SnapshotID, *exSnap.Resource.ResourceID)
		if exSnap.Resource != nil && *exSnap.Resource.ResourceType == "volume" &&
			*exSnap.Resource.ResourceID == sourceVolumeId {
			timestamp, err = ptypes.TimestampProto(*exSnap.CreateTime)
			if err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			if *exSnap.Status == snapshot.SnapshotStatusAvailable {
				isReadyToUse = true
			} else {
				isReadyToUse = false
			}
			return &csi.CreateSnapshotResponse{
				Snapshot: &csi.Snapshot{
					SizeBytes:      int64(*exSnap.Size) * server.Mib,
					SnapshotId:     *exSnap.SnapshotID,
					SourceVolumeId: *exSnap.Resource.ResourceID,
					CreationTime:   timestamp,
					ReadyToUse:     isReadyToUse,
				},
			}, nil
		}
		return nil, status.Errorf(codes.AlreadyExists,
			"snapshot name=[%s] id=[%s] already exists, but is incompatible with the volume id=[%s]",
			snapshotName, *exSnap.SnapshotID, sourceVolumeId)
	}
	// Create a new full snapshot
	glog.Infof("Creating snapshot %s from volume %s in zone %s...", snapshotName, sourceVolumeId, sm.GetZone())
	snapId, err := sm.CreateSnapshot(snapshotName, sourceVolumeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create snapshot [%s] from source volume [%s] error: %s",
			snapshotName, sourceVolumeId, err.Error())
	}
	snapInfo, err := sm.FindSnapshot(snapId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "find snapshot [%s] error: %s", snapId, err.Error())
	}
	timestamp, err = ptypes.TimestampProto(*snapInfo.CreateTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if *snapInfo.Status == snapshot.SnapshotStatusAvailable {
		isReadyToUse = true
	} else {
		isReadyToUse = false
	}
	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SizeBytes:      int64(*snapInfo.Size) * server.Mib,
			SnapshotId:     *snapInfo.SnapshotID,
			SourceVolumeId: *snapInfo.Resource.ResourceID,
			CreationTime:   timestamp,
			ReadyToUse:     isReadyToUse,
		},
	}, nil
}

// CreateSnapshot allows the CO to delete a snapshot.
// This operation MUST be idempotent.
// Snapshot id is REQUIRED
func (cs *controllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	glog.Info("----- Start DeleteSnapshot -----")
	defer glog.Info("===== End DeleteSnapshot =====")
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); err != nil {
		glog.Errorf("invalid create snapshot request: %v", req)
		return nil, err
	}
	// 0. Preflight
	// Check snapshot id
	if len(req.GetSnapshotId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "snapshot ID missing in request")
	}
	snapshotId := req.GetSnapshotId()

	// Create snapshot manager object
	sm, err := snapshot.NewSnapshotManagerFromFile(cs.cloudServer.GetConfigFilePath())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// 1. For idempotent:
	// MUST reply OK when snapshot does not exist
	exSnap, err := sm.FindSnapshot(snapshotId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exSnap == nil {
		return &csi.DeleteSnapshotResponse{}, nil
	}
	// 2. Delete snapshot
	glog.Infof("Deleting snapshot %s in zone %s...", snapshotId, sm.GetZone())
	// When return with retry message at deleting snapshot, retry after several seconds.
	// Retry times is 10.
	// Retry interval is changed from 1 second to 10 seconds.
	for i := 1; i <= 10; i++ {
		err = sm.DeleteSnapshot(snapshotId)
		if err != nil {
			glog.Infof("Failed to delete snapshot %s in %s with error: %v", snapshotId, sm.GetZone(), err)
			if strings.Contains(err.Error(), server.RetryString) {
				time.Sleep(time.Duration(i*2) * time.Second)
			} else {
				return nil, status.Error(codes.Internal, err.Error())
			}
		} else {
			return &csi.DeleteSnapshotResponse{}, nil
		}
	}
	return nil, status.Error(codes.Internal, "Exceed retry times: "+err.Error())
}

func (cs *controllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
