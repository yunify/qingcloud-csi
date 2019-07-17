package common

import (
	"github.com/golang/glog"
	"time"
)

// EntryFunction print timestamps
// TODO: set log prefix, k8s.io/klog/klogr
func EntryFunction(functionName string) func() {
	start := time.Now()
	glog.Infof("*************** enter %s at %s ***************", functionName, start.String())
	return func() {
		glog.Infof("=============== exit %s (%s since %s) ===============", functionName, time.Since(start),
			start.String())
	}
}
