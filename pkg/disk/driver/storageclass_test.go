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

package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"reflect"
	"strconv"
	"testing"
)

func TestNewDefaultQingStorageClassFromType(t *testing.T) {
	tests := []struct {
		name     string
		diskType VolumeType
		sc       *QingStorageClass
	}{
		{
			name:     "normal",
			diskType: DefaultVolumeType,
			sc: &QingStorageClass{
				diskType: DefaultVolumeType,
				maxSize:  VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:  VolumeTypeToMinSize[DefaultVolumeType],
				stepSize: VolumeTypeToStepSize[DefaultVolumeType],
				fsType:   common.DefaultFileSystem,
				replica:  DefaultDiskReplicaType,
			},
		},
		{
			name:     "number",
			diskType: 2,
			sc: &QingStorageClass{
				diskType: HighCapacityVolumeType,
				maxSize:  VolumeTypeToMaxSize[HighCapacityVolumeType],
				minSize:  VolumeTypeToMinSize[HighCapacityVolumeType],
				stepSize: VolumeTypeToStepSize[HighCapacityVolumeType],
				fsType:   common.DefaultFileSystem,
				replica:  DefaultDiskReplicaType,
			},
		},
		{
			name:     "invalid volume type",
			diskType: 99,
			sc:       nil,
		},
	}
	for _, test := range tests {
		res := NewDefaultQingStorageClassFromType(test.diskType)
		if !reflect.DeepEqual(test.sc, res) {
			t.Errorf("name %s: expect %v, but actually %v", test.name, test.sc, res)
		}
	}
}

func TestNewQingStorageClassFromMap(t *testing.T) {
	tests := []struct {
		name    string
		opt     map[string]string
		sc      *QingStorageClass
		isError bool
	}{
		{
			name: "normal",
			opt: map[string]string{
				StorageClassTypeName:     strconv.Itoa(int(DefaultVolumeType)),
				StorageClassMaxSizeName:  strconv.Itoa(VolumeTypeToMaxSize[DefaultVolumeType]),
				StorageClassMinSizeName:  strconv.Itoa(VolumeTypeToMinSize[DefaultVolumeType]),
				StorageClassStepSizeName: strconv.Itoa(VolumeTypeToStepSize[DefaultVolumeType]),
				StorageClassFsTypeName:   common.FileSystemExt3,
				StorageClassReplicaName:  strconv.Itoa(DiskSingleReplicaType),
			},
			sc: &QingStorageClass{
				diskType: DefaultVolumeType,
				maxSize:  VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:  VolumeTypeToMinSize[DefaultVolumeType],
				stepSize: VolumeTypeToStepSize[DefaultVolumeType],
				fsType:   common.FileSystemExt3,
				replica:  DiskSingleReplicaType,
			},
			isError: false,
		},
		{
			name: "only type",
			opt: map[string]string{
				StorageClassTypeName: strconv.Itoa(int(DefaultVolumeType)),
			},
			sc: &QingStorageClass{
				diskType: DefaultVolumeType,
				maxSize:  VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:  VolumeTypeToMinSize[DefaultVolumeType],
				stepSize: VolumeTypeToStepSize[DefaultVolumeType],
				fsType:   common.FileSystemExt4,
				replica:  DiskMultiReplicaType,
			},
			isError: false,
		},
		{
			name: "specific parameter",
			opt: map[string]string{
				StorageClassTypeName:     "5",
				StorageClassMaxSizeName:  "23",
				StorageClassMinSizeName:  "22",
				StorageClassStepSizeName: "4",
				StorageClassFsTypeName:   common.FileSystemExt3,
				StorageClassReplicaName:  strconv.Itoa(DiskSingleReplicaType),
			},
			sc: &QingStorageClass{
				diskType: 5,
				maxSize:  23,
				minSize:  22,
				stepSize: 4,
				fsType:   common.FileSystemExt3,
				replica:  DiskSingleReplicaType,
			},
			isError: false,
		},
		{
			name: "specific invalid type",
			opt: map[string]string{
				StorageClassTypeName:     "4",
				StorageClassMaxSizeName:  "23",
				StorageClassMinSizeName:  "22",
				StorageClassStepSizeName: "4",
				StorageClassFsTypeName:   common.FileSystemExt3,
				StorageClassReplicaName:  strconv.Itoa(DiskSingleReplicaType),
			},
			sc:      nil,
			isError: true,
		},
		{
			name: "empty tag",
			opt: map[string]string{
				StorageClassTypeName:     strconv.Itoa(int(DefaultVolumeType)),
				StorageClassMaxSizeName:  strconv.Itoa(VolumeTypeToMaxSize[DefaultVolumeType]),
				StorageClassMinSizeName:  strconv.Itoa(VolumeTypeToMinSize[DefaultVolumeType]),
				StorageClassStepSizeName: strconv.Itoa(VolumeTypeToStepSize[DefaultVolumeType]),
				StorageClassFsTypeName:   common.FileSystemExt3,
				StorageClassReplicaName:  strconv.Itoa(DiskSingleReplicaType),
				StorageClassTagsName:     "",
			},
			sc: &QingStorageClass{
				diskType: DefaultVolumeType,
				maxSize:  VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:  VolumeTypeToMinSize[DefaultVolumeType],
				stepSize: VolumeTypeToStepSize[DefaultVolumeType],
				fsType:   common.FileSystemExt3,
				replica:  DiskSingleReplicaType,
			},
			isError: false,
		},
		{
			name: "one tag",
			opt: map[string]string{
				StorageClassTypeName:     strconv.Itoa(int(DefaultVolumeType)),
				StorageClassMaxSizeName:  strconv.Itoa(VolumeTypeToMaxSize[DefaultVolumeType]),
				StorageClassMinSizeName:  strconv.Itoa(VolumeTypeToMinSize[DefaultVolumeType]),
				StorageClassStepSizeName: strconv.Itoa(VolumeTypeToStepSize[DefaultVolumeType]),
				StorageClassFsTypeName:   common.FileSystemExt3,
				StorageClassReplicaName:  strconv.Itoa(DiskSingleReplicaType),
				StorageClassTagsName:     "tag-12345567",
			},
			sc: &QingStorageClass{
				diskType: DefaultVolumeType,
				maxSize:  VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:  VolumeTypeToMinSize[DefaultVolumeType],
				stepSize: VolumeTypeToStepSize[DefaultVolumeType],
				fsType:   common.FileSystemExt3,
				replica:  DiskSingleReplicaType,
				tags:     []string{"tag-12345567"},
			},
			isError: false,
		},
		{
			name: "multiple tags",
			opt: map[string]string{
				StorageClassTypeName:     strconv.Itoa(int(DefaultVolumeType)),
				StorageClassMaxSizeName:  strconv.Itoa(VolumeTypeToMaxSize[DefaultVolumeType]),
				StorageClassMinSizeName:  strconv.Itoa(VolumeTypeToMinSize[DefaultVolumeType]),
				StorageClassStepSizeName: strconv.Itoa(VolumeTypeToStepSize[DefaultVolumeType]),
				StorageClassFsTypeName:   common.FileSystemExt3,
				StorageClassReplicaName:  strconv.Itoa(DiskSingleReplicaType),
				StorageClassTagsName:     "tag-12345567,tag-22345567,  tag-32345567 ",
			},
			sc: &QingStorageClass{
				diskType: DefaultVolumeType,
				maxSize:  VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:  VolumeTypeToMinSize[DefaultVolumeType],
				stepSize: VolumeTypeToStepSize[DefaultVolumeType],
				fsType:   common.FileSystemExt3,
				replica:  DiskSingleReplicaType,
				tags:     []string{"tag-12345567", "tag-22345567", "tag-32345567"},
			},
			isError: false,
		},
		{
			name: "without replica",
			opt: map[string]string{
				StorageClassTypeName:     strconv.Itoa(int(DefaultVolumeType)),
				StorageClassMaxSizeName:  strconv.Itoa(VolumeTypeToMaxSize[DefaultVolumeType]),
				StorageClassMinSizeName:  strconv.Itoa(VolumeTypeToMinSize[DefaultVolumeType]),
				StorageClassStepSizeName: strconv.Itoa(VolumeTypeToStepSize[DefaultVolumeType]),
				StorageClassFsTypeName:   common.FileSystemExt3,
			},
			sc: &QingStorageClass{
				diskType: DefaultVolumeType,
				maxSize:  VolumeTypeToMaxSize[DefaultVolumeType],
				minSize:  VolumeTypeToMinSize[DefaultVolumeType],
				stepSize: VolumeTypeToStepSize[DefaultVolumeType],
				fsType:   common.FileSystemExt3,
				replica:  DiskMultiReplicaType,
			},
			isError: false,
		},
	}
	for _, test := range tests {
		res, err := NewQingStorageClassFromMap(test.opt)
		if (err != nil) != test.isError {
			t.Errorf("name %s: expect %t, but actually %t", test.name, test.isError, err != nil)
		}
		if !reflect.DeepEqual(test.sc, res) {
			t.Errorf("name %s: expect %v, but actually %v", test.name, test.sc, res)
		}
	}
}

func TestQingStorageClass_FormatVolumeSizeByte(t *testing.T) {
	sc := NewDefaultQingStorageClassFromType(SSDEnterpriseVolumeType)
	tests := []struct {
		name       string
		inputSize  int64
		formatSize int64
	}{
		{
			name:       "normal",
			inputSize:  123 * common.Gib,
			formatSize: 130 * common.Gib,
		},
		{
			name:       "ceil",
			inputSize:  2001 * common.Gib,
			formatSize: 2000 * common.Gib,
		},
		{
			name:       "minus",
			inputSize:  -1 * common.Gib,
			formatSize: 10 * common.Gib,
		},
	}
	for _, test := range tests {
		res := sc.FormatVolumeSizeByte(test.inputSize)
		if test.formatSize != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.formatSize, res)
		}
	}
}

func TestQingStorageClass_GetRequiredVolumeSizeByte(t *testing.T) {
	sc := NewDefaultQingStorageClassFromType(SSDEnterpriseVolumeType)
	tests := []struct {
		name      string
		capRrange *csi.CapacityRange
		size      int64
		isError   bool
	}{
		{
			name: "normal",
			capRrange: &csi.CapacityRange{
				RequiredBytes: 20 * common.Gib,
				LimitBytes:    20 * common.Gib,
			},
			size:    20 * common.Gib,
			isError: false,
		},
		{
			name: "without limit",
			capRrange: &csi.CapacityRange{
				RequiredBytes: 23 * common.Gib,
			},
			size:    30 * common.Gib,
			isError: false,
		},
		{
			name: "failed",
			capRrange: &csi.CapacityRange{
				RequiredBytes: 2 * common.Gib,
				LimitBytes:    4 * common.Gib,
			},
			size:    -1,
			isError: true,
		},
	}
	for _, test := range tests {
		res, err := sc.GetRequiredVolumeSizeByte(test.capRrange)
		if (err != nil) != test.isError {
			t.Errorf("name %s: expect %t, but actually %t", test.name, test.isError, err != nil)
		}
		if test.size != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.size, res)
		}
	}
}
