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

package rpcserver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/yunify/qingcloud-csi/pkg/cloud"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

type IdentityServer struct {
	driver *driver.DiskDriver
	cloud  cloud.CloudManager
}

// NewIdentityServer
// Create identity server
func NewIdentityServer(d *driver.DiskDriver, c cloud.CloudManager) *IdentityServer {
	return &IdentityServer{
		driver: d,
		cloud:  c,
	}
}

var _ csi.IdentityServer = &IdentityServer{}

// Plugin MUST implement this RPC call
func (is *IdentityServer) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	zones, err := is.cloud.GetZoneList()
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	klog.V(5).Infof("get active zone lists [%v]", zones)
	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{Value: true},
	}, nil
}

// Get plugin capabilities: CONTROLLER, ACCESSIBILITY, EXPANSION
func (d *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.
	GetPluginCapabilitiesResponse, error) {
	klog.V(5).Infof("Using default capabilities")
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: d.driver.GetPluginCapability(),
	}, nil
}

func (d *IdentityServer) GetPluginInfo(ctx context.Context,
	req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	klog.V(5).Infof("Using GetPluginInfo")

	if d.driver.GetName() == "" {
		return nil, status.Error(codes.Unavailable, "Driver name not configured")
	}

	if d.driver.GetVersion() == "" {
		return nil, status.Error(codes.Unavailable, "Driver is missing version")
	}

	return &csi.GetPluginInfoResponse{
		Name:          d.driver.GetName(),
		VendorVersion: d.driver.GetVersion(),
	}, nil
}
