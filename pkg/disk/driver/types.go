package driver

import "github.com/container-storage-interface/spec/lib/go/csi"

const (
	DefaultInstanceIdFilePath = "/etc/qingcloud/instance-id"
)

var DefaultVolumeAccessModeType = []csi.VolumeCapability_AccessMode_Mode{
	csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
}
var DefaultControllerServiceCapability = []csi.ControllerServiceCapability_RPC_Type{
	csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
	csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
}
var DefaultNodeServiceCapability = []csi.NodeServiceCapability_RPC_Type{
	csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
	csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
}
var DefaultPluginCapability = []*csi.PluginCapability{
	{
		Type: &csi.PluginCapability_Service_{
			Service: &csi.PluginCapability_Service{
				Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
			},
		},
	},
	{
		Type: &csi.PluginCapability_VolumeExpansion_{
			VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
				Type: csi.PluginCapability_VolumeExpansion_OFFLINE,
			},
		},
	},
}

const (
	HighPerformanceDiskType      int = 0
	HighCapacityDiskType         int = 2
	SuperHighPerformanceDiskType int = 3
	StandardDiskType             int = 100
	SSDEnterpriseDiskType        int = 200
	NeonSANDiskType              int = 5
)

// convert volume type to string
// https://docs.qingcloud.com/product/api/action/volume/create_volumes.html
var VolumeTypeName = map[int]string{
	0:   "HighPerformance",
	2:   "HighCapacity",
	3:   "SuperHighPerformance",
	100: "Standard",
	200: "SSDEnterprise",
	5:   "NeonSAN",
}

var VolumeTypeToStepSize = map[int]int{
	0:   10,
	2:   50,
	3:   10,
	100: 10,
	200: 10,
	5:   100,
}

var VolumeTypeToMinSize = map[int]int{
	0:   10,
	2:   100,
	3:   10,
	100: 10,
	200: 10,
	5:   100,
}

var VolumeTypeToMaxSize = map[int]int{
	0:   2000,
	2:   5000,
	3:   2000,
	100: 2000,
	200: 2000,
	5:   50000,
}
