package block

import (
	"fmt"
	"strconv"
)

type qingStorageClass struct {
	VolumeType    int `json:"type"`
	VolumeMaxSize int `json:"maxSize"`
	VolumeMinSize int `json:"minSize"`
}

func NewDefaultQingStorageClass() *qingStorageClass {
	return &qingStorageClass{
		VolumeType:    0,
		VolumeMaxSize: 500,
		VolumeMinSize: 10,
	}
}

func NewQingStorageClassFromMap(opt map[string]string) (*qingStorageClass, error) {
	sc := NewDefaultQingStorageClass()
	// volume type
	if sVolType, ok := opt["type"]; ok {
		if iVolType, err := strconv.Atoi(sVolType); err != nil {
			return nil, err
		} else {
			sc.VolumeType = iVolType
		}
	}

	// Get volume maxsize +optional
	if sMaxSize, ok := opt["maxSize"]; ok {
		if iMaxSize, err := strconv.Atoi(sMaxSize); err != nil {
			return nil, err
		} else {
			if iMaxSize < 0 {
				return nil, fmt.Errorf("MaxSize must not less than zero")
			}
			sc.VolumeMaxSize = iMaxSize
		}
	}

	// Get volume minsize +optional
	if sMinSize, ok := opt["minSize"]; ok {
		if iMinSize, err := strconv.Atoi(sMinSize); err != nil {
			return nil, err
		} else {
			if iMinSize < 0 {
				return nil, fmt.Errorf("MinSize must not less than zero")
			}
			sc.VolumeMinSize = iMinSize
		}
	}

	// Ensure volume minSize less than volume maxSize
	if sc.VolumeMaxSize < sc.VolumeMinSize {
		return nil, fmt.Errorf("Volume maxSize must greater than or equal to volume minSize")
	}
	return sc, nil
}

func (sc qingStorageClass) FormatVolumeSize(size int) int {
	if size <= sc.VolumeMinSize {
		return sc.VolumeMinSize
	} else if size >= sc.VolumeMaxSize {
		return sc.VolumeMaxSize
	}
	if size%10 != 0 {
		size = (size/10 + 1) * 10
	}
	if size >= sc.VolumeMaxSize {
		return sc.VolumeMaxSize
	}
	return size
}
