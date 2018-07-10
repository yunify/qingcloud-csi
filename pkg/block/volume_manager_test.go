package block

import (
	"runtime"
	"testing"
)

var (
	volumeId1   string = "vol-5pmaukiv"
	volumeName1 string = "qingcloud-csi-test"
	instanceId1 string = "i-msu2th7i"
	instanceId2 string = "i-hgz8mri2"
)

var getvm = func() VolumeManager {
	// get storage class
	var filepath string
	if runtime.GOOS == "linux" {
		filepath = "../../ut-config.yaml"
	}
	if runtime.GOOS == "darwin" {
		filepath = "../../ut-config.yaml"
	}
	qcConfig, err := ReadConfigFromFile(filepath)
	if err != nil {
		return nil
	}
	vm, err := NewVolumeManagerWithConfig(qcConfig)
	if err != nil {
		return nil
	}

	return vm
}

func TestFindVolume(t *testing.T) {
	vm := getvm()
	_, err := vm.FindVolume(volumeId1)
	if err != nil {
		t.Error(err.Error())
	}
	// testcase
	testcase := []struct {
		name   string
		id     string
		result bool
	}{
		{
			name:   "Avaiable",
			id:     volumeId1,
			result: true,
		},
		{
			name:   "Not found",
			id:     volumeId1 + "fake",
			result: false,
		},
		{
			name:   "By name",
			id:     volumeName1,
			result: false,
		},
	}

	// test findVolume
	for _, v := range testcase {
		vol, err := vm.FindVolume(v.id)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		res := (vol != nil)
		if res != v.result {
			t.Errorf("name: %s, expect %t, actually %t", v.name, v.result, res)
		}
	}
}

func TestFindVolumeByName(t *testing.T) {
	testcase := []struct {
		name     string
		testname string
		result   bool
	}{
		{
			name:     "Avaiable",
			testname: volumeName1,
			result:   true,
		},
		{
			name:     "Ceased",
			testname: "sanity",
			result:   false,
		},
		{
			name:     "Volume id",
			testname: volumeId1,
			result:   false,
		},
		{
			name:     "Substring",
			testname: string((volumeName1)[:2]),
			result:   false,
		},
		{
			name:     "Null string",
			testname: "",
			result:   false,
		},
	}

	vm := getvm()
	// test findVolume
	for _, v := range testcase {
		vol, err := vm.FindVolumeByName(v.testname)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		res := (vol != nil)
		if res != v.result {
			t.Errorf("name %s, expect %t, actually %t", v.name, v.testname, res)
		}
	}
}

func TestCreateVolume(t *testing.T) {

	sc := NewDefaultQingStorageClass()
	vm := getvm()

	testcases := []struct {
		name    string
		volName string
		reqSize int
		result  bool
		volId   string
	}{
		{
			name:    "create volume name test-1",
			volName: "test-1",
			reqSize: 1,
			result:  true,
			volId:   "",
		},
		{
			name:    "create volume name test-1 repeatedly",
			volName: "test-1",
			reqSize: 3,
			result:  false,
			volId:   "",
		},
		{
			name:    "create volume name test-2",
			volName: "test-2",
			reqSize: 20,
			result:  true,
			volId:   "",
		},
	}
	for i, v := range testcases {
		volId, err := vm.CreateVolume(v.volName, v.reqSize, *sc)
		if err != nil {
			t.Errorf("test %s: %s", v.name, err.Error())
		} else {
			testcases[i].volId = volId
			vol, _ := vm.FindVolume(volId)
			if *vol.VolumeName != v.volName {
				t.Errorf("test %s: expect %t", v.name, v.result)
			}
		}
	}
	// clear process
	for _, v := range testcases {
		err := vm.DeleteVolume(v.volId)
		if err != nil {
			t.Errorf("test %s: delete error %s", v.name, err.Error())
		}
	}
}

func TestAttachVolume(t *testing.T) {
	vm := getvm()
	// testcase
	testcases := []struct {
		name       string
		volumeId   string
		instanceId string
		isError    bool
	}{
		{
			name:       "Attach success",
			volumeId:   volumeId1,
			instanceId: instanceId1,
			isError:    false,
		},
		{
			name:       "Attach repeatedly, idempotent",
			volumeId:   volumeId1,
			instanceId: instanceId1,
			isError:    false,
		},
		{
			name:       "Attach another instance",
			volumeId:   volumeId1,
			instanceId: instanceId2,
			isError:    true,
		},
		{
			name:       "Attach not exist instance",
			volumeId:   volumeId1,
			instanceId: "ins-123456",
			isError:    true,
		},
	}
	for _, v := range testcases {
		err := vm.AttachVolume(v.volumeId, v.instanceId)
		if err != nil && !v.isError {
			t.Errorf("error name %s: %s", v.name, err.Error())
		}
	}

}

func TestIsAttachedToInstance(t *testing.T) {
	vm := getvm()
	testcases := []struct {
		name       string
		volumeId   string
		instanceId string
		result     bool
		isError    bool
	}{
		{
			name:       "Attach success",
			volumeId:   volumeId1,
			instanceId: instanceId1,
			result:     true,
			isError:    false,
		},
		{
			name:       "Attach another instance",
			volumeId:   volumeId1,
			instanceId: instanceId2,
			result:     false,
			isError:    false,
		},
		{
			name:       "Not found volume",
			volumeId:   volumeId1 + "fake",
			instanceId: instanceId1,
			result:     false,
			isError:    true,
		},
		{
			name:       "Not found instance",
			volumeId:   volumeId1,
			instanceId: instanceId1 + "fake",
			result:     false,
			isError:    false,
		},
	}
	for _, v := range testcases {
		flag, err := vm.IsAttachedToInstance(v.volumeId, v.instanceId)
		if err != nil {
			if !v.isError {
				t.Errorf("error name %s: %s", v.name, err.Error())
			}
		}
		if flag != v.result {
			t.Errorf("name %s: expect %t", v.name, v.result)
		}
	}
}

func TestDetachVolume(t *testing.T) {
	vm := getvm()
	testcase := []struct {
		name       string
		volumeId   string
		instanceId string
		isError    bool
	}{
		{
			name:       "detach normally",
			volumeId:   volumeId1,
			instanceId: instanceId1,
			isError:    false,
		},
		{
			name:       "detach repeatedly, idempotent",
			volumeId:   volumeId1,
			instanceId: instanceId1,
			isError:    false,
		},
		{
			name:       "volume not found",
			volumeId:   "fake",
			instanceId: instanceId1,
			isError:    true,
		},
		{
			name:       "instance not found",
			volumeId:   volumeId1,
			instanceId: "fake",
			isError:    true,
		},
	}

	for _, v := range testcase {
		err := vm.DetachVolume(v.volumeId, v.instanceId)
		if err != nil && !v.isError {
			t.Errorf("error name %s: %s", v.name, err.Error())
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	vm := getvm()
	// testcase
	testcase := []struct {
		name    string
		id      string
		isError bool
	}{
		{
			name:    "delete first volume",
			id:      volumeId1,
			isError: false,
		},
		{
			name:    "delete first volume repeatedly",
			id:      volumeId1,
			isError: true,
		},
		{
			name:    "delete not exist volume",
			id:      "vol-1234567",
			isError: true,
		},
	}
	for _, v := range testcase {
		err := vm.DeleteVolume(v.id)
		if err != nil && !v.isError {
			t.Errorf("error name %s: %s", v.name, err.Error())
		}
	}
}
