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

package server

import (
	"fmt"
	"strconv"
)

type QingStorageClass struct {
	VolumeType     int    `json:"type"`
	VolumeMaxSize  int    `json:"maxSize"`
	VolumeMinSize  int    `json:"minSize"`
	VolumeStepSize int    `json:"stepSize"`
	VolumeFsType   string `json:"fsType"`
}

// NewDefaultQingStorageClass create default qingStorageClass object
func NewDefaultQingStorageClass() *QingStorageClass {
	return &QingStorageClass{
		VolumeType:     0,
		VolumeMaxSize:  500,
		VolumeMinSize:  10,
		VolumeStepSize: 10,
		VolumeFsType:   FileSystemDefault,
	}
}

// NewQingStorageClassFromMap create qingStorageClass object from map
func NewQingStorageClassFromMap(opt map[string]string) (*QingStorageClass, error) {
	sc := NewDefaultQingStorageClass()
	// volume type
	if sVolType, ok := opt["type"]; ok {
		iVolType, err := strconv.Atoi(sVolType)
		if err != nil {
			return nil, err
		}
		sc.VolumeType = iVolType
	}

	// Get volume FsType
	// Default is ext4
	if sFsType, ok := opt["fsType"]; ok {
		if !IsValidFileSystemType(sFsType) {
			return nil, fmt.Errorf("Does not support fsType \"%s\"", sFsType)
		}
		sc.VolumeFsType = sFsType
	}

	// Get volume maxsize
	if sMaxSize, ok := opt["maxSize"]; ok {
		iMaxSize, err := strconv.Atoi(sMaxSize)
		if err != nil {
			return nil, err
		}
		if iMaxSize < 0 {
			return nil, fmt.Errorf("MaxSize must not less than zero")
		}
		sc.VolumeMaxSize = iMaxSize
	}

	// Get volume minsize
	if sMinSize, ok := opt["minSize"]; ok {
		iMinSize, err := strconv.Atoi(sMinSize)
		if err != nil {
			return nil, err
		}
		if iMinSize < 0 {
			return nil, fmt.Errorf("MinSize must not less than zero")
		}
		sc.VolumeMinSize = iMinSize
	}

	// Get volume step
	if sStepSize, ok := opt["stepSize"]; ok {
		iStepSize, err := strconv.Atoi(sStepSize)
		if err != nil {
			return nil, err
		}
		if iStepSize <= 0 {
			return nil, fmt.Errorf("StepSize must greate than zero")
		}
		sc.VolumeStepSize = iStepSize
	}

	// Ensure volume minSize less than volume maxSize
	if sc.VolumeMaxSize < sc.VolumeMinSize {
		return nil, fmt.Errorf("Volume maxSize must greater than or equal to volume minSize")
	}
	return sc, nil
}

// FormatVolumeSize transfer to proper volume size
func (sc QingStorageClass) FormatVolumeSize(size int, step int) int {
	if size <= sc.VolumeMinSize {
		return sc.VolumeMinSize
	} else if size >= sc.VolumeMaxSize {
		return sc.VolumeMaxSize
	}
	if size%step != 0 {
		size = (size/step + 1) * step
	}
	if size >= sc.VolumeMaxSize {
		return sc.VolumeMaxSize
	}
	return size
}
