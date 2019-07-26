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
