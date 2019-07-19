package common

import (
	"k8s.io/kubernetes/pkg/volume"
	"testing"
)

func TestMetrics(t *testing.T) {
	volumePath := "/mnt"
	metricsStatFs := volume.NewMetricsStatFS(volumePath)
	metrics, err := metricsStatFs.GetMetrics()
	if err != nil {
		t.Log(err.Error())
	}

	t.Logf("inode value %d %d %d", metrics.Inodes.Value(), metrics.InodesFree.Value(), metrics.InodesUsed.Value())
	inode64, _ := metrics.Inodes.AsInt64()
	inode64free, _ := metrics.InodesFree.AsInt64()
	inode64used, _ := metrics.InodesUsed.AsInt64()
	t.Logf("inode as int64 %d %d %d", inode64, inode64free, inode64used)

	t.Logf("capacity value %d %d %d", metrics.Capacity.Value(), metrics.Available.Value(), metrics.Used.Value())
	cap64, _ := metrics.Capacity.AsInt64()
	cap64free, _ := metrics.Available.AsInt64()
	cap64used, _ := metrics.Used.AsInt64()
	t.Logf("capacity as int64 %d %d %d", cap64, cap64free, cap64used)
}
