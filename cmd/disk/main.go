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

package main

import (
	"flag"
	"github.com/yunify/qingcloud-csi/pkg/cloudprovider"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"github.com/yunify/qingcloud-csi/pkg/disk/rpcserver"
	"k8s.io/klog"
	"math/rand"
	"os"
	"time"
)

const (
	version              = "v1.1.0"
	defaultProvisionName = "disk.csi.qingcloud.com"
	defaultConfigPath    = "/etc/config/config.yaml"
)

var (
	endpoint   = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName = flag.String("drivername", defaultProvisionName, "name of the driver")
	nodeId     = flag.String("nodeid", "",
		"If driver cannot get instance ID from /etc/qingcloud/instance-id, we would use this flag.")
	configPath = flag.String("config", defaultConfigPath, "server config file path")
	maxVolume  = flag.Int64("maxvolume", 10,
		"Maximum number of volumes that controller can publish to the node.")
	timeout = flag.Duration("timeout", time.Second*60, "timeout duration for retrying, default 60s")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	handle()
	os.Exit(0)
}

func handle() {
	// Get Instance Id
	instanceId, err := driver.GetInstanceIdFromFile(driver.DefaultInstanceIdFilePath)
	if err != nil {
		klog.Warningf("Failed to get instance id from file, use --nodeId flag. error: %s", err)
		instanceId = *nodeId
	}
	// Get qingcloud config object
	cloud, err := cloudprovider.NewQingCloudManagerFromFile(*configPath)
	if err != nil {
		klog.Fatal(err)
	}

	// Set DiskDriverInput
	diskDriverInput := &driver.InitDiskDriverInput{
		Name:          *driverName,
		Version:       version,
		NodeId:        instanceId,
		MaxVolume:     *maxVolume,
		VolumeCap:     driver.DefaultVolumeAccessModeType,
		ControllerCap: driver.DefaultControllerServiceCapability,
		NodeCap:       driver.DefaultNodeServiceCapability,
		PluginCap:     driver.DefaultPluginCapability,
	}

	// For resize
	mounter := common.NewSafeMounter()
	driver := driver.GetDiskDriver()
	driver.InitDiskDriver(diskDriverInput)
	rpcserver.Run(driver, cloud, mounter, *endpoint)
}
