/*
Copyright (C) 2019 Yunify, Inc.

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

package rpcserver

import (
	"context"
	"flag"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/cloud/mock"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"reflect"
	"testing"
)

func init() {
	klog.InitFlags(nil)
	flag.Parse()
}

func TestIdentityServer_Probe(t *testing.T) {
	tests := []struct {
		name  string
		zones map[string]*qcservice.Zone
		err   error
	}{
		{
			name: "normal",
			zones: map[string]*qcservice.Zone{
				"mock1": {},
				"mock2": {},
			},
			err: nil,
		},
		{
			name:  "failed",
			zones: nil,
			err:   status.Error(codes.FailedPrecondition, "cannot find any zones"),
		},
	}

	for _, test := range tests {
		cm := &mock.MockCloudManager{}
		cm.SetZones(test.zones)
		is := NewIdentityServer(nil, cm)
		resp, err := is.Probe(context.Background(), &csi.ProbeRequest{})
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("testcase %s: expect %s, but actually %t", test.name, test.err, err)
		}

		if err == nil && resp.GetReady().GetValue() != true {
			t.Errorf("testcase %s: expect %t, but actually %t", test.name, true, resp.GetReady().GetValue())
		}
	}
}

func TestIdentityServer_GetPluginCapabilities(t *testing.T) {
	tests := []struct {
		name   string
		config *driver.InitDiskDriverInput
	}{
		{
			name: "normal",
			config: &driver.InitDiskDriverInput{
				PluginCap: []*csi.PluginCapability{
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
						Type: &csi.PluginCapability_VolumeExpansion_{
							VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
								Type: csi.PluginCapability_VolumeExpansion_ONLINE,
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
		},
		{
			name:   "empty",
			config: &driver.InitDiskDriverInput{},
		},
	}

	for _, test := range tests {
		driver := driver.GetDiskDriver()
		driver.InitDiskDriver(test.config)
		is := NewIdentityServer(driver, nil)
		resp, _ := is.GetPluginCapabilities(context.Background(), &csi.GetPluginCapabilitiesRequest{})
		if !reflect.DeepEqual(resp.GetCapabilities(), test.config.PluginCap) {
			t.Errorf("testcase %s: expect cap %v, but actually %v", test.name, test.config.PluginCap, resp.GetCapabilities())
		}
	}
}

func TestIdentityServer_GetPluginInfo(t *testing.T) {
	tests := []struct {
		name   string
		config *driver.InitDiskDriverInput
		err    error
	}{
		{
			name: "normal",
			config: &driver.InitDiskDriverInput{
				Name:    "test-driver",
				Version: "v19.2.0",
			},
			err: nil,
		},
		{
			name: "lack of driver name",
			config: &driver.InitDiskDriverInput{
				Version: "v19.2.0",
			},
			err: status.Error(codes.Unavailable, "Driver name not configured"),
		},
		{
			name: "lack of driver version",
			config: &driver.InitDiskDriverInput{
				Name: "mock_driver",
			},
			err: status.Error(codes.Unavailable, "Driver is missing version"),
		},
	}

	for _, test := range tests {
		driver := driver.GetDiskDriver()
		driver.InitDiskDriver(test.config)
		is := NewIdentityServer(driver, nil)
		resp, err := is.GetPluginInfo(context.Background(), &csi.GetPluginInfoRequest{})
		if !reflect.DeepEqual(test.err, err) {
			t.Errorf("testcase %s: expect error %s, but actually %s", test.name, test.err, err)
		}
		if err == nil && resp.GetName() != test.config.Name {
			t.Errorf("testcase %s: expect name %s, but actually %s", test.name, test.config.Name, resp.GetName())
		}
		if err == nil && resp.GetVendorVersion() != test.config.Version {
			t.Errorf("testcase %s: expect version %s, but actually %s", test.name, test.config.Version, resp.GetVendorVersion())
		}
	}
}
