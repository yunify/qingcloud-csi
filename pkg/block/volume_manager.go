package block

import (
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"github.com/golang/glog"
	"fmt"
)

type volumeClaim struct{
	VolName string
	VolID string
	// high-performance: hp
	// high-performance-plus: hpp
	// high-capacity: hc
	VolType string
	// VolSizeRequest: unit GB
	VolSizeRequest int
	// VolSizeCapacity: unit GB
	VolSizeCapacity int
}

type volumeProvisioner struct {
	volumeService     *qcservice.VolumeService
	jobService 		  *qcservice.JobService
	volumeType string
}

func newVolumeProvisioner(sc *qingStorageClass)(*volumeProvisioner, error){
	// create config
	config := getConfigFromStorageClass(sc)
	// initial qingcloud iaas service
	qs, err := qcservice.Init(config)
	if err != nil{
		return nil,err
	}
	// create volume service
	vs, _ := qs.Volume(config.Zone)
	// create job service
	js, _ := qs.Job(config.Zone)
	// initial volume provisioner
	vp := volumeProvisioner{
		volumeService: vs,
		jobService: js,
		volumeType: sc.Type,
	}
	glog.Infof("volume provisioner init finish, zone: %s, type: %d", vp.volumeService.Properties.Zone, vp.volumeType)
	return &vp, nil
}

// find volume by volume ID
// return: 	nil,	nil: 	not found volumes
//			volume, nil: 	found volume
//			nil, 	error:	error
func (vm *volumeProvisioner)findVolume(id string)(volume *qcservice.Volume, err error){
	// set describe volume input
	input := qcservice.DescribeVolumesInput{}
	input.Volumes = append(input.Volumes, &id)
	// call describe volume
	output, err := vm.volumeService.DescribeVolumes(&input)
	// error
	if err != nil{
		return nil, err
	}
	if *output.RetCode != 0 {
		return nil, fmt.Errorf("call DescribeVolumes return: %d", output.RetCode)
	}
	// not found volumes
	switch *output.TotalCount {
	case 0:
		return nil,nil
	case 1:
		return output.VolumeSet[0],nil
	default:
		return nil, fmt.Errorf("call DescribeVolumes return %d volumesets", output.TotalCount)
	}
}

// create volume
func (vm *volumeProvisioner)CreateVolume(opt *volumeClaim)error{
	// set input value
	input := &qcservice.CreateVolumesInput{}
	// volume provisioner size
	size := FormatVolumeSize(opt.VolSizeRequest)
	input.Size = &size
	// volume provisioner count
	count := 1
	input.Count = &count
	// volume provisioner name
	input.VolumeName = &opt.VolName
	// volume provisioner type
	volumeType := VOLUME_TYPE_MAP[*vm.volumeService.Properties.Zone][vm.volumeType]
	input.VolumeType = &volumeType
	// create volume
	glog.Infof("call CreateVolume request size: %d GB, zone: %s, type: %d, count: %d, name: %s",
		input.Size, vm.volumeService.Properties.Zone, input.VolumeType, input.Count, input.VolumeName)
	output, err := vm.volumeService.CreateVolumes(input)
	if err != nil{
		return err
	}
	// check output
	if *output.RetCode != 0{
		glog.Warningf("call CreateVolumes return %d, name %s",
			*output.RetCode, opt.VolName)
	}
	// check volume exist
	opt.VolID = *output.Volumes[0]
	volumeInfo, err := vm.findVolume(opt.VolID)
	if err != nil{
		return err
	}else{
		opt.VolSizeCapacity = *volumeInfo.Size
		return nil
	}
}

// delete volume
func (vm *volumeProvisioner)DeleteVolume(id string)error{
	// set input value
	input := &qcservice.DeleteVolumesInput{}
	input.Volumes = append(input.Volumes, &id)
	// delete volume
	glog.Infof("call DeleteVolume request id: %s, zone: %s",
		id, vm.volumeService.Properties.Zone)
	output, err := vm.volumeService.DeleteVolumes(input)
	if err != nil{
		return err
	}

	// check output
	if *output.RetCode != 0{
		glog.Errorf("call DeleteVolumes return %d, id %s",
			*output.RetCode, id)
	}
	return nil
}

// check volume attaching to instance
func (vm *volumeProvisioner)isAttachedToInstance(volumeId *string, instanceId *string) bool{
	// get volume item
	volumeItem, err := vm.findVolume(*volumeId)
	if err != nil{
		glog.Errorf("find volume error: %s", err.Error())
	}
	if volumeItem == nil{
		return false
	}
	if *volumeItem.Instance.InstanceID == *instanceId{
		return true
	}else{
		return false
	}
}

// attach volume
func (vm *volumeProvisioner)AttachVolume(volumeId *string, instanceId *string) error{
	// check volume status
	if vm.isAttachedToInstance(volumeId, instanceId){
		glog.Infof("volume %s has been attached to instance %s in zone %s",
			*volumeId, *instanceId, *vm.volumeService.Properties.Zone)
		return nil
	}
	// set input parameter
	input:=&qcservice.AttachVolumesInput{}
	input.Volumes = append(input.Volumes, volumeId)
	input.Instance = instanceId
	// attach volume
	glog.Infof("call AttachVolume request volume id: %s, instance id: %s, zone: %s",
		*volumeId, *instanceId, *vm.volumeService.Properties.Zone)
	output, err := vm.volumeService.AttachVolumes(input)
	if err != nil{
		return err
	}
	// check output
	if *output.RetCode != 0{
		glog.Errorf("call AttachVolume return %d, volume id %s",
			*output.RetCode, *volumeId)
	}
	return nil
}
