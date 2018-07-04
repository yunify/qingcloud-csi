package block

import (
	"testing"
	"os"
	"google.golang.org/grpc/codes"
	"k8s.io/kubernetes/pkg/util/mount"
	"syscall"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)

var targetPath = "/root/adf"

func isNotDirErr(err error) bool {
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOTDIR {
		return true
	}
	return false
}

func TestBindMount(t *testing.T){
	// 1. Mount
	// check targetPath is mounted
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)
	//  notMnt, err := mount.New("").IsNotMountPoint(targetPath)
	flag := isNotDirErr(err)
	t.Logf("%v |%v |%v", notMnt, err, flag)

	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(targetPath, 0750); err != nil {
				t.Error(err.Error())
			}
			notMnt = true
		} else {
			t.Error(codes.Internal, err.Error())
		}
	}
	if !notMnt {
		t.Logf("%s %v", targetPath, notMnt)
	}
}

func TestGbToByte(t *testing.T){
	testcases :=[]struct{
		gb int
		byte int64
		res bool
	}{
		{10, 10*gib, true},
		{-1, 0, true},
		{1, gib,true},
	}

	for _, v:=range testcases{
		res := GbToByte(v.gb)
		if res == v.byte{
			t.Logf("pass Gib %d, Byte=%d", v.gb, res)
		}else{
			t.Errorf("faile Gib %d, expect Byte %d, but actually %d", v.gb, v.byte, res)
		}
	}
}

func TestByteCeilToGb(t *testing.T){
	testcases :=[]struct{
		gb int
		byte int64
		res bool
	}{
		{10, 10*gib- 2, true},
		{0, -1, true},
		{1, gib,true},
	}

	for _, v:=range testcases{
		res := ByteCeilToGb(v.byte)
		if res == v.gb{
			t.Logf("pass Gib %d, Byte %d", v.gb, res)
		}else{
			t.Errorf("faile Byte %d, expect Gb %d, but actually %d", v.byte, v.gb, res)
		}
	}
}

func TestHasSameAccessMode(t *testing.T){
	testcases := []struct{
		access []*csi.VolumeCapability_AccessMode
		cap []*csi.VolumeCapability
		res bool
	}{
		{
			[]*csi.VolumeCapability_AccessMode{&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
			[]*csi.VolumeCapability{
				{nil, &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}}},
				true,
		},
		{
			[]*csi.VolumeCapability_AccessMode{&csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER}},
			[]*csi.VolumeCapability{
				{nil, &csi.VolumeCapability_AccessMode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER}}},
			false,
		},

	}
	for _, v:=range testcases{
		res := HasSameAccessMode(v.access,v.cap)
		if res == v.res{
			t.Logf("success")
		}else{
			t.Errorf("failed, expect %t, but actually %t", v.res, res)
		}
	}
}