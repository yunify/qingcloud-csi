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

package storageclass

import (
	"fmt"
	"github.com/yunify/qingcloud-csi/pkg/server"
	"strconv"
)

type QingStorageClass struct {
	VolumeType     int    `json:"type"`
	VolumeMaxSize  int    `json:"maxSize"`
	VolumeMinSize  int    `json:"minSize"`
	VolumeStepSize int    `json:"stepSize"`
	VolumeFsType   string `json:"fsType"`
	VolumeReplica  int    `json:"replica"`
}

// NewDefaultQingStorageClass create default qingStorageClass object
func NewDefaultQingStorageClass() *QingStorageClass {
	return NewDefaultQingStorageClassFromType(server.HighCapacityDiskType)
}

// NewDefaultQingStorageClassFromType create default qingStorageClass by specified volume type
func NewDefaultQingStorageClassFromType(volumeType int) *QingStorageClass {
	if server.IsValidVolumeType(volumeType) != true {
		return nil
	}
	return &QingStorageClass{
		VolumeType:     volumeType,
		VolumeMaxSize:  server.VolumeTypeToMaxSize[volumeType],
		VolumeMinSize:  server.VolumeTypeToMinSize[volumeType],
		VolumeStepSize: server.VolumeTypeToStepSize[volumeType],
		VolumeFsType:   server.FileSystemDefault,
		VolumeReplica:  server.DefaultReplica,
	}
}

// NewQingStorageClassFromMap create qingStorageClass object from map
func NewQingStorageClassFromMap(opt map[string]string) (*QingStorageClass, error) {
	sVolType, volTypeOk := opt["type"]
	sMaxSize, maxSizeOk := opt["maxSize"]
	sMinSize, minSizeOk := opt["minSize"]
	sStepSize, stepSizeOk := opt["stepSize"]
	sFsType, fsTypeOk := opt["fsType"]
	sReplica, replicaOk := opt["replica"]
	if volTypeOk == false {
		return NewDefaultQingStorageClass(), nil
	}
	// Convert volume type to integer
	iVolType, err := strconv.Atoi(sVolType)
	if err != nil {
		return nil, err
	}
	sc := NewDefaultQingStorageClassFromType(iVolType)
	if maxSizeOk == true && minSizeOk == true && stepSizeOk == true {
		// Get volume maxsize
		iMaxSize, err := strconv.Atoi(sMaxSize)
		if err != nil {
			return nil, err
		}
		if iMaxSize < 0 {
			return nil, fmt.Errorf("MaxSize must not less than zero")
		}
		sc.VolumeMaxSize = iMaxSize
		// Get volume minsize
		iMinSize, err := strconv.Atoi(sMinSize)
		if err != nil {
			return nil, err
		}
		if iMinSize < 0 {
			return nil, fmt.Errorf("MinSize must not less than zero")
		}
		sc.VolumeMinSize = iMinSize
		// Ensure volume minSize less than volume maxSize
		if sc.VolumeMaxSize < sc.VolumeMinSize {
			return nil, fmt.Errorf("volume maxSize must greater than or equal to volume minSize")
		}
		// Get volume step size
		iStepSize, err := strconv.Atoi(sStepSize)
		if err != nil {
			return nil, err
		}
		if iStepSize <= 0 {
			return nil, fmt.Errorf("StepSize must greate than zero")
		}
		sc.VolumeStepSize = iStepSize
	}

	if fsTypeOk == true {
		if !server.IsValidFileSystemType(sFsType) {
			return nil, fmt.Errorf("unsupported fsType %s", sFsType)
		}
		sc.VolumeFsType = sFsType
	}

	// Get volume replicas
	if replicaOk == true {
		iReplica, err := strconv.Atoi(sReplica)
		if err != nil {
			return nil, err
		}
		if !server.IsValidReplica(iReplica) {
			return nil, fmt.Errorf("unsupported replicas \"%s\"", sReplica)
		}
		sc.VolumeReplica = iReplica
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
