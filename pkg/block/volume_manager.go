package block

import (
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"fmt"
)

const (
	VolumeTypePerformance = 0
	VolumeTypeHighPerformance = 3 // Only support BJ2
	VolumeTypeCapacity = 1 // BJ1, AS1: 1; BJ2, GD1: 2
)

type volumeClaim struct{
	VolName string
	VolID string
	VolType int
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

// check existence volume by volume ID
// Return: true: volume exist, false: volume not exist
func (vm *volumeProvisioner)IsVolumeIdExist(id string)(bool,error){
	input := qcservice.DescribeVolumesInput{}
	input.Volumes = append(input.Volumes, &id)
	// consult volume info
	output, err := vm.volumeService.DescribeVolumes(&input)
	if err != nil{
		return false, err
	}
	// IaaS response return code not equal to 0
	if *output.RetCode != 0{
		glog.Errorf("call DescribeVolumes() return %d, volume ID: %s, volume.zone: %s",
			*output.RetCode, id, vm.volumeService.Properties.Zone)
	}

	if *output.TotalCount == 1 &&
		*output.VolumeSet[0].VolumeID == id{
		return true, nil
	}else{
		return false, nil
	}
}

// get volume info
func (vm *volumeProvisioner)getVolumeInfoById(id string)(*qcservice.Volume,error){
	input := qcservice.DescribeVolumesInput{}
	input.Volumes = append(input.Volumes, &id)
	// consult volume info
	output, err := vm.volumeService.DescribeVolumes(&input)
	if err != nil{
		return nil, err
	}
	// IaaS response return code not equal to 0
	if *output.RetCode != 0{
		glog.Errorf("call DescribeVolumes() return %d, volume ID: %s, volume.zone: %s",
			*output.RetCode, id, vm.volumeService.Properties.Zone)
	}

	if *output.TotalCount == 1 &&
		*output.VolumeSet[0].VolumeID == id{
		return output.VolumeSet[0], nil
	}else{
		return nil, errors.New(fmt.Sprintf("total count %d when id = %s in zone %s",
			*output.TotalCount, id, vm.volumeService.Properties.Zone))
	}
}

// create volume
func (vm *volumeProvisioner)CreateVolume(opt *volumeClaim)error{
	// set input value
	input := &qcservice.CreateVolumesInput{}
	size := FormatVolumeSize(opt.VolSizeRequest)
	input.Size = &size
	count := 1
	input.Count = &count
	input.VolumeName = &opt.VolName
	input.VolumeType = &opt.VolType
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
	if len(output.Volumes) == 1{
		opt.VolID = *output.Volumes[0]
		volumeInfo, err := vm.getVolumeInfoById(opt.VolID)
		if err != nil{
			return err
		}
		opt.VolSizeCapacity = *volumeInfo.Size
		return nil
	}else{
		return errors.New(fmt.Sprintf("call CreateVolumes output %d, name %s",
			opt.VolID, opt.VolName))
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
		glog.Warningf("call DeleteVolumes return %d, id %s",
			*output.RetCode, id)
	}
	return nil
}
