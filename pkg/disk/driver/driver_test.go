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
	"testing"
)

func TestDiskDriver_ValidatePluginCapability(t *testing.T) {
	tests := []struct {
		name    string
		driver  DiskDriver
		cap     csi.PluginCapability_Service_Type
		isValid bool
	}{
		{
			name: "check topology",
			driver: DiskDriver{
				pluginCap: []*csi.PluginCapability{
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
				},
			},
			cap:     csi.PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS,
			isValid: true,
		},
		{
			name: "check controller service",
			driver: DiskDriver{
				pluginCap: []*csi.PluginCapability{
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
				},
			},
			cap:     csi.PluginCapability_Service_CONTROLLER_SERVICE,
			isValid: true,
		},
		{
			name: "check expansion",
			driver: DiskDriver{
				pluginCap: []*csi.PluginCapability{
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
				},
			},
			cap:     csi.PluginCapability_Service_UNKNOWN,
			isValid: false,
		},
		{
			name: "lake of controller service",
			driver: DiskDriver{
				pluginCap: []*csi.PluginCapability{
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
				},
			},
			cap:     csi.PluginCapability_Service_CONTROLLER_SERVICE,
			isValid: false,
		},
	}
	for _, test := range tests {
		isValid := test.driver.ValidatePluginCapabilityService(test.cap)
		if test.isValid != isValid {
			t.Errorf("testcase %s: expect %t but actually %t", test.name, test.isValid, isValid)
		}
	}
}
