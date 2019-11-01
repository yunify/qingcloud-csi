/*
Copyright (C) 2019 Yunify, Inc.

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

package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"reflect"
	"testing"
)

func TestIsValidReplica(t *testing.T) {
	testcases := []struct {
		name    string
		replica int
		expect  bool
	}{
		{
			name:    "single",
			replica: DiskSingleReplicaType,
			expect:  true,
		},
		{
			name:    "multi",
			replica: DiskMultiReplicaType,
			expect:  true,
		},
		{
			name:    "fake1",
			replica: 0,
			expect:  false,
		},
		{
			name:    "fake2",
			replica: 3,
			expect:  false,
		},
	}

	for _, v := range testcases {
		res := IsValidReplica(v.replica)
		if res != v.expect {
			t.Errorf("test %s: expect %t, but actually %t", v.name, v.expect, res)
		}
	}
}

func TestIsValidFileSystemType(t *testing.T) {
	testcases := []struct {
		name   string
		fsType string
		expect bool
	}{
		{
			name:   "EXT3",
			fsType: common.FileSystemExt3,
			expect: true,
		},
		{
			name:   "EXT4",
			fsType: common.FileSystemExt4,
			expect: true,
		},
		{
			name:   "XFS",
			fsType: common.FileSystemXfs,
			expect: true,
		},
		{
			name:   "ext5",
			fsType: "ext5",
			expect: false,
		},
		{
			name:   "Ext3",
			fsType: "Ext3",
			expect: false,
		},
		{
			name:   "NTFS",
			fsType: "NTFS",
			expect: false,
		},
	}

	for _, v := range testcases {
		res := IsValidFileSystemType(v.fsType)
		if res != v.expect {
			t.Errorf("test %s: expect %t, but actually %t", v.name, v.expect, res)
		}
	}
}

func TestNewVolumeCapabilityAccessMode(t *testing.T) {
	tests := []struct {
		name string
		amm  csi.VolumeCapability_AccessMode_Mode
		am   *csi.VolumeCapability_AccessMode
	}{
		{
			name: "single node write",
			amm:  csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			am: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
		},
		{
			name: "multi node multi writer",
			amm:  csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
			am: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
			},
		},
	}
	for _, test := range tests {
		res := NewVolumeCapabilityAccessMode(test.amm)
		if !reflect.DeepEqual(res, test.am) {
			t.Errorf("name %s: expect %s, but actually %s", test.name, test.am.String(), res.GetMode().String())
		}
	}
}

func TestNewControllerServiceCapability(t *testing.T) {
	tests := []struct {
		name string
		csct csi.ControllerServiceCapability_RPC_Type
		csc  *csi.ControllerServiceCapability
	}{
		{
			name: "volume",
			csct: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			csc: &csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
					},
				},
			},
		},
		{
			name: "snapshot",
			csct: csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
			csc: &csi.ControllerServiceCapability{
				Type: &csi.ControllerServiceCapability_Rpc{
					Rpc: &csi.ControllerServiceCapability_RPC{
						Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
					},
				},
			},
		},
	}
	for _, test := range tests {
		res := NewControllerServiceCapability(test.csct)
		if !reflect.DeepEqual(res.GetType(), test.csc.GetType()) {
			t.Errorf("name %s: expect %s, but actually %s", test.name, test.csc.GetType(), res.GetType())
		}
	}
}

func TestNewNodeServiceCapability(t *testing.T) {
	tests := []struct {
		name string
		nsct csi.NodeServiceCapability_RPC_Type
		nsc  *csi.NodeServiceCapability
	}{
		{
			name: "expand volume",
			nsct: csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
			nsc: &csi.NodeServiceCapability{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_EXPAND_VOLUME,
					},
				},
			},
		},
		{
			name: "volume stats",
			nsct: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
			nsc: &csi.NodeServiceCapability{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
					},
				},
			},
		},
	}
	for _, test := range tests {
		res := NewNodeServiceCapability(test.nsct)
		if !reflect.DeepEqual(res.GetType(), test.nsc.GetType()) {
			t.Errorf("name %s: expect %s, but actually %s", test.name, test.nsct.String(), res.String())
		}
	}
}
