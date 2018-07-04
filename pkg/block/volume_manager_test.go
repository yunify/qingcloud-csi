package block

import (
	"runtime"
	"strconv"
	"testing"
)

var getvp = func() *volumeManager {
	// get storage class
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = "C:\\Users\\wangx\\Documents\\config.json"
	}
	if runtime.GOOS == "linux" {
		filepath = "/root/config.json"
	}
	if runtime.GOOS == "darwin" {
		filepath = "./config.yaml"
	}
	qcConfig, err := ReadConfigFromFile(filepath)
	if err != nil {
		return nil
	}
	vm, err := NewVolumeManagerWithConfig(qcConfig)
	if err != nil{
		return nil
	}

	return vm
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
		flag, err := vp.FindVolume(v.id)
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
		{"sanity", false},
	}

	vp := getvp()
	// test findVolume
	for _, v := range testcase {
		flag, err := vp.FindVolumeByName(v.name)
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
	// storageclass
	sc := NewDefaultQingStorageClass()
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
		volId, err := vp.CreateVolume(v.vc.VolName, v.vc.VolSize, *sc)
		if (err == nil) == v.createSuccess {
			t.Logf("testcase passed, %s", volId)
		} else {
			t.Error(err)
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	vp := getvp()
	// testcase
	testcase := []struct {
		id           string
	}{
		{"vol-oaihhpgo"},
		{"vol-wmxjlndr"},
		{"vol-30ltz79j"},
	}
	for _, v:= range testcase{
		err := vp.DeleteVolume(v.id)
		if err != nil{
			t.Error(err)
		}else{
			t.Logf("testcase delete %s success", v.id)
		}
	}
}

func TestAttachVolume(t *testing.T) {
	vp := getvp()
	// testcase
	testcases := []struct {
		volumeId 	string
		instanceId	string
		result 		bool
	}{
		{"vol-fhlkhxpr", "i-msu2th7i", true},
	}
	for _, v:=range testcases{
		err := vp.AttachVolume(v.volumeId, v.instanceId)
		if err != nil {
			t.Error(err)
		} else {
			t.Logf("testcase attach volume %s to instance %s success", v.volumeId, v.instanceId)
		}
	}

}

func TestIsAttachedToInstance(t *testing.T) {
	vp := getvp()
	volumeID := "vol-fhlkhxpr"
	instanceID := "i-msu2th7i"
	flag, err := vp.IsAttachedToInstance(volumeID, instanceID)
	if err != nil {
		t.Error(err)
	} else {
		if flag == true {
			t.Logf("volume %s is attached to instance %s", volumeID, instanceID)
		} else {
			t.Errorf("volume %s is not attached to instance %s", volumeID, instanceID)
		}
	}
}

func TestDetachVolume(t *testing.T) {
	vp := getvp()
	volumeID := "vol-fhlkhxpr"
	instanceID := "i-msu2th7i"
	err := vp.DetachVolume(volumeID, instanceID)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("testcase detach volume %s from instance %s success",
			volumeID, instanceID)
	}
}
