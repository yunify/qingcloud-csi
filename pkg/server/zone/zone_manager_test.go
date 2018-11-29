package zone

import (
	"os"
	"runtime"
	"testing"
)

var getzm = func() ZoneManager {
	// get storage class
	var filePath string
	if runtime.GOOS == "linux" {
		filePath = os.Getenv("GOPATH") + "/src/github.com/yunify/qingcloud-csi/deploy/block/kubernetes/config.yaml"
	}
	if runtime.GOOS == "darwin" {
		filePath = os.Getenv("GOPATH") + "/src/github.com/yunify/qingcloud-csi/deploy/block/kubernetes/config.yaml"
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
