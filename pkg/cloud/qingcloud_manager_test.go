/*
Copyright (C) 2018 Yunify, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this work except in compliance with the License.
You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud

import (
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"k8s.io/klog"
	"os"
	"path"
	"testing"
)

const (
	findSnapshotId = "ss-yv7iaiqw"
	findVolumeId   = "vol-b81mirdr"
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

const (
	tagId                 = "tag-qaf8td3d"
	tagId2                = "tag-1phpmfym"
	tagIdInOtherZone      = "tag-glozcqzd"
	resourceId            = "vol-zqxq0i9k"
	resourceType          = "volume"
	resourceIdInOtherZone = "vol-t6jgk2fp"
)

func TestQingCloudManager_FindTag(t *testing.T) {
	tests := []struct {
		name     string
		tagId    string
		foundTag bool
		isError  bool
	}{
		{
			name:     "valid tag",
			tagId:    tagId,
			foundTag: true,
			isError:  false,
		},
		{
			name:     "other zone",
			tagId:    tagIdInOtherZone,
			foundTag: false,
			isError:  false,
		},
	}
	for _, v := range tests {
		tagInfo, err := cfg.FindTag(v.tagId)
		if (tagInfo != nil) != v.foundTag && (err == nil) != v.isError {
			t.Errorf("name %s, expect [%t,%t], but actually [%t,%t]", v.name, v.foundTag, v.isError,
				tagInfo != nil, err == nil)
		}
	}
}

func TestQingCloudManager_IsValidTags(t *testing.T) {
	tests := []struct {
		name    string
		tagId   []string
		isValid bool
	}{
		{
			name:    "multiple tags",
			tagId:   []string{tagId, tagId2},
			isValid: true,
		},
		{
			name:    "single tags",
			tagId:   []string{tagId},
			isValid: true,
		},
		{
			name:    "tags in other zone",
			tagId:   []string{tagIdInOtherZone},
			isValid: false,
		},
	}
	for _, test := range tests {
		res := cfg.IsValidTags(test.tagId)
		if test.isValid != res {
			t.Errorf("name %s, expect %t, but actually %t", test.name, test.isValid, res)
		}
	}
}

func TestQingCloudManager_AttachTags(t *testing.T) {
	tests := []struct {
		name         string
		tagId        []string
		resourceId   string
		resourceType string
		isError      bool
	}{
		{
			name:         "add multiple tags",
			tagId:        []string{tagId, tagId2},
			resourceId:   resourceId,
			resourceType: ResourceTypeVolume,
			isError:      false,
		},
		{
			name:         "re-attach tags",
			tagId:        []string{tagId, tagId2},
			resourceId:   resourceId,
			resourceType: ResourceTypeVolume,
			isError:      false,
		},
		{
			name:         "attach other zone resource",
			tagId:        []string{tagId},
			resourceId:   resourceIdInOtherZone,
			resourceType: ResourceTypeVolume,
			isError:      true,
		},
		{
			name:         "invalid resource type",
			tagId:        []string{tagId},
			resourceId:   resourceId,
			resourceType: ResourceTypeSnapshot,
			isError:      true,
		},
		{
			name:         "empty tags slice",
			tagId:        []string{},
			resourceId:   resourceId,
			resourceType: ResourceTypeVolume,
			isError:      false,
		},
		{
			name:         "nil tags slice",
			tagId:        nil,
			resourceId:   resourceId,
			resourceType: ResourceTypeVolume,
			isError:      false,
		},
	}
	for _, v := range tests {
		err := cfg.AttachTags(v.tagId, v.resourceId, v.resourceType)
		if (err != nil) != v.isError {
			t.Errorf("name %s, expect %t, but actually %t/%s", v.name, v.isError, err != nil, err.Error())
		}
	}
}

func TestQingCloudManager_CloneVolume(t *testing.T) {
	tests := []struct {
		name     string
		volName  string
		volType  int
		srcVolId string
		zone     string
		isError  bool
	}{
		{
			name:     "normal",
			volName:  "clone-test",
			volType:  driver.StandardVolumeType.Int(),
			srcVolId: findVolumeId,
			zone:     zone,
			isError:  false,
		},
	}
	for _, test := range tests {
		_, err := cfg.CloneVolume(test.volName, test.volType, test.srcVolId, test.zone)
		if (err != nil) != test.isError {
			t.Errorf("name %s: expect %t, but actually %t/%s", test.name, test.isError, err != nil, err.Error())
		}
	}
}
