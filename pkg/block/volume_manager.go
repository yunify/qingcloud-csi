package block

import (
//	qcclient "github.com/yunify/qingcloud-sdk-go/client"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"github.com/golang/glog"
)

const (
	VolumeTypePerformance = 0
	VolumeTypeHighPerformance = 3 // Only support BJ2
	VolumeTypeCapacity = 1 // BJ1, AS1: 1; BJ2, GD1: 2
)

type persistentVolume struct{
	VolName string `json:"volName"`
	VolID string `json:"volID"`
	VolType string `json:"volType"`
	VolSize int64 `json:"volSize"`
}

type volumeManager struct {
	volumeService     *qcservice.VolumeService
	persistentVolume  *persistentVolume
}

func newVolumeManager(config *qcconfig.Config)(*volumeManager, error){
	qcService, err := qcservice.Init(config)
	if err != nil{
		return nil,err
	}
	// create volume service
	volumeService, err := qcService.Volume(config.Zone)
	if err != nil {
		return nil, err
	}

	qc := volumeManager{
		volumeService: volumeService,
		persistentVolume: &persistentVolume{},
	}
	glog.Infof("newVolumeManager init finish, zone: %v", config.Zone)
	return &qc, nil
}

// check existence volume by volume ID
// Return: true: volume exist, false: volume not exist
func (vm *volumeManager)IsVolumeIdExist()(bool,error){
	input := qcservice.DescribeVolumesInput{}
	id := vm.persistentVolume.VolID
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
	// IaaS response total count more than 1
	if *output.TotalCount == 0 {

	}
	if *output.TotalCount > 1 {
		glog.Errorf("has duplicated volume ID %s in zone %s",
			id, vm.volumeService.Properties.Zone)
	}
	if *output.TotalCount == 1 &&
		*output.VolumeSet[0].VolumeID == id{
		return true, nil
	}else{
		return false, nil
	}
}
