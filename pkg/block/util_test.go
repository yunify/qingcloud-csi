package block

import (
	"testing"
	"os"
	"google.golang.org/grpc/codes"
	"k8s.io/kubernetes/pkg/util/mount"
	"syscall"
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
