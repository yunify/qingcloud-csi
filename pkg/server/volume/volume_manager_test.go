// +-------------------------------------------------------------------------
// | Copyright (C) 2018 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package volume

import (
	"github.com/yunify/qingcloud-csi/pkg/server"
	"github.com/yunify/qingcloud-csi/pkg/server/storageclass"
	"runtime"
	"testing"
)

var (
	// Tester should set these variables before executing unit test.
	volumeId1       string = "vol-8boq0cz6"
	volumeName1     string = "qingcloud-csi-test"
	instanceId1     string = "i-0nuxqgal"
	instanceId2     string = "i-tta11nep"
	resizeVolumeId  string = "vol-tysu3tg2"
	volFromSnapNorm string = "volFromSnapNorm"
	snapshotId1     string = "ss-dmrxy2mn"
)

var getvm = func() VolumeManager {
	// get storage class
	var filePath string
	if runtime.GOOS == "linux" {
		filePath = "/root/.qingcloud/config.yaml"
	}
	if runtime.GOOS == "darwin" {
		filePath = "/etc/qingcloud/client.yaml"
	}
	vm, err := NewVolumeManagerFromFile(filePath)
	if err != nil {
		return nil
	}
	return vm
}

func TestFindVolume(t *testing.T) {
	vm := getvm()
	// testcase
	testcases := []struct {
		name   string
		id     string
		result bool
	}{
		{
			name:   "Available",
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
	for _, v := range testcases {
		vol, err := vm.FindVolume(v.id)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		res := vol != nil
		if res != v.result {
			t.Errorf("name: %s, expect %t, actually %t", v.name, v.result, res)
		}
	}
}

func TestFindVolumeByName(t *testing.T) {
	testcases := []struct {
		name   string
		volume string
		result bool
	}{
		{
			name:   "Available",
			volume: volumeName1,
			result: true,
		},
		{
			name:   "Ceased",
			volume: "sanity",
			result: false,
		},
		{
			name:   "Volume id",
			volume: volumeId1,
			result: false,
		},
		{
			name:   "Substring",
			volume: string((volumeName1)[:2]),
			result: false,
		},
		{
			name:   "Null string",
			volume: "",
			result: false,
		},
	}

	vm := getvm()
	// test findVolume
	for _, v := range testcases {
		vol, err := vm.FindVolumeByName(v.volume)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		res := vol != nil
		if res != v.result {
			t.Errorf("name %s, expect %t, actually %t", v.name, v.result, res)
		}
	}
}

func TestCreateVolume(t *testing.T) {
	sc := storageclass.NewDefaultQingStorageClass()
	vm := getvm()

	testcases := []struct {
		name         string
		volName      string
		reqSize      int
		storageClass storageclass.QingStorageClass
		result       bool
		volId        string
	}{
		{
			name:         "create volume name test-1",
			volName:      "test-1",
			reqSize:      1,
			storageClass: *sc,
			result:       true,
			volId:        "",
		},
		{
			name:         "create volume name test-1 repeatedly",
			volName:      "test-1",
			reqSize:      3,
			storageClass: *sc,
			result:       true,
			volId:        "",
		},
		{
			name:         "create volume name test-2",
			volName:      "test-2",
			reqSize:      20,
			storageClass: *sc,
			result:       true,
			volId:        "",
		},

		{
			name:    "create volume name test-3 for single replica",
			volName: "test-3",
			reqSize: 20,
			storageClass: storageclass.QingStorageClass{
				VolumeType:     100,
				VolumeMaxSize:  500,
				VolumeMinSize:  10,
				VolumeStepSize: 10,
				VolumeFsType:   server.FileSystemDefault,
				VolumeReplica:  server.SingleReplica,
			},
			result: true,
			volId:  "",
		},
	}
	for i, v := range testcases {
		volId, err := vm.CreateVolume(v.volName, v.reqSize, v.storageClass)
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
	testcases := []struct {
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

	for _, v := range testcases {
		err := vm.DetachVolume(v.volumeId, v.instanceId)
		if err != nil && !v.isError {
			t.Errorf("error name %s: %s", v.name, err.Error())
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	vm := getvm()
	// testcase
	testcases := []struct {
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
	for _, v := range testcases {
		err := vm.DeleteVolume(v.id)
		if err != nil && !v.isError {
			t.Errorf("error name %s: %s", v.name, err.Error())
		}
	}
}

func TestResizeVolume(t *testing.T) {
	vm := getvm()
	// testcase
	testcases := []struct {
		name    string
		id      string
		size    int
		isError bool
	}{
		{
			name:    "resize normally",
			id:      resizeVolumeId,
			size:    30,
			isError: false,
		},
	}
	for _, v := range testcases {
		err := vm.ResizeVolume(v.id, v.size)
		if err != nil && !v.isError {
			t.Errorf("name %s: expect [%t] but actually [%s]", v.name, v.isError, err)
		}
	}
}

func TestCreateVolumeFromSnapshot(t *testing.T) {
	vm := getvm()
	testcases := []struct {
		name       string
		volumeName string
		snapshotId string
		isError    bool
	}{
		{
			name:       "create normally",
			volumeName: volFromSnapNorm,
			snapshotId: snapshotId1,
			isError:    false,
		},
		{
			name:       "zero value input",
			volumeName: volFromSnapNorm,
			snapshotId: "",
			isError:    true,
		},
	}
	for _, v := range testcases {
		_, err := vm.CreateVolumeFromSnapshot(v.volumeName, v.snapshotId)
		if (err != nil) != v.isError {
			t.Errorf("name %s: expect %t, but actually %s", v.name, v.isError, err)
		}
	}
}
