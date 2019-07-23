package cloud

import (
	"k8s.io/klog"
	"os"
	"path"
	"testing"
)

const (
	findSnapshotId = "ss-yv7iaiqw"
	findVolumeId   = "vol-ez60gw9f"
	findVolumeName = "er"
	zone           = "pek3d"
	deleteVolumeId = findVolumeId
	findInstId     = "i-5o5kwa53"
)

var cfg CloudManager

func init() {
	klog.InitFlags(nil)
	var err error
	cfg, err = NewQingCloudManagerFromFile(path.Join(os.Getenv("HOME"), ".qingcloud/config.yaml"))
	if err != nil {
		klog.Fatal(err.Error())
	}
}

func TestQingCloudManager_CreateVolume(t *testing.T) {
	tests := []struct {
		name    string
		volName string
		volSize int
		volRepl int
		volType int
		volZone string
		isError bool
	}{
		{
			name:    "Create in region's zone",
			volName: "csi-create-in-region",
			volSize: 10,
			volRepl: 1,
			volType: 0,
			volZone: "pek3c",
			isError: false,
		},
	}

	for _, test := range tests {
		volId, err := cfg.CreateVolume(test.volName, test.volSize, test.volRepl, test.volType, test.volZone)
		if err != nil {
			if !test.isError {
				t.Errorf("testcase %s: expect error %t, but actually error: %s", test.name, test.isError, err)
			}
		} else {
			if volId == "" {
				t.Errorf("testcase %s: cannot get volume %s id", test.name, test.volName)
			}
		}
	}
}

func TestQingCloudManager_FindVolume(t *testing.T) {
	tests := []struct {
		name    string
		volId   string
		volName string
		volZone string
		isError bool
	}{
		{
			name:    "Create in region's zone",
			volId:   findVolumeId,
			volName: findVolumeName,
			volZone: zone,
			isError: false,
		},
	}
	for _, test := range tests {
		volInfo, err := cfg.FindVolume(test.volId)
		if err != nil {
			if !test.isError {
				t.Errorf("testcase %s: expect error %t, but actually error: %s", test.name, test.isError, err)
			}
		} else {
			if volInfo == nil {
				t.Errorf("testcase %s: cannot get volume %s info", test.name, test.volId)
			}
			if *volInfo.VolumeName != test.volName {
				t.Errorf("testcase %s: expect volume name %s, but actually %s", test.name, test.volName,
					*volInfo.VolumeName)
			}
		}
	}
}

func TestQingCloudManager_FindVolumeByName(t *testing.T) {
	tests := []struct {
		name    string
		volId   string
		volName string
		volZone string
		isError bool
	}{
		{
			name:    "Create in region's zone",
			volId:   findVolumeId,
			volName: findVolumeName,
			volZone: zone,
			isError: false,
		},
	}
	for _, test := range tests {
		volInfo, err := cfg.FindVolumeByName(test.volName)
		if err != nil {
			if !test.isError {
				t.Errorf("testcase %s: expect error %t, but actually error: %s", test.name, test.isError, err)
			}
		} else {
			if volInfo == nil {
				t.Errorf("testcase %s: cannot get volume %s info", test.name, test.volName)
			}
			if *volInfo.VolumeID != test.volId {
				t.Errorf("testcase %s: expect volume name %s, but actually %s", test.name, test.volId,
					*volInfo.VolumeID)
			}
			if *volInfo.ZoneID != test.volZone {
				t.Errorf("testcase %s: expect volume zone %s, but actually %s", test.name, test.volZone,
					*volInfo.ZoneID)
			}
		}
	}
}

func TestQingCloudManager_FindInstance(t *testing.T) {
	tests := []struct {
		name     string
		instId   string
		instZone string
		isError  bool
	}{
		{
			name:     "Find in region's zone",
			instId:   findInstId,
			instZone: zone,
			isError:  false,
		},
	}
	for _, test := range tests {
		instInfo, err := cfg.FindInstance(test.instId)
		if err != nil {
			if !test.isError {
				t.Errorf("testcase %s: expect error %t, but actually error: %s", test.name, test.isError, err)
			}
		} else {
			if instInfo == nil {
				t.Errorf("testcase %s: cannot get instance %s info", test.name, test.instId)
			}
			if *instInfo.ZoneID != test.instZone {
				t.Errorf("testcase %s: expect zone name %s, but actually %s", test.name, test.instZone,
					*instInfo.ZoneID)
			}
		}
	}
}

func TestQingCloudManager_DeleteVolume(t *testing.T) {
	tests := []struct {
		name    string
		volid   string
		isError bool
	}{
		{
			name:    "Find in region's zone",
			volid:   findVolumeId,
			isError: false,
		},
	}
	for _, test := range tests {
		err := cfg.DeleteVolume(test.volid)
		if test.isError != (err != nil) {
			t.Errorf("testcase %s: expect error %t, but actually error: %s", test.name, test.isError, err)
		}
	}
}

func TestQingCloudManager_CreateVolumeFromSnapshot(t *testing.T) {
	tests := []struct {
		name    string
		volName string
		snapId  string
		zone    string
		isError bool
	}{
		{
			name:    "create in another zone of region",
			volName: "csi-another-zone",
			snapId:  findSnapshotId,
			zone:    zone,
			isError: false,
		},
	}
	for _, test := range tests {
		volId, err := cfg.CreateVolumeFromSnapshot(test.volName, test.snapId, test.zone)
		if err != nil {
			if !test.isError {
				t.Errorf("testcase %s: expect error %t, but actually error: %s", test.name, test.isError, err)
			}
		} else {
			if volId == "" {
				t.Errorf("testcase %s: cannot get volume %s id", test.name, test.volName)
			}
		}
	}
}
