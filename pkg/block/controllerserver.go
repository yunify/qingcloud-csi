package block

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

// csi.CreateVolumeRequest: name 				+Required
//							capability			+Required

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	glog.Info("Run CreateVolume")
	// 0. Prepare
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("Invalid create volume req: %v", req)
		return nil, err
	}
	// Required volume capability
	if req.VolumeCapabilities == nil  {
		return nil, status.Error(codes.InvalidArgument, "Volume capabilities missing in request")
	}else if !HasSameAccessMode(cs.Driver.GetVolumeCapabilityAccessModes(), req.GetVolumeCapabilities()){
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
	if limitByte == 0{
		limitByte = Int64_Max
	}

	// should not fail when requesting to create a volume with already exisiting name and same capacity
	// should fail when requesting to create a volume with already exisiting name and different capacity.
	if exVol, err:= vm.FindVolumeByName(volumeName); err == nil && exVol != nil{
		glog.Warningf("Volume name %s with capacity [%d,%d] already exist with volume Id %s capacity %d",
			volumeName, requireByte, limitByte, *exVol.VolumeID, int64(*exVol.Size) * gib)
		if *exVol.Size >= requireGb && int64(*exVol.Size)*gib <= limitByte{
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
			fmt.Sprintf("Volume with the same name: %s but with different size already exist", volumeName))
	}else if err != nil {
		return nil, status.Error(codes.Internal,
			fmt.Sprintf("Find volume by name error %s, %s", volumeName, err.Error()))
	}

	// Create volume
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

func (cs *controllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	glog.Info("Run DeleteVolume")
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.Warningf("invalid delete volume req: %v", req)
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
	// For sanity: should succeed when an invalid volume id is used
	volInfo, err := vm.FindVolume(volumeId)
	if err != nil{
		return nil, status.Error(codes.Internal, err.Error())
	}
	if volInfo == nil{
		return &csi.DeleteVolumeResponse{}, nil
	}
	// Delete block volume
	if err = vm.DeleteVolume(volumeId); err != nil {
		glog.Infof("Failed to delete block volume: %s in %s with error: %v", volumeId, vm.volumeService.Config.Zone, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	glog.Infof("Run ControllerPublishVolume")
	// 0. Preflight
	// check volume id arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	// check nodeId arguments
	if len(req.GetNodeId()) == 0{
		return nil, status.Error(codes.InvalidArgument, "Node ID missing in request")
	}
	// check no volume capability
	if req.GetVolumeCapability() == nil{
		return nil, status.Error(codes.InvalidArgument, "No volume capability is provided ")
	}
	// create volume manager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// create instance manager object
	im, err:= NewInstanceManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// if volume id not exist
	volumeId := req.GetVolumeId()
	exVol, err :=vm.FindVolume(volumeId)
	if err == nil && exVol == nil{
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Volume %s does not exist", volumeId))
	}

	// if instance id not exist
	nodeId := req.GetNodeId()
	if exIns, err := im.FindInstance(nodeId) ; err == nil && exIns == nil{
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Instance %s does not exist", nodeId))
	}
	// Volume published to another node
	if len(*exVol.Instance.InstanceID) != 0 && *exVol.Instance.InstanceID != nodeId{
		return nil, status.Error(codes.FailedPrecondition, "Volume published to another node")
	}


	if req.GetVolumeCapability() == nil{
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

func (cs *controllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	glog.Infof("Run ControllerUnpublishVolume")
	// 0. Preflight
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetNodeId()) == 0{
		return nil, status.Error(codes.InvalidArgument, "Node ID missing in request")
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
	err = vm.DetachVolume(volumeId, nodeId)
	if err != nil {
		glog.Errorf("failed to detach block image: %s from instance %s with error: %v",
			volumeId,nodeId, err)
		return nil, err
	}

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (cs *controllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	glog.V(5).Infof("Using default ValidateVolumeCapabilities")
	// check input arguments
	if len(req.GetVolumeId()) == 0{
		return nil, status.Error(codes.InvalidArgument, "No volume id is provided")
	}
	if len(req.GetVolumeCapabilities()) == 0{
		return nil, status.Error(codes.InvalidArgument, "No volume capabilities are provided")
	}
	for _, c := range req.GetVolumeCapabilities() {
		found := false
		for _, c1 := range cs.Driver.GetVolumeCapabilityAccessModes(){
			if c1.GetMode() == c.GetAccessMode().GetMode() {
				found = true
			}
		}
		if !found {
			return &csi.ValidateVolumeCapabilitiesResponse{
				Supported: false,
				Message:   "Driver doesnot support mode:" + c.GetAccessMode().GetMode().String(),
			}, status.Error(codes.InvalidArgument, "Driver doesnot support mode:"+c.GetAccessMode().GetMode().String())
		}
		// TODO: Ignoring mount & block tyeps for now.
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Supported: true,
	}, nil
}
