package driver_test

import (
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"testing"
)

func TestVolumeType_String(t *testing.T) {
	tests := []struct {
		volType driver.VolumeType
		volName string
	}{
		{
			volType: driver.VolumeType(2),
			volName: "HighCapacity",
		},
		{
			volType: driver.VolumeType(5),
			volName: "NeonSAN",
		},
		{
			volType: driver.VolumeType(-1),
			volName: "",
		},
	}
	for _, v := range tests {
		if v.volType.String() != v.volName {
			t.Errorf("VolumeType(%d).String() expect %s, but actually %s", v.volType.Int(), v.volName, v.volType.String())
		}
	}
}
