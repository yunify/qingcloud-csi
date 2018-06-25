package block

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"path"
)

type controllerServer struct {
	*csicommon.DefaultControllerServer
}

func (cs *controllerServer) CreateVolume(
	ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("invalid create volume req: %v", req)
		return nil, err
	}
	// Check sanity of request Name, Volume Capabilities
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume Name cannot be empty")
	}
	if req.VolumeCapabilities == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume Capabilities cannot be empty")
	}

	// Create QingCloud storage class object
	sc, err := NewStorageClassFromMap(req.GetParameters())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// Create volume provisioner object
	vp, err := newVolumeProvisioner(sc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// Need to check for already existing volume name, and if found
	// check for the requested capacity and already allocated capacity
	if exVol, err := vp.findVolumeByName(req.GetName()); err == nil && exVol != nil {
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
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("Volume with the same name: %s but with different size already exist", req.GetName()))
	} else if err != nil {
		return nil, status.Error(codes.Internal,
			fmt.Sprintf("Find volume by name error %s, %v", req.GetName(), err.Error()))
	}

	// Create QingCloud volume
	newVol := blockVolume{req.Name, "", 0, sc.Zone, *sc}

	// Get volume size
	volSizeBytes := int64(gib)
	if req.GetVolumeCapabilities() != nil {
		volSizeBytes = int64(req.GetCapacityRange().GetRequiredBytes())
	}
	volSizeGB := int(volSizeBytes / gib)

	// Create volume
	err = vp.CreateVolume(volSizeGB, &newVol)
	if err != nil {
		return nil, err
	}

	// Store volInfo into a persistent file.
	if err := persistVolInfo(newVol.VolID, path.Join(PluginFolder, "controller"), &newVol); err != nil {
		glog.Warningf("failed to store volInfo with error: %v", err)
	}
	blockVolumes[newVol.VolID] = newVol
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			Id:            newVol.VolID,
			CapacityBytes: int64(newVol.VolSize) * gib,
			Attributes:    req.GetParameters(),
		},
	}, nil
}

func (cs *controllerServer) DeleteVolume(
	ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.Warningf("invalid delete volume req: %v", req)
		return nil, err
	}
	// For now the image get unconditionally deleted, but here retention policy can be checked
	volumeID := req.GetVolumeId()
	blockVol := &blockVolume{}
	if err := loadVolInfo(volumeID, path.Join(PluginFolder, "controller"), blockVol); err != nil {
		return nil, err
	}
	// Deleting block image
	glog.Infof("deleting volume %s", blockVol.VolName)
	// Create volume provisioner object
	vp, err := newVolumeProvisioner(&blockVol.Sc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// Delete block volume
	if err = vp.DeleteVolume(blockVol.VolID); err != nil {
		glog.Infof("failed to delete block volume: %s in %s with error: %v", blockVol.VolID, blockVol.Zone, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	// Remove persistent storage file for the unmapped volume
	if err := deleteVolInfo(blockVol.VolID, path.Join(PluginFolder, "controller")); err != nil {
		return nil, err
	}
	delete(blockVolumes, blockVol.VolID)
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerUnpublishVolume(
	ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	glog.Infof("ControllerUnpublishVolume")
	return &csi.ControllerUnpublishVolumeResponse{}, nil
}

func (cs *controllerServer) ControllerPublishVolume(
	ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	glog.Infof("ControllerPublishVolume")
	return &csi.ControllerPublishVolumeResponse{}, nil
}
