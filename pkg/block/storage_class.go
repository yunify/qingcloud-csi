package block

import (
	"fmt"
	"strconv"
)

type qingStorageClass struct {
	VolumeType      int    `json:"type"`
	VolumeMaxSize   int    `json:"maxSize"`
	VolumeMinSize   int    `json:"minSize"`
}


func NewDefaultQingStorageClass() *qingStorageClass {
	return &qingStorageClass{
		VolumeType:      0,
		VolumeMaxSize:   500,
		VolumeMinSize:   10,
	}
}

func NewQingStorageClassFromMap(opt map[string]string) (*qingStorageClass, error) {
	var ok bool
	sc := NewDefaultQingStorageClass()
	// volume type
	sVolType, ok := opt["type"]
	if !ok {
		return nil, fmt.Errorf("Missing required parameter type")
	}
	iVolType, err := strconv.Atoi(sVolType)
	if err != nil {
		return nil, err
	} else {
		sc.VolumeType = iVolType
	}
	// Get volume maxsize +optional
	sMaxSize, ok := opt["maxSize"]
	iMaxSize, err := strconv.Atoi(sMaxSize)
	if err != nil {
		return nil, err
	} else {
		sc.VolumeMaxSize = iMaxSize
	}
	// Get volume minsize +optional
	sMinSize, ok := opt["minSize"]
	iMinSize, err := strconv.Atoi(sMinSize)
	if err != nil {
		return nil, err
	} else {
		sc.VolumeMinSize = iMinSize
	}
	// Ensure volume minSize less than volume maxSize
	if sc.VolumeMinSize >= sc.VolumeMaxSize {
		return nil, fmt.Errorf("Volume minSize must less than volume maxSize")
	}
	return sc, nil
}

func (sc qingStorageClass) formatVolumeSize(size int) int {
	if size <= sc.VolumeMinSize {
		return sc.VolumeMinSize
	} else if size >= sc.VolumeMaxSize {
		return sc.VolumeMaxSize
	}
	if size%10 != 0 {
		size = (size/10 + 1) * 10
	}
	return size
}
