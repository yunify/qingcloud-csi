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

package block

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/kubernetes-csi/drivers/pkg/csi-common"
	"github.com/yunify/qingcloud-csi/pkg/server"
)

const version = "0.2.0"

type block struct {
	driver *csicommon.CSIDriver

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability

	server *server.ServerConfig
}

// GetBlockDriver
// Create block driver
func GetBlockDriver() *block {
	return &block{}
}

// NewIdentityServer
// Create identity server
func NewIdentityServer(d *csicommon.CSIDriver, svr *server.ServerConfig) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
		server:                svr,
	}
}

// NewControllerServer
// Create controller server
func NewControllerServer(d *csicommon.CSIDriver, svr *server.ServerConfig) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
		server:                  svr,
	}
}

// NewNodeServer
// Create node server
func NewNodeServer(d *csicommon.CSIDriver, svr *server.ServerConfig) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
		server:            svr,
	}
}

// Run
// Initial and start CSI driver
func (blk *block) Run(driverName, nodeID, endpoint string, server *server.ServerConfig) {
	glog.Infof("Driver: %v version: %v", driverName, version)

	// Initialize default library driver
	blk.driver = csicommon.NewCSIDriver(driverName, version, nodeID)
	if blk.driver == nil {
		glog.Fatalln("Failed to initialize CSI Driver.")
	}

	blk.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	})
	blk.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	// Create GRPC servers
	blk.ids = NewIdentityServer(blk.driver, server)
	blk.ns = NewNodeServer(blk.driver, server)
	blk.cs = NewControllerServer(blk.driver, server)

	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(endpoint, blk.ids, blk.cs, blk.ns)
	s.Wait()

}
