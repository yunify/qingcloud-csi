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
