package block

import "testing"
import (
	"flag"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"golang.org/x/net/context"
)

func TestCsiCreateVolume(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	drv := csicommon.NewCSIDriver("fake", "fake", "fake")
	cs := controllerServer{csicommon.NewDefaultControllerServer(drv)}
	req := csi.CreateVolumeRequest{}
	req.Name = "wx-sanity"
	req.VolumeCapabilities = []*csi.VolumeCapability{
		{nil, &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER}},
	}
	req.CapacityRange = &csi.CapacityRange{1 * gib, 0}
	cs.Driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME})
	cs.Driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER})
	_, err := cs.CreateVolume(context.Background(), &req)

	if err != nil {
		t.Errorf(err.Error())
	} else {
		t.Logf("Pass")
	}
}
