package block

import (
	"fmt"
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	Instance_Status_PENDING    string = "pending"
	Instance_Status_RUNNING    string = "running"
	Instance_Status_STOPPED    string = "stopped"
	Instance_Status_SUSPENDED  string = "suspended"
	Instance_Status_TERMINATED string = "terminated"
	Instance_Status_CEASED     string = "ceased"
)

type instanceManager struct {
	instanceService *qcservice.InstanceService
	jobService      *qcservice.JobService
}

func NewInstanceManagerWithConfig(config *qcconfig.Config) (*instanceManager, error) {
	// initial qingcloud iaas service
	qs, err := qcservice.Init(config)
	if err != nil {
		return nil, err
	}
	// create volume service
	is, _ := qs.Instance(config.Zone)
	// create job service
	js, _ := qs.Job(config.Zone)
	// initial volume provisioner
	im := instanceManager{
		instanceService: is,
		jobService:      js,
	}
	glog.Infof("Finish initial volume manager")
	return &im, nil
}

func NewInstanceManager() (*instanceManager, error) {
	// create config
	config, err := ReadConfigFromFile(ConfigFilePath)
	if err != nil {
		return nil, err
	}
	// initial Qingcloud iaas service
	qs, err := qcservice.Init(config)
	if err != nil {
		return nil, err
	}
	// create volume service
	is, _ := qs.Instance(config.Zone)
	// create job service
	js, _ := qs.Job(config.Zone)
	// initial volume provider
	im := instanceManager{
		instanceService: is,
		jobService:      js,
	}
	glog.Infof("instance provider init finish, zone: %s",
		*im.instanceService.Properties.Zone)
	return &im, nil
}

// Find instance by instance ID
// Return: 	nil,	nil: 	not found instance
//			instance, nil: 	found instance
//			nil, 	error:	internal error
func (iv *instanceManager) FindInstance(id string) (instance *qcservice.Instance, err error) {
	// set describe instance input
	input := qcservice.DescribeInstancesInput{}
	input.Instances = append(input.Instances, &id)
	// call describe instance
	output, err := iv.instanceService.DescribeInstances(&input)
	// error
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf(*output.Message)
	}
	// not found instances
	switch *output.TotalCount {
	case 0:
		return nil, nil
	case 1:
		if *output.InstanceSet[0].Status == Instance_Status_CEASED || *output.InstanceSet[0].Status == Instance_Status_TERMINATED {
			return nil, nil
		} else {
			return output.InstanceSet[0], nil
		}
		return output.InstanceSet[0], nil
	default:
		return nil, fmt.Errorf("Find duplicate instances id %s in %s", id, iv.instanceService.Config.Zone)
	}
}
