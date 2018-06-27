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

func (cs *controllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	glog.Info("Run CreateVolume")
	// 0. Prepare
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("Invalid create volume req: %v", req)
		return nil, err
	}
	// Check sanity of request Name, Volume Capabilities
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume Name cannot be empty")
	}
	if req.VolumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume Capabilities cannot be empty")
	}
	volumeName := req.GetName()

	// create StorageClass object
	sc, err := NewQingStorageClassFromMap(req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// create VolumeManager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Need to check for already existing volume name, and if found
	// check for the requested capacity and already allocated capacity
	if exVol, err := vm.FindVolumeByName(volumeName); err == nil && exVol != nil {
		// Since err is nil, it means the volume with the same name already exists
		// need to check if the size of exisiting volume is the same as in new
		// request
		if int64(*exVol.Size)*gib >= req.GetCapacityRange().GetRequiredBytes() {
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
	} else if err != nil {
		return nil, status.Error(codes.Internal,
			fmt.Sprintf("Find volume by name error %s, %s", volumeName, err.Error()))
	}
	// Get volume size
	volSizeBytes := int64(gib)
	if req.GetVolumeCapabilities() != nil {
		volSizeBytes = int64(req.GetCapacityRange().GetRequiredBytes())
	}
	volSizeGB := int(volSizeBytes / gib)

	// Create volume
	volumeId, err := vm.CreateVolume(volumeName, volSizeGB, *sc)
	if err != nil {
		return nil, err
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            volumeId,
			CapacityBytes: int64(volSizeGB) * gib,
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
	// For now the image get unconditionally deleted, but here retention policy can be checked
	volumeId := req.GetVolumeId()

	// Deleting block image
	glog.Infof("deleting volume %s", volumeId)
	// Create VolumeManager object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
	// check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}
	if len(req.GetNodeId()) == 0{
		return nil, status.Error(codes.InvalidArgument, "Node ID missing in request")
	}
	volumeId := req.GetVolumeId()
	nodeId := req.GetNodeId()

	// 1. Attach
	// create volume provisioner object
	vm, err := NewVolumeManager()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
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
