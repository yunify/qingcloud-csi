// +-------------------------------------------------------------------------
// | Copyright (C) 2018 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package disk

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/yunify/qingcloud-csi/pkg/server"
)

const version = "v1.1.0"

type disk struct {
	driver *csicommon.CSIDriver

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability

	cloud *server.ServerConfig
}

// GetDiskDriver
// Create disk driver
func GetDiskDriver() *disk {
	return &disk{}
}

// NewIdentityServer
// Create identity server
func NewIdentityServer(d *csicommon.CSIDriver, svr *server.ServerConfig) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
		cloudServer:           svr,
	}
}

// NewControllerServer
// Create controller server
func NewControllerServer(d *csicommon.CSIDriver, svr *server.ServerConfig) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
		cloudServer:             svr,
	}
}

// NewNodeServer
// Create node server
func NewNodeServer(d *csicommon.CSIDriver, svr *server.ServerConfig) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
		cloudServer:       svr,
	}
}

// Run
// Initial and start CSI driver
func (d *disk) Run(driverName, nodeID, endpoint string, serverConfig *server.ServerConfig) {
	glog.Infof("Driver: %v version: %v", driverName, version)

	// Initialize default library driver
	d.driver = csicommon.NewCSIDriver(driverName, version, nodeID)
	if d.driver == nil {
		glog.Fatalln("Failed to initialize CSI Driver.")
	}

	d.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	})
	d.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	// Create GRPC servers
	d.ids = NewIdentityServer(d.driver, serverConfig)
	d.ns = NewNodeServer(d.driver, serverConfig)
	d.cs = NewControllerServer(d.driver, serverConfig)

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(endpoint, d.ids, d.cs, d.ns)
	s.Wait()

}
