package block

import (
	"fmt"
	"github.com/golang/glog"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	BlockVolume_Status_PENDING   string = "pending"
	BlockVolume_Status_AVAILABLE string = "available"
	BlockVolume_Status_INUSE     string = "in-use"
	BlockVolume_Status_SUSPENDED string = "suspended"
	BlockVolume_Status_DELETED   string = "deleted"
	BlockVolume_Status_CEASED    string = "ceased"
)

type volumeManager struct {
	volumeService *qcservice.VolumeService
	jobService    *qcservice.JobService
}

func NewVolumeManagerWithConfig(config *qcconfig.Config) (*volumeManager, error) {
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
	vp := volumeManager{
		volumeService: vs,
		jobService:    js,
	}
	glog.Infof("Finished new volume manager")
	return &vp, nil
}

func NewVolumeManager() (*volumeManager, error) {
	config, err := ReadConfigFromFile(ConfigFilePath)
	if err != nil{
		return nil, err
	}
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
	vp := volumeManager{
		volumeService: vs,
		jobService:    js,
	}
	glog.Infof("Finished new volume manager")
	return &vp, nil
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
		return nil,
			fmt.Errorf("call DescribeVolumes err: volume id %s in %s", id, vm.volumeService.Config.Zone)
	}
	switch *output.TotalCount {
	// Not found volumes
	case 0:
		return nil, nil
	// Found one volume
	case 1:
		return output.VolumeSet[0], nil
	// Found duplicate volumes
	default:
		return nil,
			fmt.Errorf("call DescribeVolumes err: find duplicate volumes, volume id %s in %s", id, vm.volumeService.Config.Zone)
	}
}

// Find volume by volume name
// Return: 	nil, 		nil: 	not found volumes
//			volumes,	nil:	found volume
//			nil,		error:	internal error
func (vm *volumeManager) FindVolumeByName(name string) (volume *qcservice.Volume, err error) {
	// Set input arguements
	input := qcservice.DescribeVolumesInput{}
	input.SearchWord = &name
	// Call DescribeVolumes
	output, err := vm.volumeService.DescribeVolumes(&input)
	// Handle error
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		return nil,
		fmt.Errorf("call DescribeVolumes err: volume name %s in %s", name, vm.volumeService.Config.Zone)
	}
	// Not found volumes
	switch *output.TotalCount {
	case 0:
		return nil, nil
	case 1:
		return output.VolumeSet[0], nil
	default:
		return nil,
		fmt.Errorf("call DescribeVolumes err: find duplicate volumes, volume name %s in %s", name, vm.volumeService.Config.Zone)
	}
}

// create volume
func (vm *volumeManager) CreateVolume(volumeName string, requestSize int, sc qingStorageClass) ( volumeId string, err error) {
	// 0. Set CreateVolume args
	// set input value
	input := &qcservice.CreateVolumesInput{}
	// create volume count
	count := 1
	input.Count = &count
	// volume provisioner size
	size := sc.formatVolumeSize(requestSize)
	input.Size = &size
	// create volume name
	input.VolumeName = &volumeName
	// volume provisioner type
	input.VolumeType = &sc.VolumeType

	// 1. Create volume
	glog.Infof("CreateVolume request size: %d GB, zone: %s, type: %d, count: %d, name: %s",
		*input.Size, *vm.volumeService.Properties.Zone, *input.VolumeType, *input.Count, *input.VolumeName)
	output, err := vm.volumeService.CreateVolumes(input)
	if err != nil {
		return "", err
	}
	// check output
	if *output.RetCode != 0 {
		glog.Warningf("call CreateVolumes return %d, name %s",
			*output.RetCode, volumeName)
	}else{
		volumeId = *output.Volumes[0]
		glog.Info("call CreateVolume name %s id %s succeed", volumeName, volumeId)
	}

	return *output.Volumes[0], nil
}

// delete volume
func (vm *volumeManager) DeleteVolume(id string) error {
	// set input value
	input := &qcservice.DeleteVolumesInput{}
	input.Volumes = append(input.Volumes, &id)
	// delete volume
	glog.Infof("DeleteVolume request id: %s, zone: %s",
		id, *vm.volumeService.Properties.Zone)
	output, err := vm.volumeService.DeleteVolumes(input)
	if err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		glog.Errorf("call DeleteVolumes %s failed, return %d",
			id, *output.RetCode)
	}else{
		glog.Infof("call DeleteVolume %s succeed", id)
	}
	return nil
}

// check volume attaching to instance
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
			codes.NotFound, "volume %s not found in %s", volumeId, zone)
	}

	if volumeItem.Instance != nil && *volumeItem.Instance.InstanceID == instanceId {
		return true, nil
	} else {
		return false, nil
	}
}

// attach volume
func (vm *volumeManager) AttachVolume(volumeId string, instanceId string) error{
	glog.Infof("volumeId: %s, instanceId: %s", volumeId, instanceId)
	zone := *vm.volumeService.Properties.Zone
	// check volume status
	flag, err := vm.IsAttachedToInstance(volumeId, instanceId)
	if err != nil {
		return  err
	}
	if flag {
		vol, err := vm.FindVolume(volumeId)
		if err == nil && vol.Instance != nil {
			glog.Infof("volume %s has been attached to instance %s device %s in zone %s", volumeId, instanceId, *vol.Instance.Device, zone)
			return nil
		}
		return nil
	}
	// set input parameter
	input := &qcservice.AttachVolumesInput{}
	input.Volumes = append(input.Volumes, &volumeId)
	input.Instance = &instanceId
	// attach volume
	glog.Infof("call AttachVolume request volume id: %s, instance id: %s, zone: %s", volumeId, instanceId, zone)
	output, err := vm.volumeService.AttachVolumes(input)
	if err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		return fmt.Errorf("call AttachVolume return %d, volume id %s", *output.RetCode, volumeId)
	}
	return  nil
}

// detach volume
func (vm *volumeManager) DetachVolume(volumeId string, instanceId string) error {
	// check volume status
	if ok, _ := vm.IsAttachedToInstance(volumeId, instanceId); ok == false {
		return fmt.Errorf("volume %s is not attached to instance %s in zone %s",
				volumeId, instanceId, *vm.volumeService.Properties.Zone)
	}
	// set input parameter
	input := &qcservice.DetachVolumesInput{}
	input.Volumes = append(input.Volumes, &volumeId)
	input.Instance = &instanceId
	// attach volume
	glog.Infof("call DetachVolume request volume id: %s, instance id: %s, zone: %s",
		volumeId, instanceId, *vm.volumeService.Properties.Zone)
	output, err := vm.volumeService.DetachVolumes(input)
	if err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		glog.Errorf("call DetachVolume return %d, volume id %s",
			*output.RetCode, volumeId)
	}
	return nil
}
