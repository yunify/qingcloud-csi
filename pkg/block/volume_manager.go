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
	"fmt"
	"github.com/golang/glog"
	qcclient "github.com/yunify/qingcloud-sdk-go/client"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	BlockVolumeStatusPending   string = "pending"
	BlockVolumeStatusAvailable string = "available"
	BlockVolumeStatusInuse     string = "in-use"
	BlockVolumeStatusSuspended string = "suspended"
	BlockVolumeStatusDeleted   string = "deleted"
	BlockVolumeStatusCeased    string = "ceased"
)

type VolumeManager interface {
	FindVolume(id string) (volume *qcservice.Volume, err error)
	FindVolumeByName(name string) (volume *qcservice.Volume, err error)
	CreateVolume(volumeName string, requestSize int, sc qingStorageClass) (volumeId string, err error)
	DeleteVolume(id string) error
	IsAttachedToInstance(volumeId string, instanceId string) (flag bool, err error)
	AttachVolume(volumeId string, instanceId string) error
	DetachVolume(volumeId string, instanceId string) error
	GetZone() string
	waitJob(jobId string) error
}

type volumeManager struct {
	volumeService *qcservice.VolumeService
	jobService    *qcservice.JobService
}

// NewVolumeManagerFromConfig
// Create volume manager from config
func NewVolumeManagerFromConfig(config *qcconfig.Config) (VolumeManager, error) {
	// initial qingcloud iaas service
	qs, err := qcservice.Init(config)
	if err != nil {
		return nil, err
	}
	// create volume service
	vs, _ := qs.Volume(config.Zone)
	// create job service
	js, _ := qs.Job(config.Zone)
	// initial volume provisioner
	vm := volumeManager{
		volumeService: vs,
		jobService:    js,
	}
	glog.Infof("Finished initial volume manager")
	return &vm, nil
}

// NewVolumeManagerFromFile
// Create volume manager from file
func NewVolumeManagerFromFile(filePath string) (VolumeManager, error) {
	config, err := ReadConfigFromFile(filePath)
	if err != nil {
		glog.Errorf("Failed read config file [%s], error: [%s]", filePath, err.Error())
		return nil, err
	}
	glog.Infof("Succeed read config file [%s]", filePath)
	return NewVolumeManagerFromConfig(config)
}

// Find volume by volume ID
// Return: 	nil,	nil: 	not found volumes
//			volume, nil: 	found volume
//			nil, 	error:	internal error
func (vm *volumeManager) FindVolume(id string) (volume *qcservice.Volume, err error) {
	// Set DescribeVolumes input
	input := qcservice.DescribeVolumesInput{}
	input.Volumes = append(input.Volumes, &id)
	// Call describe volume
	output, err := vm.volumeService.DescribeVolumes(&input)
	// Error:
	// 1. Error is not equal to nil.
	if err != nil {
		return nil, err
	}
	// 2. Return code is not equal to 0.
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf("Call IaaS DescribeVolumes err: volume id %s in %s", id, vm.volumeService.Config.Zone)
	}
	switch *output.TotalCount {
	// Not found volumes
	case 0:
		return nil, nil
	// Found one volume
	case 1:
		if *output.VolumeSet[0].Status == BlockVolumeStatusCeased || *output.VolumeSet[0].Status == BlockVolumeStatusDeleted {
			return nil, nil
		}
		return output.VolumeSet[0], nil
	// Found duplicate volumes
	default:
		return nil,
			fmt.Errorf("Call IaaS DescribeVolumes err: find duplicate volumes, volume id %s in %s", id, vm.volumeService.Config.Zone)
	}
}

// Find volume by volume name
// In Qingcloud IaaS platform, it is possible that two volumes have the same name.
// In Kubernetes, the CO will set a unique PV name.
// CSI driver take the PV name as a volume name.
// Return: 	nil, 		nil: 	not found volumes
//			volumes,	nil:	found volume
//			nil,		error:	internal error
func (vm *volumeManager) FindVolumeByName(name string) (volume *qcservice.Volume, err error) {
	if len(name) == 0 {
		return nil, nil
	}
	// Set input arguments
	input := qcservice.DescribeVolumesInput{}
	input.SearchWord = &name
	// Call DescribeVolumes
	output, err := vm.volumeService.DescribeVolumes(&input)
	// Handle error
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf("Call IaaS DescribeVolumes err: volume name %s in %s", name, vm.volumeService.Config.Zone)
	}
	// Not found volumes
	for _, v := range output.VolumeSet {
		if *v.VolumeName != name {
			continue
		}
		if *v.Status == BlockVolumeStatusCeased || *v.Status == BlockVolumeStatusDeleted {
			continue
		}
		return v, nil
	}
	return nil, nil
}

// CreateVolume
// 1. format volume size
// 2. create volume
// 3. wait job
func (vm *volumeManager) CreateVolume(volumeName string, requestSize int, sc qingStorageClass) (volumeId string, err error) {
	// 0. Set CreateVolume args
	// set input value
	input := &qcservice.CreateVolumesInput{}
	// create volume count
	count := 1
	input.Count = &count
	// volume provisioner size
	size := sc.FormatVolumeSize(requestSize, sc.VolumeStepSize)
	input.Size = &size
	// create volume name
	input.VolumeName = &volumeName
	// volume provisioner type
	input.VolumeType = &sc.VolumeType
	// volume replicas
	replica := QingCloudReplName[sc.VolumeReplica]
	input.Repl = &replica

	// 1. Create volume
	glog.Infof("Call IaaS CreateVolume request size: %d GB, zone: %s, type: %d, count: %d, replica: %s, name: %s",
		*input.Size, *vm.volumeService.Properties.Zone, *input.VolumeType, *input.Count, *input.Repl, *input.VolumeName)
	output, err := vm.volumeService.CreateVolumes(input)
	if err != nil {
		return "", err
	}
	// wait job
	glog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := vm.waitJob(*output.JobID); err != nil {
		return "", err
	}
	// check output
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return "", fmt.Errorf(*output.Message)
	}
	volumeId = *output.Volumes[0]
	glog.Infof("Call IaaS CreateVolume name %s id %s succeed", volumeName, volumeId)
	return *output.Volumes[0], nil
}

// DeleteVolume
// 1. delete volume by id
// 2. wait job
func (vm *volumeManager) DeleteVolume(id string) error {
	// set input value
	input := &qcservice.DeleteVolumesInput{}
	input.Volumes = append(input.Volumes, &id)
	// delete volume
	glog.Infof("Call IaaS DeleteVolume request id: %s, zone: %s",
		id, *vm.volumeService.Properties.Zone)
	output, err := vm.volumeService.DeleteVolumes(input)
	if err != nil {
		return err
	}
	// wait job
	glog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := vm.waitJob(*output.JobID); err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	glog.Infof("Call IaaS DeleteVolume %s succeed", id)
	return nil
}

// IsAttachedToInstance
// 1. get volume information
// 2. compare input instance id with instance field in volume information
func (vm *volumeManager) IsAttachedToInstance(volumeId string, instanceId string) (flag bool, err error) {
	// zone
	zone := vm.volumeService.Config.Zone

	// get volume item
	volumeItem, err := vm.FindVolume(volumeId)
	if err != nil {
		return false, status.Errorf(codes.Internal, err.Error())
	}
	// check volume exist
	if volumeItem == nil {
		return false, status.Errorf(
			codes.NotFound, "Volume %s not found in %s", volumeId, zone)
	}

	if volumeItem.Instance != nil && *volumeItem.Instance.InstanceID == instanceId {
		return true, nil
	}
	return false, nil
}

// AttachVolume
// 1. get volume information
// 2. attach volume on instance
// 3. wait job
func (vm *volumeManager) AttachVolume(volumeId string, instanceId string) error {
	zone := *vm.volumeService.Properties.Zone
	// check volume status
	vol, err := vm.FindVolume(volumeId)
	if err != nil {
		return err
	}
	if vol == nil {
		return fmt.Errorf("Cannot found volume %s", volumeId)
	}
	if *vol.Instance.InstanceID == "" {
		// set input parameter
		input := &qcservice.AttachVolumesInput{}
		input.Volumes = append(input.Volumes, &volumeId)
		input.Instance = &instanceId
		// attach volume
		glog.Infof("Call IaaS AttachVolume request volume id: %s, instance id: %s, zone: %s", volumeId, instanceId, zone)
		output, err := vm.volumeService.AttachVolumes(input)
		if err != nil {
			return err
		}
		// check output
		if *output.RetCode != 0 {
			glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
			return fmt.Errorf(*output.Message)
		}
		// wait job
		glog.Infof("Call IaaS WaitJob %s", *output.JobID)
		return vm.waitJob(*output.JobID)
	} else {
		if *vol.Instance.InstanceID == instanceId {
			return nil
		}
		return fmt.Errorf("Volume %s has been attached to another instance %s.", volumeId, *vol.Instance.InstanceID)
	}
}

// detach volume
// 1. get volume information
// 2. If volume not attached, return nil.
//   If volume attached, check instance id.
// 3. attach volume
// 4. wait job
func (vm *volumeManager) DetachVolume(volumeId string, instanceId string) error {
	zone := *vm.volumeService.Properties.Zone
	// check volume status
	vol, err := vm.FindVolume(volumeId)
	if err != nil {
		return err
	}
	if vol == nil {
		return fmt.Errorf("Cannot found volume %s", volumeId)
	}
	if *vol.Instance.InstanceID == "" {
		return nil
	} else {
		if *vol.Instance.InstanceID == instanceId || instanceId == "" {
			// set input parameter
			input := &qcservice.DetachVolumesInput{}
			input.Volumes = append(input.Volumes, &volumeId)
			input.Instance = vol.Instance.InstanceID
			// attach volume
			glog.Infof("Call IaaS DetachVolume request volume id: %s, instance id: %s, zone: %s", volumeId, instanceId, zone)
			output, err := vm.volumeService.DetachVolumes(input)
			if err != nil {
				return err
			}
			// check output
			if *output.RetCode != 0 {
				glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
				return fmt.Errorf(*output.Message)
			}
			// wait job
			glog.Infof("Call IaaS WaitJob %s", *output.JobID)
			return vm.waitJob(*output.JobID)
		}
		return fmt.Errorf("Volume %s has been attached to another instance %s", volumeId, *vol.Instance.InstanceID)
	}
}

// GetZone
// Get current zone in Qingcloud IaaS
func (vm *volumeManager) GetZone() string {
	if vm == nil || vm.volumeService == nil || vm.volumeService.Properties == nil || vm.volumeService.Properties.Zone == nil {
		return ""
	}
	return *vm.volumeService.Properties.Zone
}

func (vm *volumeManager) waitJob(jobId string) error {
	err := qcclient.WaitJob(vm.jobService, jobId, OperationWaitTimeout, WaitInterval)
	if err != nil {
		glog.Error("Call Iaas WaitJob: ", jobId)
		return err
	}
	return nil
}
