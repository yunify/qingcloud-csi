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
