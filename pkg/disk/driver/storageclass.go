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
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"strconv"
	"strings"
)

const (
	StorageClassTypeName     = "type"
	StorageClassMaxSizeName  = "maxSize"
	StorageClassMinSizeName  = "minSize"
	StorageClassStepSizeName = "stepSize"
	StorageClassFsTypeName   = "fsType"
	StorageClassReplicaName  = "replica"
	StorageClassTagsName     = "tags"
)

type QingStorageClass struct {
	diskType VolumeType
	maxSize  int
	minSize  int
	stepSize int
	fsType   string
	replica  int
	tags     []string
}

// NewDefaultQingStorageClassFromType create default qingStorageClass by specified volume type
func NewDefaultQingStorageClassFromType(diskType VolumeType) *QingStorageClass {
	if diskType.IsValid() != true {
		return nil
	}
	return &QingStorageClass{
		diskType: diskType,
		maxSize:  VolumeTypeToMaxSize[diskType],
		minSize:  VolumeTypeToMinSize[diskType],
		stepSize: VolumeTypeToStepSize[diskType],
		fsType:   common.DefaultFileSystem,
		replica:  DefaultDiskReplicaType,
	}
}

// NewQingStorageClassFromMap create qingStorageClass object from map
func NewQingStorageClassFromMap(opt map[string]string) (*QingStorageClass, error) {
	volType := -1
	maxSize, minSize, stepSize := -1, -1, -1
	fsType := ""
	replica := -1
	var tags []string
	for k, v := range opt {
		switch strings.ToLower(k) {
		case strings.ToLower(StorageClassTypeName):
			// Convert to integer
			iv, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			volType = iv
		case strings.ToLower(StorageClassMaxSizeName):
			// Convert to integer
			iv, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			maxSize = iv
		case strings.ToLower(StorageClassMinSizeName):
			// Convert to integer
			iv, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			minSize = iv
		case strings.ToLower(StorageClassStepSizeName):
			// Convert to integer
			iv, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			stepSize = iv
		case strings.ToLower(StorageClassFsTypeName):
			if len(v) != 0 && !IsValidFileSystemType(v) {
				return nil, fmt.Errorf("unsupported filesystem type %s", v)
			}
			fsType = v
		case strings.ToLower(StorageClassReplicaName):
			iv, err := strconv.Atoi(v)
			if err != nil {
				return nil, err
			}
			replica = iv
		case strings.ToLower(StorageClassTagsName):
			if len(v) > 0 {
				tags = strings.Split(strings.ReplaceAll(v, " ", ""), ",")
			}
		}
	}

	if volType == -1 {
		return NewDefaultQingStorageClassFromType(DefaultVolumeType), nil
	} else {
		t := VolumeType(volType)
		if !t.IsValid() {
			return nil, fmt.Errorf("unsupported volume type %d", volType)
		}
		sc := NewDefaultQingStorageClassFromType(t)
		// For backward compatiblility, ignore error
		sc.setTypeSize(maxSize, minSize, stepSize)
		sc.setFsType(fsType)
		sc.setReplica(replica)
		sc.setTags(tags)
		return sc, nil
	}
}

func (sc QingStorageClass) GetDiskType() VolumeType {
	return sc.diskType
}

func (sc QingStorageClass) GetMinSizeByte() int64 {
	return int64(sc.minSize) * common.Gib
}

func (sc QingStorageClass) GetMaxSizeByte() int64 {
	return int64(sc.maxSize) * common.Gib
}
func (sc QingStorageClass) GetStepSizeByte() int64 {
	return int64(sc.stepSize) * common.Gib
}

func (sc QingStorageClass) GetFsType() string {
	return sc.fsType
}

func (sc QingStorageClass) GetReplica() int {
	return sc.replica
}

func (sc QingStorageClass) GetTags() []string {
	return sc.tags
}

func (sc *QingStorageClass) setFsType(fs string) error {
	if !IsValidFileSystemType(fs) {
		return fmt.Errorf("unsupported filesystem type %s", fs)
	}
	sc.fsType = fs
	return nil
}

func (sc *QingStorageClass) setReplica(repl int) error {
	if !IsValidReplica(repl) {
		return fmt.Errorf("unsupported replica %d", repl)
	}
	sc.replica = repl
	return nil
}

func (sc *QingStorageClass) setTypeSize(maxSize, minSize, stepSize int) error {
	if maxSize < 0 || minSize <= 0 || stepSize < 0 {
		return nil
	}
	// Ensure volume minSize less than volume maxSize
	if sc.maxSize < sc.minSize {
		return fmt.Errorf("max size must greater than or equal to min size")
	}
	sc.maxSize, sc.minSize, sc.stepSize = maxSize, minSize, stepSize
	return nil
}

func (sc *QingStorageClass) setTags(tagsStr []string) {
	sc.tags = tagsStr
}

// FormatVolumeSize transfer to proper volume size
func (sc QingStorageClass) FormatVolumeSizeByte(sizeByte int64) int64 {
	if sizeByte <= sc.GetMinSizeByte() {
		return sc.GetMinSizeByte()
	} else if sizeByte > sc.GetMaxSizeByte() {
		return sc.GetMaxSizeByte()
	}
	if sizeByte%sc.GetStepSizeByte() != 0 {
		sizeByte = (sizeByte/sc.GetStepSizeByte() + 1) * sc.GetStepSizeByte()
	}
	if sizeByte > sc.GetMaxSizeByte() {
		return sc.GetMaxSizeByte()
	}
	return sizeByte
}

// Required Volume Size
func (sc QingStorageClass) GetRequiredVolumeSizeByte(capRange *csi.CapacityRange) (int64, error) {
	if capRange == nil {
		return int64(sc.minSize) * common.Gib, nil
	}
	res := int64(0)
	if capRange.GetRequiredBytes() > 0 {
		res = capRange.GetRequiredBytes()
	}
	res = sc.FormatVolumeSizeByte(res)
	if capRange.GetLimitBytes() > 0 && res > capRange.GetLimitBytes() {
		return -1, fmt.Errorf("volume required bytes %d greater than limit bytes %d", res, capRange.GetLimitBytes())
	}
	return res, nil
}
