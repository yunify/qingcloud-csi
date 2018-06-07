package block

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"testing"
)

var getvp = func() *volumeProvisioner {
	// get storage class
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = "C:\\Users\\wangx\\Documents\\config.json"
	}
	if runtime.GOOS == "linux" {
		filepath = "/root/config.json"
	}
	if runtime.GOOS == "darwin" {
		filepath = "./config.json"
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Errorf("Open file error: %s", err.Error())
		os.Exit(-1)
	}
	sc := qingStorageClass{}
	err = json.Unmarshal(content, &sc)
	if err != nil {
		fmt.Errorf("get storage class error: %s", err.Error())
		os.Exit(-1)
	}

	// get volume provisioner
	vp, err := newVolumeProvisioner(&sc)
	if err != nil {
		fmt.Errorf("new volume provisioner error: %s", err.Error())
		os.Exit(-1)
	}
	return vp
}

func TestFindVolume(t *testing.T) {
	// testcase
	testcase := []struct {
		id    string
		exist bool
	}{
		{"vol-fhlkhxpr", true},
		{"vol-vol-fhlkhxpw", false},
	}

	vp := getvp()
	// test findVolume
	for _, v := range testcase {
		flag, err := vp.findVolume(v.id)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		if (flag != nil) == v.exist {
			t.Logf("volume id %s, expect %t, actually %t", v.id, v.exist, flag != nil)
		} else {
			t.Errorf("volume id %s, expect %t, actually %t", v.id, v.exist, flag != nil)
		}
	}
}

func TestFindVolumeByName(t *testing.T) {
	testcase := []struct {
		name  string
		exist bool
	}{
		{"hp-test", true},
		{"hp-test-false", false},
	}

	vp := getvp()
	// test findVolume
	for _, v := range testcase {
		flag, err := vp.findVolumeByName(v.name)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		if (flag != nil) == v.exist {
			t.Logf("volume id %s, expect %t, actually %t", v.name, v.exist, flag != nil)
		} else {
			t.Errorf("volume id %s, expect %t, actually %t", v.name, v.exist, flag != nil)
		}
	}
}

func TestCreateVolume(t *testing.T) {
	// testcase
	testcase := []struct {
		vc            blockVolume
		createSuccess bool
	}{
		{blockVolume{VolName: "pvc-test-", VolSize: 12}, true},
		{blockVolume{VolName: "pvc-test-", VolSize: 121}, true},
		{blockVolume{VolName: "pvc-test-", VolSize: -1}, true},
	}
	vp := getvp()
	for i, v := range testcase {
		v.vc.VolName += strconv.Itoa(i)
		err := vp.CreateVolume(testcase[i].vc.VolSize, &v.vc)
		if (err == nil) == v.createSuccess {
			t.Logf("testcase passed, %v", v)
		} else {
			t.Error(err)
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	vp := getvp()
	volumeID := "vol-vhvb1il0"
	err := vp.DeleteVolume(volumeID)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("testcase delete %s success", volumeID)
	}
}

func TestAttachVolume(t *testing.T) {
	vp := getvp()
	volumeID := "vol-fhlkhxpr"
	instanceID := "i-msu2th7i"
	err := vp.AttachVolume(&volumeID, &instanceID)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("testcase attach volume %s to instance %s success",
			volumeID, instanceID)
	}
}

func TestDetachVolume(t *testing.T) {
	vp := getvp()
	volumeID := "vol-fhlkhxpr"
	instanceID := "i-msu2th7i"
	err := vp.DetachVolume(&volumeID, &instanceID)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("testcase detach volume %s from instance %s success",
			volumeID, instanceID)
	}
}
