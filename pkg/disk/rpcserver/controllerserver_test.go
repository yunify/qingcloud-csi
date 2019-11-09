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
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"github.com/yunify/qingcloud-sdk-go/service"
	"reflect"
	"testing"
)

func getFakeVolume() *service.Volume {
	volId := "vol-riv17xkh"
	volName := "pvc-016fa900-e142-41cf-b29b-58d1ccf31147"
	volType := 200
	volRepl := "rpp-00000002"
	volSize := 10
	volZone := "pek3b"
	return &service.Volume{
		VolumeID:   &volId,
		VolumeName: &volName,
		VolumeType: &volType,
		Repl:       &volRepl,
		Size:       &volSize,
		ZoneID:     &volZone,
	}
}

func getMockControllerServer() *ControllerServer {
	d := driver.GetDiskDriver()
	d.InitDiskDriver(
		&driver.InitDiskDriverInput{
			Name:      "disk.csi.qingcloud",
			Version:   "v1.1.0",
			NodeId:    "i-12345678",
			MaxVolume: 10,
		},
	)
	return NewControllerServer(d, nil, DefaultBackOff, 5)
}

func TestDiskControllerServer_PickTopology(t *testing.T) {
	cs := getMockControllerServer()
	tests := []struct {
		name       string
		topRequire *csi.TopologyRequirement
		topResult  *driver.Topology
	}{
		{
			name: "normal in Kubernetes",
			topRequire: &csi.TopologyRequirement{
				Requisite: []*csi.Topology{
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3b",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3c",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3d",
						},
					},
				},
				Preferred: []*csi.Topology{
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3b",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3c",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3d",
						},
					},
				},
			},
			topResult: driver.NewTopology("pek3b", driver.Enterprise1InstanceType),
		},
		{
			name: "csi spec example 2",
			topRequire: &csi.TopologyRequirement{
				Requisite: []*csi.Topology{
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3b",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3c",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3d",
						},
					},
				},
				Preferred: []*csi.Topology{
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3c",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3d",
						},
					},
				},
			},
			topResult: driver.NewTopology("pek3c", driver.Enterprise1InstanceType),
		},
	}
	for _, test := range tests {
		res, _ := cs.PickTopology(test.topRequire)
		if !reflect.DeepEqual(test.topResult, res) {
			t.Errorf("name %s: expect %v, but actually %v", test.name, test.topResult, res)
		}
	}
}

func TestDiskControllerServer_IsValidTopology(t *testing.T) {
	cs := getMockControllerServer()
	vol := getFakeVolume()

	tests := []struct {
		name       string
		volInfo    *service.Volume
		topRequire *csi.TopologyRequirement
		isValid    bool
	}{
		{
			name:       "nil topology",
			volInfo:    vol,
			topRequire: nil,
			isValid:    true,
		},
		{
			name:    "valid topology",
			volInfo: vol,
			topRequire: &csi.TopologyRequirement{
				Requisite: []*csi.Topology{
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3b",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3c",
						},
					},
				},
			},
			isValid: true,
		},
		{
			name:    "wrong zone",
			volInfo: vol,
			topRequire: &csi.TopologyRequirement{
				Requisite: []*csi.Topology{
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3d",
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
							cs.driver.GetTopologyZoneKey():         "pek3c",
						},
					},
				},
			},
			isValid: false,
		},
		{
			name:    "only instance type",
			volInfo: vol,
			topRequire: &csi.TopologyRequirement{
				Requisite: []*csi.Topology{
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.Enterprise1InstanceType],
						},
					},
					{
						Segments: map[string]string{
							cs.driver.GetTopologyInstanceTypeKey(): driver.InstanceTypeName[driver.HighPerformanceInstanceType],
						},
					},
				},
			},
			isValid: false,
		},
		{
			name:    "non topology requirement",
			volInfo: vol,
			topRequire: &csi.TopologyRequirement{
				Requisite: []*csi.Topology{},
			},
			isValid: true,
		},
	}

	for _, test := range tests {
		res := cs.IsValidTopology(test.volInfo, test.topRequire)
		if test.isValid != res {
			t.Errorf("name %s: expect %t but actually %t", test.name, test.isValid, res)
		}
	}
}

func TestDiskControllerServer_GetVolumeTopology(t *testing.T) {
	cs := getMockControllerServer()
	vol := getFakeVolume()
	tests := []struct {
		name     string
		volume   *service.Volume
		topology []*csi.Topology
	}{
		{
			name:   "normal",
			volume: vol,
			topology: []*csi.Topology{
				{
					Segments: map[string]string{
						cs.driver.GetTopologyInstanceTypeKey(): "SuperHighPerformance",
						cs.driver.GetTopologyZoneKey():         "pek3b",
					},
				},
				{
					Segments: map[string]string{
						cs.driver.GetTopologyInstanceTypeKey(): "Enterprise1",
						cs.driver.GetTopologyZoneKey():         "pek3b",
					},
				},
				{
					Segments: map[string]string{
						cs.driver.GetTopologyInstanceTypeKey(): "Enterprise2",
						cs.driver.GetTopologyZoneKey():         "pek3b",
					},
				},
				{
					Segments: map[string]string{
						cs.driver.GetTopologyInstanceTypeKey(): "Premium",
						cs.driver.GetTopologyZoneKey():         "pek3b",
					},
				},
			},
		},
		{
			name:     "nil",
			volume:   nil,
			topology: nil,
		},
	}
	for _, v := range tests {
		res := cs.GetVolumeTopology(v.volume)
		if !reflect.DeepEqual(res, v.topology) {
			t.Errorf("name %s: expect %v, but actually %v", v.name, v.topology, res)
		}
	}
}
