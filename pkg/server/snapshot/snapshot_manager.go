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

package snapshot

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/yunify/qingcloud-csi/pkg/server"
	qcclient "github.com/yunify/qingcloud-sdk-go/client"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
)

// https://github.com/yunify/qingcloud-sdk-go/blob/c8f8d40dd4793219c129b7516d6f8ae130bc83c9/service/types.go#L2763
// Available status values: pending, available, suspended, deleted, ceased
const (
	SnapshotStatusPending   string = "pending"
	SnapshotStatusAvailable string = "available"
	SnapshotStatusSuspended string = "suspended"
	SnapshotStatusDeleted   string = "deleted"
	SnapshotStatusCeased    string = "ceased"
)

// https://github.com/yunify/qingcloud-sdk-go/blob/c8f8d40dd4793219c129b7516d6f8ae130bc83c9/service/types.go#L2770
// Available transition status values: creating, suspending, resuming, deleting, recovering
const (
	SnapshotTransitionStatusCreating   string = "creating"
	SnapshotTransitionStatusSuspending string = "suspending"
	SnapshotTransitionStatusResuming   string = "resuming"
	SnapshotTransitionStatusDeleting   string = "deleting"
	SnapshotTransitionStatusRecovering string = "recovering"
)

const (
	SnapshotFull      int = 1
	SnapshotIncrement int = 0
)

type SnapshotManager interface {
	FindSnapshot(id string) (snapshot *qcservice.Snapshot, err error)
	FindSnapshotByName(name string) (snapshot *qcservice.Snapshot, err error)
	CreateSnapshot(snapshotName string, resourceId string) (snapshotId string, err error)
	DeleteSnapshot(snapshotId string) error
	GetZone() string
	waitJob(jobId string) error
}

type snapshotManager struct {
	snapshotService *qcservice.SnapshotService
	jobService      *qcservice.JobService
}

// NewSnapshotManagerFromConfig
// Create snapshot manager from config
func NewSnapshotManagerFromConfig(config *qcconfig.Config) (SnapshotManager, error) {
	// initial qingcloud iaas service
	qs, err := qcservice.Init(config)
	if err != nil {
		return nil, err
	}
	// create snapshot service
	ss, _ := qs.Snapshot(config.Zone)
	// create job service
	js, _ := qs.Job(config.Zone)
	// initial snapshot manager
	sm := snapshotManager{
		snapshotService: ss,
		jobService:      js,
	}
	glog.Infof("Finished initial snapshot manager")
	return &sm, nil
}

// NewSnapshotManagerFromFile
// Create snapshot manager from file
func NewSnapshotManagerFromFile(filePath string) (SnapshotManager, error) {
	config, err := server.ReadConfigFromFile(filePath)
	if err != nil {
		glog.Errorf("Failed read config file [%s], error: [%s]", filePath, err.Error())
		return nil, err
	}
	glog.Infof("Succeed read config file [%s]", filePath)
	return NewSnapshotManagerFromConfig(config)
}

// Find snapshot by snapshot id
// Return: 	nil,	nil: 	not found snapshot
//			snapshot, nil: 	found snapshot
//			nil, 	error:	internal error
func (sm *snapshotManager) FindSnapshot(id string) (snapshot *qcservice.Snapshot, err error) {
	// Set DescribeSnapshot input
	input := qcservice.DescribeSnapshotsInput{}
	input.Snapshots = append(input.Snapshots, &id)
	// Call describe snapshot
	output, err := sm.snapshotService.DescribeSnapshots(&input)
	// 1. Error is not equal to nil.
	if err != nil {
		return nil, err
	}
	// 2. Return code is not equal to 0.
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf("Call IaaS DescribeSnapshot err: snapshot id %s in %s",
			id, sm.snapshotService.Config.Zone)
	}
	switch *output.TotalCount {
	// Not found snapshot
	case 0:
		return nil, nil
	// Found one snapshot
	case 1:
		if *output.SnapshotSet[0].Status == SnapshotStatusCeased ||
			*output.SnapshotSet[0].Status == SnapshotStatusDeleted {
			return nil, nil
		}
		return output.SnapshotSet[0], nil
	// Found duplicate snapshots
	default:
		return nil,
			fmt.Errorf("Call IaaS DescribeSnapshot err: find duplicate snapshot, snapshot id %s in %s",
				id, sm.snapshotService.Config.Zone)
	}
}

// Find snapshot by snapshot name
// In Qingcloud IaaS platform, it is possible that two snapshots have the same name.
// In Kubernetes, the CO will set a unique PV name.
// CSI driver take the PV name as a snapshot name.
// Return: 	nil, 		nil: 	not found snapshots
//			snapshots,	nil:	found snapshot
//			nil,		error:	internal error
func (sm *snapshotManager) FindSnapshotByName(name string) (snapshot *qcservice.Snapshot, err error) {
	if len(name) == 0 {
		return nil, nil
	}
	// Set input arguments
	input := qcservice.DescribeSnapshotsInput{}
	input.SearchWord = &name
	// Call DescribeSnapshot
	output, err := sm.snapshotService.DescribeSnapshots(&input)
	// Handle error
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf("Call IaaS DescribeSnapshots err: snapshot name %s in %s",
			name, sm.snapshotService.Config.Zone)
	}
	// Not found snapshots
	for _, v := range output.SnapshotSet {
		if *v.SnapshotName != name {
			continue
		}
		if *v.Status == SnapshotStatusCeased || *v.Status == SnapshotStatusDeleted {
			continue
		}
		return v, nil
	}
	return nil, nil
}

// CreateSnapshot
// 1. format snapshot size
// 2. create snapshot
// 3. wait job
func (sm *snapshotManager) CreateSnapshot(snapshotName string, resourceId string) (snapshotId string, err error) {
	// 0. Set CreateSnapshot args
	// set input value
	input := &qcservice.CreateSnapshotsInput{}
	// snapshot name
	input.SnapshotName = &snapshotName
	// full snapshot
	snapshotType := int(SnapshotFull)
	input.IsFull = &snapshotType
	// resource volume id
	input.Resources = []*string{&resourceId}

	// 1. Create snapshot
	glog.Infof("Call IaaS CreateSnapshot request snapshot name: %s, zone: %s, resource id %s, is full snapshot %T",
		*input.SnapshotName, sm.GetZone(), *input.Resources[0], *input.IsFull == SnapshotFull)
	output, err := sm.snapshotService.CreateSnapshots(input)
	if err != nil {
		return "", err
	}
	// check output
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return "", fmt.Errorf(*output.Message)
	}
	snapshotId = *output.Snapshots[0]
	glog.Infof("Call IaaS CreateSnapshots snapshot name %s snapshot id %s succeed", snapshotName, snapshotId)
	return snapshotId, nil
}

// DeleteSnapshot
// 1. delete snapshot by id
// 2. wait job
func (sm *snapshotManager) DeleteSnapshot(snapshotId string) error {
	// set input value
	input := &qcservice.DeleteSnapshotsInput{}
	input.Snapshots = append(input.Snapshots, &snapshotId)
	// delete snapshot
	glog.Infof("Call IaaS DeleteSnapshot request id: %s, zone: %s",
		snapshotId, *sm.snapshotService.Properties.Zone)
	output, err := sm.snapshotService.DeleteSnapshots(input)
	if err != nil {
		return err
	}
	// wait job
	glog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := sm.waitJob(*output.JobID); err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		glog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	glog.Infof("Call IaaS DeleteSnapshot %s succeed", snapshotId)
	return nil
}

// GetZone
// Get current zone in Qingcloud IaaS
func (sm *snapshotManager) GetZone() string {
	if sm == nil || sm.snapshotService == nil || sm.snapshotService.Properties == nil || sm.snapshotService.Properties.
		Zone == nil {
		return ""
	}
	return *sm.snapshotService.Properties.Zone
}

func (vm *snapshotManager) waitJob(jobId string) error {
	err := qcclient.WaitJob(vm.jobService, jobId, server.OperationWaitTimeout, server.WaitInterval)
	if err != nil {
		glog.Error("Call Iaas WaitJob: ", jobId)
		return err
	}
	return nil
}
