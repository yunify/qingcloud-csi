package server

import "time"

const (
	// In Qingcloud bare host, the path of the file containing instance id.
	InstanceFilePath     = "/etc/qingcloud/instance-id"
	RetryString          = "please try later"
	Int64Max             = int64(^uint64(0) >> 1)
	WaitInterval         = 10 * time.Second
	OperationWaitTimeout = 180 * time.Second
)

const (
	Kib    int64 = 1024
	Mib    int64 = Kib * 1024
	Gib    int64 = Mib * 1024
	Gib100 int64 = Gib * 100
	Tib    int64 = Gib * 1024
	Tib100 int64 = Tib * 100
)

const (
	FileSystemExt3    string = "ext3"
	FileSystemExt4    string = "ext4"
	FileSystemXfs     string = "xfs"
	FileSystemDefault string = FileSystemExt4
)

const (
	SingleReplica  int = 1
	MultiReplica   int = 2
	DefaultReplica int = MultiReplica
)

const (
	QingCloudSingleReplica string = "rpp-00000001"
	QingCloudMultiReplica  string = "rpp-00000002"
)

var QingCloudReplName = map[int]string{
	1: QingCloudSingleReplica,
	2: QingCloudMultiReplica,
}

type ServerConfig struct {
	instanceId       string
	configFilePath   string
	maxVolumePerNode int64
}

const (
	HighPerformanceDiskType      int = 0
	HighCapacityDiskType         int = 2
	SuperHighPerformanceDiskType int = 3
	BasicDiskType                int = 100
	SSDEnterpriseDiskType        int = 200
	NeonSANDiskType              int = 5
)

// convert volume type to string
// https://docs.qingcloud.com/product/api/action/volume/create_volumes.html
var VolumeTypeToString = map[int]string{
	0:   "HighPerformance",
	2:   "HighCapacity",
	3:   "SuperHighPerformance",
	100: "Basic",
	200: "SSDEnterprise",
	5:   "NeonSAN",
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
