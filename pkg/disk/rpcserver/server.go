package rpcserver

import (
	"github.com/yunify/qingcloud-csi/pkg/cloudprovider"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"k8s.io/kubernetes/pkg/util/mount"
)

// Run
// Initial and start CSI drivermakee

func Run(driver *driver.DiskDriver, cloud cloudprovider.CloudManager, mounter *mount.SafeFormatAndMount,
	endpoint string) {
	// Initialize default library driver
	ids := NewIdentityServer(driver, cloud)
	cs := NewControllerServer(driver, cloud)
	ns := NewNodeServer(driver, cloud, mounter)

	s := common.NewNonBlockingGRPCServer()
	s.Start(endpoint, ids, cs, ns)
	s.Wait()
}
