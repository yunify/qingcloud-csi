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
	"io/ioutil"
	"k8s.io/klog"
	"strings"
)

// Check replica
// Support: 2 MultiReplicas, 1 SingleReplica
func IsValidReplica(replica int) bool {
	switch replica {
	case DiskMultiReplicaType:
		return true
	case DiskSingleReplicaType:
		return true
	default:
		return false
	}
}

// Check file system type
// Support: ext3, ext4 and xfs
func IsValidFileSystemType(fs string) bool {
	switch fs {
	case common.FileSystemExt3:
		return true
	case common.FileSystemExt4:
		return true
	case common.FileSystemXfs:
		return true
	default:
		return false
	}
}

// FormatVolumeSize transfer to proper volume size
func FormatVolumeSize(volType VolumeType, volSize int) int {
	_, ok := VolumeTypeName[volType]
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

func GetInstanceIdFromFile(filepath string) (instanceId string, err error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	instanceId = string(bytes[:])
	instanceId = strings.Replace(instanceId, "\n", "", -1)
	klog.Infof("Getting instance-id: \"%s\"", instanceId)
	return instanceId, nil
}

func NewVolumeCapabilityAccessMode(mode csi.VolumeCapability_AccessMode_Mode) *csi.VolumeCapability_AccessMode {
	return &csi.VolumeCapability_AccessMode{Mode: mode}
}

func NewControllerServiceCapability(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
	return &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func NewNodeServiceCapability(cap csi.NodeServiceCapability_RPC_Type) *csi.NodeServiceCapability {
	return &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

func validateVolumeCapabilities(vcs []*csi.VolumeCapability) error {
	isMnt := false
	isBlk := false

	if vcs == nil {
		return fmt.Errorf("volume capabilities is nil")
	}

	for _, vc := range vcs {
		if err := validateVolumeCapability(vc); err != nil {
			return err
		}
		if blk := vc.GetBlock(); blk != nil {
			isBlk = true
		}
		if mnt := vc.GetMount(); mnt != nil {
			isMnt = true
		}
	}

	if isBlk && isMnt {
		return fmt.Errorf("both mount and block volume capabilities specified")
	}

	return nil
}

func validateVolumeCapability(vc *csi.VolumeCapability) error {
	if err := validateAccessMode(vc.GetAccessMode()); err != nil {
		return err
	}
	blk := vc.GetBlock()
	mnt := vc.GetMount()
	if mnt == nil && blk == nil {
		return fmt.Errorf("must specify an access type")
	}
	if mnt != nil && blk != nil {
		return fmt.Errorf("specified both mount and block access types")
	}
	return nil
}

func validateAccessMode(am *csi.VolumeCapability_AccessMode) error {
	if am == nil {
		return fmt.Errorf("access mode is nil")
	}

	switch am.GetMode() {
	case csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER:
	case csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY:
	case csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY:
	case csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER:
		return fmt.Errorf("MULTI_NODE_MULTI_WRITER access mode is not yet supported for PD")
	default:
		return fmt.Errorf("%v access mode is not supported for for PD", am.GetMode())
	}
	return nil
}
