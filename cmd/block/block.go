package main

import (
	"flag"
	"os"
)

func init(){
	flag.Set("logtostderr", "true")
}

var (
	endpoint   = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName = flag.String("drivername", "csi-qingcloud-performance", "name of the driver")
	nodeID     = flag.String("nodeid", "", "node id")
)

func main() {
	flag.Parse()
	handle()
	os.Exit(0)
}

func handle(){

}
