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

package cloud

import "time"

const (
	// In Qingcloud bare host, the path of the file containing instance id.

	WaitJobInterval = 10 * time.Second
	WaitJobTimeout  = 180 * time.Second
)

// Instance
const (
	InstanceStatusPending    string = "pending"
	InstanceStatusRunning    string = "running"
	InstanceStatusStopped    string = "stopped"
	InstanceStatusSuspended  string = "suspended"
	InstanceStatusTerminated string = "terminated"
	InstanceStatusCreased    string = "ceased"
)

// Snapshot
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

// Volume
const (
	DiskStatusPending   string = "pending"
	DiskStatusAvailable string = "available"
	DiskStatusInuse     string = "in-use"
	DiskStatusSuspended string = "suspended"
	DiskStatusDeleted   string = "deleted"
	DiskStatusCeased    string = "ceased"
)

var DiskReplicaTypeName = map[int]string{
	1: "rpp-00000001",
	2: "rpp-00000002",
}

// Zone
const (
	ZoneStatusActive  = "active"
	ZoneStatusFaulty  = "faulty"
	ZoneStatusDefunct = "defunct"
)

const (
	EnableDescribeSnapshotVerboseMode  = 1
	DisableDescribeSnapshotVerboseMode = 0
)

const (
	EnableDescribeInstanceAppCluster  = 1
	DisableDescribeInstanceAppCluster = 0
)

const (
	EnableDescribeInstanceVerboseMode  = 1
	DisableDescribeInstanceVerboseMode = 0
)

const (
	ResourceTypeVolume   = "volume"
	ResourceTypeSnapshot = "snapshot"
)
