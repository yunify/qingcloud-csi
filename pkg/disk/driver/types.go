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
	HighPerformanceVolumeType      VolumeType = 0
	HighCapacityVolumeType         VolumeType = 2
	SuperHighPerformanceVolumeType VolumeType = 3
	NeonSANVolumeType              VolumeType = 5
	NeonSANHDDVolumeType           VolumeType = 6
	NeonSANRDMAVolumeType          VolumeType = 7
	ThirdPartyStorageType          VolumeType = 20
	StandardVolumeType             VolumeType = 100
	SSDEnterpriseVolumeType        VolumeType = 200
	DefaultVolumeType                         = StandardVolumeType
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
	HighPerformanceVolumeType:      "HighPerformance",
	HighCapacityVolumeType:         "HighCapacity",
	SuperHighPerformanceVolumeType: "SuperHighPerformance",
	NeonSANVolumeType:              "NeonSAN",
	NeonSANHDDVolumeType:           "NeonSANHDD",
	NeonSANRDMAVolumeType:          "NeonSANRDMA",
	ThirdPartyStorageType:          "ThirdPartyStorage",
	StandardVolumeType:             "Standard",
	SSDEnterpriseVolumeType:        "SSDEnterprise",
}

var VolumeTypeToStepSize = map[VolumeType]int{
	HighPerformanceVolumeType:      10,
	HighCapacityVolumeType:         10,
	SuperHighPerformanceVolumeType: 10,
	NeonSANVolumeType:              10,
	NeonSANHDDVolumeType:           10,
	NeonSANRDMAVolumeType:          10,
	ThirdPartyStorageType:          10,
	StandardVolumeType:             10,
	SSDEnterpriseVolumeType:        10,
}

var VolumeTypeToMinSize = map[VolumeType]int{
	HighPerformanceVolumeType:      20,
	HighCapacityVolumeType:         20,
	SuperHighPerformanceVolumeType: 20,
	NeonSANVolumeType:              20,
	NeonSANHDDVolumeType:           20,
	NeonSANRDMAVolumeType:          450,
	ThirdPartyStorageType:          10,
	StandardVolumeType:             10,
	SSDEnterpriseVolumeType:        20,
}

var VolumeTypeToMaxSize = map[VolumeType]int{
	HighPerformanceVolumeType:      2000,
	HighCapacityVolumeType:         5000,
	SuperHighPerformanceVolumeType: 2000,
	NeonSANVolumeType:              24000,
	NeonSANHDDVolumeType:           32000,
	NeonSANRDMAVolumeType:          32000,
	ThirdPartyStorageType:          200,
	StandardVolumeType:             2000,
	SSDEnterpriseVolumeType:        2000,
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
	SANInstanceType                     InstanceType = 3
	SuperHighPerformanceSANInstanceType InstanceType = 6
	HighPerformanceSANInstanceType      InstanceType = 7
	StandardInstanceType                InstanceType = 101
	Enterprise1InstanceType             InstanceType = 201
	Enterprise2InstanceType             InstanceType = 202
	EnterpriseCompute3InstanceType      InstanceType = 203
	PremiumInstanceType                 InstanceType = 301
	NvidiaAmpereG3InstanceType          InstanceType = 1003
)

var InstanceTypeName = map[InstanceType]string{
	0:    "HighPerformance",
	1:    "SuperHighPerformance",
	3:    "SAN",
	6:    "SuperHighPerformanceSAN",
	7:    "HighPerformanceSAN",
	101:  "Standard",
	201:  "Enterprise1",
	202:  "Enterprise2",
	203:  "EnterpriseCompute3",
	301:  "Premium",
	1003: "NvidiaAmpereG3",
}

var InstanceTypeValue = map[string]InstanceType{
	"HighPerformance":         0,
	"SuperHighPerformance":    1,
	"SAN":                     3,
	"SuperHighPerformanceSAN": 6,
	"HighPerformanceSAN":      7,
	"Standard":                101,
	"Enterprise1":             201,
	"Enterprise2":             202,
	"EnterpriseCompute3":      203,
	"Premium":                 301,
	"NvidiaAmpereG3":          1003,
}

var InstanceTypeAttachPreferred = map[InstanceType]VolumeType{
	HighPerformanceInstanceType:         HighPerformanceVolumeType,
	SuperHighPerformanceInstanceType:    SuperHighPerformanceVolumeType,
	SANInstanceType:                     NeonSANVolumeType,
	SuperHighPerformanceSANInstanceType: NeonSANVolumeType,
	HighPerformanceSANInstanceType:      NeonSANRDMAVolumeType,
	StandardInstanceType:                StandardVolumeType,
	Enterprise1InstanceType:             SSDEnterpriseVolumeType,
	Enterprise2InstanceType:             SSDEnterpriseVolumeType,
	EnterpriseCompute3InstanceType:      SSDEnterpriseVolumeType,
	PremiumInstanceType:                 SSDEnterpriseVolumeType,
	NvidiaAmpereG3InstanceType:          NeonSANHDDVolumeType,
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
		NvidiaAmpereG3InstanceType,
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
		SANInstanceType,
		SuperHighPerformanceSANInstanceType,
		StandardInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		EnterpriseCompute3InstanceType,
		PremiumInstanceType,
		NvidiaAmpereG3InstanceType,
	},
	NeonSANHDDVolumeType: {
		HighPerformanceInstanceType,
		SuperHighPerformanceInstanceType,
		SANInstanceType,
		HighPerformanceSANInstanceType,
		StandardInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		EnterpriseCompute3InstanceType,
		PremiumInstanceType,
		NvidiaAmpereG3InstanceType,
	},
	NeonSANRDMAVolumeType: {
		HighPerformanceInstanceType,
		SuperHighPerformanceInstanceType,
		SANInstanceType,
		HighPerformanceSANInstanceType,
		StandardInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		EnterpriseCompute3InstanceType,
		PremiumInstanceType,
		NvidiaAmpereG3InstanceType,
	},
	ThirdPartyStorageType: {
		StandardInstanceType,
	},
	StandardVolumeType: {
		HighPerformanceInstanceType,
		StandardInstanceType,
	},
	SSDEnterpriseVolumeType: {
		SuperHighPerformanceInstanceType,
		Enterprise1InstanceType,
		Enterprise2InstanceType,
		EnterpriseCompute3InstanceType,
		PremiumInstanceType,
	},
}

const (
	DiskSingleReplicaType  int = 1
	DiskMultiReplicaType   int = 2
	DefaultDiskReplicaType int = DiskMultiReplicaType
)
