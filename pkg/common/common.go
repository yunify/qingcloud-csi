package common

import (
	"k8s.io/klog"
	"time"
)

// EntryFunction print timestamps
// TODO: set log prefix, k8s.io/klog/klogr
func EntryFunction(functionName string) func() {
	start := time.Now()
	klog.Infof("*************** enter %s at %s ***************", functionName, start.String())
	return func() {
		klog.Infof("=============== exit %s (%s since %s) ===============", functionName, time.Since(start),
			start.String())
	}
}
