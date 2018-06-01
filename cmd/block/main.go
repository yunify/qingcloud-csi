package main

import (
	"flag"
	"os"
	"github.com/wnxn/qingcloud-csi/pkg/block"
)

func init(){
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

func handle(){
	driver := block.GetBlockDriver()
	driver.Run(*driverName, *nodeID, *endpoint)
}
