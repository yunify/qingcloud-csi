package block

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
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
	} else if !HasSameAccessMode(cs.Driver.GetVolumeCapabilityAccessModes(), req.GetVolumeCapabilities()) {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities not match")
	}
	// Check sanity of request Name, Volume Capabilities
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume name missing in request")
	}
	volumeName := req.GetName()


	// create VolumeManager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// create StorageClass object
	sc, err := NewQingStorageClassFromMap(req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// get request volume capacity range
	requireByte := req.GetCapacityRange().GetRequiredBytes()
	requireGb := sc.formatVolumeSize(ByteCeilToGb(requireByte))
	limitByte := req.GetCapacityRange().GetLimitBytes()
	if limitByte == 0 {
		limitByte = Int64_Max
	}
	// check volume range
	if GbToByte(requireGb) > limitByte || requireGb > sc.VolumeMaxSize{
		glog.Errorf("Request capacity range [%d, %d] bytes, storage class capacity range [%d, %d] GB, format required size: %d gb",
			requireByte, limitByte, sc.VolumeMinSize, sc.VolumeMaxSize, requireGb)
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
			volumeName, requireByte, limitByte, sc.VolumeType, vm.volumeService.Config.Zone)
		glog.Infof("Exist volume name: %s, id: %s, capacity: %d GB, type: %d, zone: %s",
			*exVol.VolumeName, *exVol.VolumeID, GbToByte(*exVol.Size), *exVol.VolumeType, vm.volumeService.Config.Zone)
		if *exVol.Size >= requireGb && int64(*exVol.Size)*gib <= limitByte && *exVol.VolumeType == sc.VolumeType{
			// exisiting volume is compatible with new request and should be reused.
			return &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					Id:            *exVol.VolumeID,
					CapacityBytes: int64(*exVol.Size) * gib,
					Attributes:    req.GetParameters(),
				},
			}, nil
		}
		return nil, status.Error(codes.AlreadyExists,
			fmt.Sprintf("Volume %s already exsit but is incompatible", volumeName))
	}

	// do create volume
	glog.Infof("Creating volume %s with %d GB in zone %s...", volumeName, requireGb, vm.volumeService.Config.Zone)
	volumeId, err := vm.CreateVolume(volumeName, requireGb, *sc)
	if err != nil {
		return nil, err
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            volumeId,
			CapacityBytes: int64(requireGb) * gib,
			Attributes:    req.GetParameters(),
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

	// Deleting block image
	glog.Infof("deleting volume %s", volumeId)
	// Create VolumeManager object
	vm, err := NewVolumeManager()
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
	if *volInfo.Status == BlockVolume_Status_INUSE{
		return nil, status.Errorf(codes.FailedPrecondition, "volume is in use by another resource")
	}
	// Do delete volume
	glog.Infof("Deleting volume %s status %s in zone %s...", volumeId, *volInfo.Status, vm.volumeService.Config.Zone)
	// When return with retry message at deleting volume, retry after several seconds.
	// Retry times is 10.
	// Retry interval is changed from 1 second to 10 seconds.
	for i := 1; i <= 10; i++ {
		err = vm.DeleteVolume(volumeId)
		if err != nil {
			glog.Infof("Failed to delete block volume: %s in %s with error: %v", volumeId, vm.volumeService.Config.Zone, err)
			if strings.Contains(err.Error(), RetryString) {
				time.Sleep(time.Duration(i) * time.Second)
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
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// create instance manager object
	im, err := NewInstanceManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// if volume id not exist
	volumeId := req.GetVolumeId()
	exVol, err := vm.FindVolume(volumeId)
	if err != nil{
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exVol == nil {
		return nil, status.Errorf(codes.NotFound, "Volume: %s does not exist", volumeId)
	}

	// if instance id not exist
	nodeId := req.GetNodeId()
	exIns, err := im.FindInstance(nodeId)
	if err != nil{
		return nil, status.Error(codes.Internal, err.Error())
	}
	if exIns == nil {
		return nil, status.Errorf(codes.NotFound, "Node: %s does not exist", nodeId)
	}else{
		if *exIns.Status!= Instance_Status_RUNNING{
			return nil, status.Errorf(codes.NotFound, "Node: %d does not running", nodeId)
		}
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
	glog.Infof("Attaching volume %s to instance %s in zone %s...", volumeId, nodeId, vm.volumeService.Config.Zone)
	err = vm.AttachVolume(volumeId, nodeId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	glog.Infof("Attaching volume %s succeed.", volumeId)

	return &csi.ControllerPublishVolumeResponse{}, nil
}

// This operation MUST be idempotent
// csi.ControllerUnpublishVolumeRequest: 	volume id	+Required
//
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
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// do detach
	glog.Infof("Detaching volume %s to instance %s in zone %s...", volumeId, nodeId, vm.volumeService.Config.Zone)
	err = vm.DetachVolume(volumeId, nodeId)
	if err != nil {
		glog.Errorf("failed to detach block image: %s from instance %s with error: %v",
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
	vm, err := NewVolumeManager()
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
				Supported: false,
				Message:   "Driver does not support mode:" + c.GetAccessMode().GetMode().String(),
			}, nil
		}
		// TODO: Ignoring mount & block tyeps for now.
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Supported: true,
	}, nil
}
