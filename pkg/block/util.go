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

package block

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	// In Qingcloud bare host, the path of the file containing instance id.
	InstanceFilePath = "/etc/qingcloud/instance-id"

	RetryString          = "please try later"
	Int64Max             = int64(^uint64(0) >> 1)
	WaitInterval         = 10 * time.Second
	OperationWaitTimeout = 180 * time.Second
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100
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

var instanceIdFromFile string
var ConfigFilePath string

// CreatePath
// Create file path if it does not exits.
func CreatePath(persistentStoragePath string) error {
	if _, err := os.Stat(persistentStoragePath); os.IsNotExist(err) {
		if err := os.MkdirAll(persistentStoragePath, os.FileMode(0755)); err != nil {
			return err
		}
	} else {
	}
	return nil
}

func readCurrentInstanceId() {
	bytes, err := ioutil.ReadFile(InstanceFilePath)
	if err != nil {
		glog.Errorf("Getting current instance-id error: %s", err.Error())
		os.Exit(1)
	}
	instanceIdFromFile = string(bytes[:])
	instanceIdFromFile = strings.Replace(instanceIdFromFile, "\n", "", -1)
	glog.Infof("Getting current instance-id: \"%s\"", instanceIdFromFile)
}

// GetCurrentInstanceId
// Get instance id
func GetCurrentInstanceId() string {
	if len(instanceIdFromFile) == 0 {
		readCurrentInstanceId()
	}
	return instanceIdFromFile
}

// ReadConfigFromFile
// Read config file from a path and return config
func ReadConfigFromFile(filePath string) (*qcconfig.Config, error) {
	config, err := qcconfig.NewDefault()
	if err != nil {
		return nil, err
	}
	if err = config.LoadConfigFromFilepath(filePath); err != nil {
		return nil, err
	}
	return config, nil
}

// ContainsVolumeCapability
// Does Array of VolumeCapability_AccessMode contain the volume capability of subCaps
func ContainsVolumeCapability(accessModes []*csi.VolumeCapability_AccessMode, subCaps *csi.VolumeCapability) bool {
	for _, cap := range accessModes {
		if cap.GetMode() == subCaps.GetAccessMode().GetMode() {
			return true
		}
	}
	return false
}

// ContainsVolumeCapabilities
// Does array of VolumeCapability_AccessMode contain volume capabilities of subCaps
func ContainsVolumeCapabilities(accessModes []*csi.VolumeCapability_AccessMode, subCaps []*csi.VolumeCapability) bool {
	for _, v := range subCaps {
		if !ContainsVolumeCapability(accessModes, v) {
			return false
		}
	}
	return true
}

// ContainsNodeServiceCapability
// Does array of NodeServiceCapability contain node service capability of subCap
func ContainsNodeServiceCapability(nodeCaps []*csi.NodeServiceCapability, subCap csi.NodeServiceCapability_RPC_Type) bool {
	for _, v := range nodeCaps {
		if strings.Contains(v.String(), subCap.String()) {
			return true
		}
	}
	return false
}

// GibToByte
// Convert GiB to Byte
func GibToByte(num int) int64 {
	if num < 0 {
		return 0
	}
	return int64(num) * gib
}

// ByteCeilToGib
// Convert Byte to Gib
func ByteCeilToGib(num int64) int {
	if num <= 0 {
		return 0
	}
	res := num / gib
	if res*gib < num {
		res += 1
	}
	return int(res)
}

// Check file system type
// Support: ext3, ext4 and xfs
func IsValidFileSystemType(fs string) bool {
	switch fs {
	case FileSystemExt3:
		return true
	case FileSystemExt4:
		return true
	case FileSystemXfs:
		return true
	default:
		return false
	}
}

// Check replica
// Support: 2 MultiReplicas, 1 SingleReplica
func IsValidReplica(replica int) bool {
	switch replica {
	case MultiReplica:
		return true
	case SingleReplica:
		return true
	default:
		return false
	}
}
