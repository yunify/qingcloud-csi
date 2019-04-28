package zone

import (
	"os"
	"path"
	"runtime"
	"testing"
)

var getzm = func() ZoneManager {
	// get storage class
	var filePath string
	if runtime.GOOS == "linux" {
		filePath = path.Join(os.Getenv("GOPATH"), "src/github.com/yunify/qingcloud-csi/deploy/disk/kubernetes/config.yaml")
	}
	if runtime.GOOS == "darwin" {
		filePath = path.Join(os.Getenv("GOPATH"), "src/github.com/yunify/qingcloud-csi/deploy/disk/kubernetes/config.yaml")
	}
	vm, err := NewZoneManagerFromFile(filePath)
	if err != nil {
		return nil
	}
	return vm
}

func TestGetZoneList(t *testing.T) {
	zm := getzm()
	// testcase
	testcases := []struct {
		name string
	}{
		{
			name: "Get zone list",
		},
	}

	// test findVolume
	for _, v := range testcases {
		zones, _ := zm.GetZoneList()
		if len(zones) <= 0 {
			t.Errorf("name %s: expected get at least one active zone, but actually [%d] active zones", v.name,
				len(zones))
		}
	}
}
