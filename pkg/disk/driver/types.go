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
	DefaultVolumeType              VolumeType = SSDEnterpriseVolumeType
	HighPerformanceVolumeType      VolumeType = 0
	HighCapacityVolumeType         VolumeType = 2
	SuperHighPerformanceVolumeType VolumeType = 3
	StandardVolumeType             VolumeType = 100
	SSDEnterpriseVolumeType        VolumeType = 200
	NeonSANVolumeType              VolumeType = 5
)

type VolumeType int

func (v VolumeType) Int() int {
	return int(v)
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
	100: "Standard",
	200: "SSDEnterprise",
	5:   "NeonSAN",
}

var VolumeTypeValue = map[string]VolumeType{
	"HighPerformance":      0,
	"HighCapacity":         2,
	"SuperHighPerformance": 3,
	"Standard":             100,
	"SSDEnterprise":        200,
	"NeonSAN":              5,
}

var VolumeTypeToStepSize = map[VolumeType]int{
	0:   10,
	2:   50,
	3:   10,
	100: 10,
	200: 10,
	5:   100,
}

var VolumeTypeToMinSize = map[VolumeType]int{
	0:   10,
	2:   100,
	3:   10,
	100: 10,
	200: 10,
	5:   100,
}

var VolumeTypeToMaxSize = map[VolumeType]int{
	0:   2000,
	2:   5000,
	3:   2000,
	100: 2000,
	200: 2000,
	5:   50000,
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
	HighPerformanceInstanceType      InstanceType = 0
	SuperHighPerformanceInstanceType InstanceType = 1
	StandardInstanceType             InstanceType = 101
	EnterpriseInstanceType           InstanceType = 201
	PremiumInstanceType              InstanceType = 301
)

var InstanceTypeName = map[InstanceType]string{
	0:   "HighPerformance",
	1:   "SuperHighPerformance",
	101: "Standard",
	201: "Enterprise",
	301: "Premium",
}

var InstanceTypeValue = map[string]InstanceType{
	"HighPerformance":      0,
	"SuperHighPerformance": 1,
	"Standard":             101,
	"Enterprise":           201,
	"Premium":              301,
}

var VolumeTypeAttachConstraint = map[VolumeType][]InstanceType{
	HighPerformanceVolumeType:      {HighPerformanceInstanceType},
	SuperHighPerformanceVolumeType: {SuperHighPerformanceInstanceType},
	HighCapacityVolumeType: {
		HighPerformanceInstanceType,
		SuperHighPerformanceInstanceType,
		StandardInstanceType,
		EnterpriseInstanceType,
		PremiumInstanceType,
	},
	StandardVolumeType: {
		StandardInstanceType,
	},
	SSDEnterpriseVolumeType: {
		EnterpriseInstanceType,
		PremiumInstanceType,
	},
	NeonSANVolumeType: {
		HighPerformanceInstanceType,
		SuperHighPerformanceInstanceType,
		StandardInstanceType,
		EnterpriseInstanceType,
		PremiumInstanceType,
	},
}
