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

package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
)

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
	csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
}

var DefaultNodeServiceCapability = []csi.NodeServiceCapability_RPC_Type{
	csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
	csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
	csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
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
	{
		Type: &csi.PluginCapability_Service_{
			Service: &csi.PluginCapability_Service{
				Type: csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
			},
		},
	},
}

const (
	DefaultVolumeType              VolumeType = StandardVolumeType
	HighPerformanceVolumeType      VolumeType = 0
	HighCapacityVolumeType         VolumeType = 2
	SuperHighPerformanceVolumeType VolumeType = 3
	NeonSANVolumeType              VolumeType = 5
	NeonSANHDDVolumeType           VolumeType = 6
	StandardVolumeType             VolumeType = 100
	SSDEnterpriseVolumeType        VolumeType = 200
)

type VolumeType int

func (v VolumeType) Int() int {
	return int(v)
}

func (v VolumeType) String() string {
	return VolumeTypeName[v]
}

func (v VolumeType) ValidateAttachedOn(i InstanceType) bool {
	for _, iType := range VolumeTypeAttachConstraint[v] {
		if iType == i {
			return true
		}
	}
	return false
}

func (v VolumeType) IsValid() bool {
	if _, ok := VolumeTypeName[v]; !ok {
		return false
	} else {
		return true
	}
}

// convert volume type to string
// https://docs.qingcloud.com/product/api/action/volume/create_volumes.html
var VolumeTypeName = map[VolumeType]string{
	0:   "HighPerformance",
	2:   "HighCapacity",
	3:   "SuperHighPerformance",
	5:   "NeonSAN",
	6:   "NeonSANHDD",
	100: "Standard",
	200: "SSDEnterprise",
}

var VolumeTypeValue = map[string]VolumeType{
	"HighPerformance":      0,
	"HighCapacity":         2,
	"SuperHighPerformance": 3,
	"NeonSAN":              5,
	"NeonSANHDD":           6,
	"Standard":             100,
	"SSDEnterprise":        200,
}

var VolumeTypeToStepSize = map[VolumeType]int{
	0:   10,
	2:   50,
	3:   10,
	5:   100,
	6:   100,
	100: 10,
	200: 10,
}

var VolumeTypeToMinSize = map[VolumeType]int{
	0:   10,
	2:   100,
	3:   10,
	5:   100,
	6:   100,
	100: 10,
	200: 10,
}

var VolumeTypeToMaxSize = map[VolumeType]int{
	0:   2000,
	2:   5000,
	3:   2000,
	5:   50000,
	6:   50000,
	100: 2000,
	200: 2000,
}

type InstanceType int

func (i InstanceType) Int() int {
	return int(i)
}

func (i InstanceType) IsValid() bool {
	if _, ok := InstanceTypeName[i]; !ok {
		return false
	} else {
		return true
	}
}

const (
	HighPerformanceInstanceType         InstanceType = 0
	SuperHighPerformanceInstanceType    InstanceType = 1
	SuperHighPreformanceSANInstanceType InstanceType = 6
	HighPerformanceSANInstanceType      InstanceType = 7
	StandardInstanceType                InstanceType = 101
	Enterprise1InstanceType             InstanceType = 201
	Enterprise2InstanceType             InstanceType = 202
	PremiumInstanceType                 InstanceType = 301
)

var InstanceTypeName = map[InstanceType]string{
	0:   "HighPerformance",
	1:   "SuperHighPerformance",
	6:   "SuperHighPerformanceSAN",
	7:   "HighPerformanceSAN",
	101: "Standard",
	201: "Enterprise1",
	202: "Enterprise2",
	301: "Premium",
}

var InstanceTypeValue = map[string]InstanceType{
	"HighPerformance":         0,
	"SuperHighPerformance":    1,
	"SuperHighPerformanceSAN": 6,
	"HighPerformanceSAN":      7,
	"Standard":                101,
	"Enterprise1":             201,
	"Enterprise2":             202,
	"Premium":                 301,
}

var VolumeTypeAttachConstraint = map[VolumeType][]InstanceType{
	HighPerformanceVolumeType: {
		HighPerformanceInstanceType,
		StandardInstanceType,
	},
	HighCapacityVolumeType: {
		HighPerformanceInstanceType,
		SuperHighPerformanceInstanceType,
		StandardInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		PremiumInstanceType,
	},
	SuperHighPerformanceVolumeType: {
		SuperHighPerformanceInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		PremiumInstanceType,
	},
	NeonSANVolumeType: {
		HighPerformanceInstanceType,
		SuperHighPerformanceInstanceType,
		SuperHighPreformanceSANInstanceType,
		StandardInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		PremiumInstanceType,
	},
	NeonSANHDDVolumeType: {
		HighPerformanceInstanceType,
		SuperHighPerformanceInstanceType,
		HighPerformanceSANInstanceType,
		StandardInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		PremiumInstanceType,
	},
	StandardVolumeType: {
		HighPerformanceInstanceType,
		StandardInstanceType,
	},
	SSDEnterpriseVolumeType: {
		SuperHighPerformanceInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		PremiumInstanceType,
	},
}

const (
	DiskSingleReplicaType  int = 1
	DiskMultiReplicaType   int = 2
	DefaultDiskReplicaType int = DiskMultiReplicaType
)
