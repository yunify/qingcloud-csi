package zone

import (
	"testing"
)

var getzm = func() ZoneManager {
	// get storage class
	filePath := "/root/.qingcloud/config.yaml"
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
