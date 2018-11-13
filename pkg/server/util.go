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

type ServerConfig struct {
	instanceId       string
	configFilePath   string
	maxVolumePerNode int64
}

// NewServerConfig create ServerConfig object to get server config
func NewServerConfig(id string, filePath string, volumeNumber int64) *ServerConfig {
	return &ServerConfig{
		instanceId:       id,
		configFilePath:   filePath,
		maxVolumePerNode: volumeNumber,
	}
}

// GetConfigFilePath get config file path
func (cfg *ServerConfig) GetConfigFilePath() string {
	return cfg.configFilePath
}

// GetMaxVolumePerNode gets maximum number of volumes that controller can publish to the node
func (cfg *ServerConfig) GetMaxVolumePerNode() int64 {
	return cfg.maxVolumePerNode
}

// GetCurrentInstanceId gets instance id
func (cfg *ServerConfig) GetCurrentInstanceId() string {
	if len(cfg.instanceId) == 0 {
		cfg.readCurrentInstanceId()
	}
	return cfg.instanceId
}

func (cfg *ServerConfig) readCurrentInstanceId() {
	bytes, err := ioutil.ReadFile(InstanceFilePath)
	if err != nil {
		glog.Errorf("Getting current instance-id error: %s", err.Error())
		os.Exit(1)
	}
	cfg.instanceId = string(bytes[:])
	cfg.instanceId = strings.Replace(cfg.instanceId, "\n", "", -1)
	glog.Infof("Getting current instance-id: \"%s\"", cfg.instanceId)
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
	return int64(num) * Gib
}

// ByteCeilToGib
// Convert Byte to Gib
func ByteCeilToGib(num int64) int {
	if num <= 0 {
		return 0
	}
	res := num / Gib
	if res*Gib < num {
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
