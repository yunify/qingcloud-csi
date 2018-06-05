package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/yunify/qingcloud-csi/pkg/block"
	"os"
	"path"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	endpoint   = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName = flag.String("drivername", "csi-qingcloud", "name of the driver")
	nodeID     = flag.String("nodeid", "", "node id")
)

func main() {
	flag.Parse()

	handle()
	os.Exit(0)
}

func handle() {
	if err := block.CreatePath(path.Join(block.PluginFolder, "controller")); err != nil {
		glog.Errorf("failed to create directory for controller %v", err)
		os.Exit(1)
	}
	if err := block.CreatePath(path.Join(block.PluginFolder, "node")); err != nil {
		glog.Errorf("failed to create directory for node %v", err)
		os.Exit(1)
	}
	driver := block.GetBlockDriver()
	driver.Run(*driverName, *nodeID, *endpoint)
}
