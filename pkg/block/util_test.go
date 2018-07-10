package block

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"testing"
)

func TestContainsVolumeCapability(t *testing.T) {
	tests := []struct {
		name         string
		accessModes  []*csi.VolumeCapability_AccessMode
		capabilities *csi.VolumeCapability
		result       bool
	}{
		{
			name: "Driver: SINGLE_NODE_WRITER, Req: SINGLE_NODE_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
			},
			capabilities: &csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
			result: true,
		},
		{
			name: "Driver: SINGLE_NODE_WRITER, Req: MULTI_NODE_MULTI_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
			},
			capabilities: &csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER}},
			result: false,
		},
		{
			name: "Driver: SINGLE_NODE_WRITER, MULTI_NODE_MULTI_WRITER, Req: MULTI_NODE_MULTI_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER},
			},
			capabilities: &csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER}},
			result: true,
		},
		{
			name: "Driver: MULTI_NODE_MULTI_WRITER, MULTI_NODE_READER_ONLY, Req: MULTI_NODE_READER_ONLY",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER},
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY},
			},
			capabilities: &csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY}},
			result: true,
		},
		{
			name: "Driver: MULTI_NODE_READER_ONLY, Req: SINGLE_NODE_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY},
			},
			capabilities: &csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
			result: false,
		},
	}
	for _, v := range tests {
		res := ContainsVolumeCapability(v.accessModes, v.capabilities)
		if res != v.result {
			t.Errorf("test %s: expect %t, but result was %t", v.name, v.result, res)
		}
	}
}

func TestContainsVolumeCapabilities(t *testing.T) {
	tests := []struct {
		name         string
		accessModes  []*csi.VolumeCapability_AccessMode
		capabilities []*csi.VolumeCapability
		result       bool
	}{
		{
			name: "Driver: SINGLE_NODE_WRITER, Req: SINGLE_NODE_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
			},
			result: true,
		},
		{
			name: "Driver: SINGLE_NODE_WRITER, Req: MULTI_NODE_MULTI_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER}},
			},
		},
		{
			name: "Driver: SINGLE_NODE_WRITER, Req: MULTI_NODE_READER_ONLY",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY}},
			},
		},
		{
			name: "Driver: SINGLE_NODE_WRITER, MULTI_NODE_MULTI_WRITER, Req: MULTI_NODE_MULTI_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER}},
			},
			result: true,
		},
		{
			name: "Driver: MULTI_NODE_MULTI_WRITER, MULTI_NODE_READER_ONLY, Req: MULTI_NODE_READER_ONLY",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER},
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY}},
			},
			result: true,
		},
		{
			name: "Driver: MULTI_NODE_READER_ONLY, Req: SINGLE_NODE_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
			},
			result: false,
		},
		{
			name: "Driver: SINGLE_NODE_WRITER, Req: SINGLE_NODE_WRITER,MULTI_NODE_READER_ONLY",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY}},
			},
			result: false,
		},
		{
			name: "Driver: SINGLE_NODE_WRITER, MULTI_NODE_WRITER, Req: MULTI_NODE_MULTI_WRITER, SINGLE_NODE_WRITER",
			accessModes: []*csi.VolumeCapability_AccessMode{
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER},
				&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER},
			},
			capabilities: []*csi.VolumeCapability{
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER}},
				&csi.VolumeCapability{
					AccessMode: &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
			},
			result: true,
		},
	}
	for _, v := range tests {
		res := ContainsVolumeCapabilities(v.accessModes, v.capabilities)
		if res != v.result {
			t.Errorf("test %s: expect %t, but result was %t", v.name, v.result, res)
		}
	}
}

func TestContainsNodeServiceCapability(t *testing.T) {
	tests := []struct {
		name     string
		nodeCaps []*csi.NodeServiceCapability
		subCap   csi.NodeServiceCapability_RPC_Type
		result   bool
	}{
		{
			name: "Node Caps: STAGE_UNSTAGE, ",
			nodeCaps: []*csi.NodeServiceCapability{
				{
					Type: &csi.NodeServiceCapability_Rpc{
						Rpc: &csi.NodeServiceCapability_RPC{
							Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
						},
					},
				},
			},
			subCap: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
			result: true,
		},
		{
			name: "Node Caps: STAGE_UNSTAGE, ",
			nodeCaps: []*csi.NodeServiceCapability{
				{
					Type: &csi.NodeServiceCapability_Rpc{
						Rpc: &csi.NodeServiceCapability_RPC{
							Type: csi.NodeServiceCapability_RPC_UNKNOWN,
						},
					},
				},
			},
			subCap: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
			result: false,
		},
	}
	for _, v := range tests {
		res := ContainsNodeServiceCapability(v.nodeCaps, v.subCap)
		if res != v.result {
			t.Errorf("test %s: expect %t, but result was %t", v.name, v.result, res)
		}
	}
}

func TestGbToByte(t *testing.T) {
	testcases := []struct {
		name string
		gb   int
		byte int64
	}{
		{"-1Gb", -1, 0},
		{"0Gb", 0, 0},
		{"1GB", 1, gib},
		{"10GB", 10, 10 * gib},
		{"100GB", 100, 100 * gib},
		{"1000GB", 1000, 1000 * gib},
	}

	for _, v := range testcases {
		res := GbToByte(v.gb)
		if res != v.byte {
			t.Errorf("test %s: expect %d, but result was %d", v.name, v.byte, res)
		}
	}
}

func TestByteCeilToGb(t *testing.T) {
	testcases := []struct {
		name string
		byte int64
		gb   int
	}{
		{"-1 Byte", -1, 0},
		{"0 Byte", 0, 0},
		{"1 Byte", 1, 1},
		{"1 Gib - 1 Byte", gib - 1, 1},
		{"1 Gib + 1 Byte", gib + 1, 2},
		{"10 Gib - 1 Byte", 10*gib - 1, 10},
		{"10 Gib", 10 * gib, 10},
		{"10 Gib + 1024 Byte", 10*gib + kib, 11},
		{"99 Gib - 1 Mib", 99*gib - mib, 99},
		{"99 Gib + 1 Mib", 99*gib + mib, 100},
	}

	for _, v := range testcases {
		res := ByteCeilToGb(v.byte)
		if res != v.gb {
			t.Errorf("test %s: expect %d Gb, but actually %d", v.name, v.gb, res)
		}
	}
}
