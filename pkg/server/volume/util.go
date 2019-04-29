package volume

const (
	DiskVolumeStatusPending   string = "pending"
	DiskVolumeStatusAvailable string = "available"
	DiskVolumeStatusInuse     string = "in-use"
	DiskVolumeStatusSuspended string = "suspended"
	DiskVolumeStatusDeleted   string = "deleted"
	DiskVolumeStatusCeased    string = "ceased"
)

const (
	HighPerformance      string = "HighPerformance"
	HighCapacity         string = "HighCapacity"
	SuperHighPerformance string = "SuperHighPerformance"
	Basic                string = "Basic"
	SSDEnterprise        string = "SSDEnterprise"
	NeonSAN              string = "NeonSAN"
)

// convert volume type to string
// https://docs.qingcloud.com/product/api/action/volume/create_volumes.html
var VolumeTypeToString = map[int]string{
	0:   HighPerformance,
	2:   HighCapacity,
	3:   SuperHighPerformance,
	100: Basic,
	200: SSDEnterprise,
	5:   NeonSAN,
}

var VolumeTypeToStepSize = map[int]int{
	0:   10,
	2:   50,
	3:   10,
	100: 10,
	200: 10,
	5:   100,
}

var VolumeTypeToMinSize = map[int]int{
	0:   10,
	2:   100,
	3:   10,
	100: 10,
	200: 10,
	5:   100,
}

var VolumeTypeToMaxSize = map[int]int{
	0:   2000,
	2:   5000,
	3:   2000,
	100: 2000,
	200: 2000,
	5:   50000,
}

// FormatVolumeSize transfer to proper volume size
func FormatVolumeSize(volType int, volSize int) int {
	_, ok := VolumeTypeToString[volType]
	if ok == false {
		return -1
	}
	volTypeMinSize := VolumeTypeToMinSize[volType]
	volTypeMaxSize := VolumeTypeToMaxSize[volType]
	volTypeStepSize := VolumeTypeToStepSize[volType]
	if volSize <= volTypeMinSize {
		return volTypeMinSize
	} else if volSize >= volTypeMaxSize {
		return volTypeMaxSize
	}
	if volSize%volTypeStepSize != 0 {
		volSize = (volSize/volTypeStepSize + 1) * volTypeStepSize
	}
	if volSize >= volTypeMaxSize {
		return volTypeMaxSize
	}
	return volSize
}
