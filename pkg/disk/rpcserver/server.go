package rpcserver

import (
	"github.com/yunify/qingcloud-csi/pkg/cloudprovider"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
)

// Run
// Initial and start CSI driver
func Run(driver *driver.DiskDriver, cloud cloudprovider.CloudManager, endpoint string) {
	// Initialize default library driver
	ids := NewIdentityServer(driver, cloud)
	cs := NewControllerServer(driver, cloud)
	ns := NewNodeServer(driver, cloud)

	s := common.NewNonBlockingGRPCServer()
	s.Start(endpoint, ids, cs, ns)
	s.Wait()
}
