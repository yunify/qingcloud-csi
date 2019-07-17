package cloudprovider

import "time"

const (
	// In Qingcloud bare host, the path of the file containing instance id.
	RetryString     = "please try later"
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

const (
	DiskSingleReplicaType  int = 1
	DiskMultiReplicaType   int = 2
	DefaultDiskReplicaType int = DiskMultiReplicaType
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
