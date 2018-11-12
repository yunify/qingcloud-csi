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

package instance

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/yunify/qingcloud-csi/pkg/server"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	InstanceStatusPending    string = "pending"
	InstanceStatusRunning    string = "running"
	InstanceStatusStopped    string = "stopped"
	InstanceStatusSuspended  string = "suspended"
	InstanceStatusTerminated string = "terminated"
	InstanceStatusCreased    string = "ceased"
)

type InstanceManager interface {
	FindInstance(id string) (instance *qcservice.Instance, err error)
}

type instanceManager struct {
	instanceService *qcservice.InstanceService
	jobService      *qcservice.JobService
}

// NewInstanceManagerFromConfig: Create instance manager from config
func NewInstanceManagerFromConfig(config *qcconfig.Config) (InstanceManager, error) {
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
	glog.Infof("Finish initial instance manager")
	return &im, nil
}

// NewInstanceManagerFromFile
// Create instance manager from file
func NewInstanceManagerFromFile(filePath string) (InstanceManager, error) {
	// create config
	config, err := server.ReadConfigFromFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewInstanceManagerFromConfig(config)
}

// Find instance by instance ID
// Return: 	nil,	nil: 	not found instance
//			instance, nil: 	found instance
//			nil, 	error:	internal error
func (iv *instanceManager) FindInstance(id string) (instance *qcservice.Instance, err error) {
	// set describe instance input
	input := qcservice.DescribeInstancesInput{}
	var seeCluster int = 1
	input.IsClusterNode = &seeCluster
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
		if *output.InstanceSet[0].Status == InstanceStatusCreased || *output.InstanceSet[0].Status == InstanceStatusTerminated {
			return nil, nil
		}
		return output.InstanceSet[0], nil
	default:
		return nil, fmt.Errorf("Find duplicate instances id %s in %s", id, iv.instanceService.Config.Zone)
	}
}
