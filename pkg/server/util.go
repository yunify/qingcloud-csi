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
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// NewServerConfig create ServerConfig object to get server config
func NewServerConfig(id string, filePath string, volumeNumber int64) *ServerConfig {
	sc := &ServerConfig{
		instanceId:       id,
		configFilePath:   filePath,
		maxVolumePerNode: volumeNumber,
	}
	// If instance file existed, plugin SHOULD get instance id
	// from instance file (/etc/qingcloud/instance-id).
	if _, err := os.Stat(InstanceFilePath); !os.IsNotExist(err) {
		sc.readInstanceId()
	}
	return sc
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
func (cfg *ServerConfig) GetInstanceId() string {
	return cfg.instanceId
}

func (cfg *ServerConfig) readInstanceId() {
	bytes, err := ioutil.ReadFile(InstanceFilePath)
	if err != nil {
		glog.Errorf("Getting instance-id error: %s", err.Error())
		os.Exit(1)
	}
	cfg.instanceId = string(bytes[:])
	cfg.instanceId = strings.Replace(cfg.instanceId, "\n", "", -1)
	glog.Infof("Getting instance-id: \"%s\"", cfg.instanceId)
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

// Check volume type
func IsValidVolumeType(volumeType int) bool {
	if _, ok := VolumeTypeToString[volumeType]; ok {
		return true
	}
	return false
}

// EntryFunction print timestamps
func EntryFunction(functionName string) func() {
	start := time.Now()
	glog.Infof("*************** enter %s at %s ***************", functionName, start.String())
	return func() {
		glog.Infof("=============== exit %s (%s since %s) ===============", functionName, time.Since(start),
			start.String())
	}
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

// Get minimal required bytes in capacity range
// Return Values:
//  -1 represent cannot get min required bytes
func GetMinRequiredBytes(requiredBytes, limitBytes []int64) int64 {
	res := int64(0)
	for _, v := range requiredBytes {
		if res < v {
			res = v
		}
	}
	for _, v := range limitBytes {
		if res > v {
			return -1
		}
	}
	return res
}

// Valid capacity bytes in capacity range
func IsValidCapacityBytes(cur int64, requiredBytes, limitBytes []int64) bool {
	res := cur
	for _, v := range requiredBytes {
		if res < v {
			return false
		}
	}
	for _, v := range limitBytes {
		if res > v {
			return false
		}
	}
	return true
}
