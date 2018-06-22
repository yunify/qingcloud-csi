package block

import (
	"github.com/golang/glog"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	Instance_Status_PENDING    string = "pending"
	Instance_Status_RUNNING    string = "running"
	Instance_Status_STOPPED    string = "stopped"
	Instance_Status_SUSPENDED  string = "suspended"
	Instance_Status_TERMINATED string = "terminated"
	Instance_Status_CEASED     string = "ceased"
)

type instanceProvider struct {
	instanceService *qcservice.InstanceService
	jobService      *qcservice.JobService
	storageClass    *qingStorageClass
}

func newInstanceProvider(sc *qingStorageClass) (*instanceProvider, error) {
	// create config
	config := sc.getConfig()
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
	iv := instanceProvider{
		instanceService: is,
		jobService:      js,
		storageClass:    sc,
	}
	glog.Infof("instance provider init finish, zone: %s",
		*iv.instanceService.Properties.Zone)
	return &iv, nil
}

// Find instance by instance ID
// Return: 	nil,	nil: 	not found instance
//			instance, nil: 	found instance
//			nil, 	error:	internal error
func (iv *instanceProvider) findInstance(id string) (instance *qcservice.Instance, err error) {
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
		return nil, status.Errorf(
			codes.Internal, "call DescribeInstances err: instance id %s in %s", id, iv.instanceService.Config.Zone)
	}
	// not found instances
	switch *output.TotalCount {
	case 0:
		return nil, nil
	case 1:
		return output.InstanceSet[0], nil
	default:
		return nil, status.Errorf(
			codes.OutOfRange, "find duplicate instances id %s in %s", id, iv.instanceService.Config.Zone)
	}
}
