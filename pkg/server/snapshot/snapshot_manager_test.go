package snapshot

import (
	"testing"
)

var (
	// Tester should set these variables before executing unit test.
	volumeId1     string = "vol-uwxrtw0d"
	volumeName1   string = "test2"
	snapshotId1   string = "ss-rnbwvjy5"
	snapshotName1 string = "test1"
)

var getsm = func() SnapshotManager {
	// get storage class
	filePath := "/root/.qingcloud/config.yaml"
	sm, err := NewSnapshotManagerFromFile(filePath)
	if err != nil {
		return nil
	}
	return sm
}

func TestSnapshotManager_FindSnapshot(t *testing.T) {
	sm := getsm()
	// testcase
	testcases := []struct {
		name   string
		id     string
		result bool
	}{
		{
			name:   "Available",
			id:     snapshotId1,
			result: true,
		},
		{
			name:   "Not found",
			id:     snapshotId1 + "fake",
			result: false,
		},
		{
			name:   "By name",
			id:     snapshotName1,
			result: false,
		},
	}

	// test findVolume
	for _, v := range testcases {
		snap, err := sm.FindSnapshot(v.id)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		res := snap != nil
		if res != v.result {
			t.Errorf("name: %s, expect %t, actually %t", v.name, v.result, res)
		}
	}
}

func TestSnapshotManager_FindSnapshotByName(t *testing.T) {
	sm := getsm()
	// testcase
	testcases := []struct {
		name     string
		snapshot string
		result   bool
	}{
		{
			name:     "Available",
			snapshot: snapshotName1,
			result:   true,
		},
		{
			name:     "Ceased",
			snapshot: "sanity",
			result:   false,
		},
		{
			name:     "Volume id",
			snapshot: snapshotId1,
			result:   false,
		},
		{
			name:     "Substring",
			snapshot: string((snapshotName1)[:2]),
			result:   false,
		},
		{
			name:     "Null string",
			snapshot: "",
			result:   false,
		},
	}

	// test findVolume
	for _, v := range testcases {
		snap, err := sm.FindSnapshotByName(v.snapshot)
		if err != nil {
			t.Error("find volume error: ", err.Error())
		}
		res := snap != nil
		if res != v.result {
			t.Errorf("name %s, expect %t, actually %t", v.name, v.result, res)
		}
	}
}

func TestSnapshotManager_CreateSnapshot(t *testing.T) {
	sm := getsm()

	testcases := []struct {
		name        string
		snapName    string
		sourceVolId string
		result      bool
		snapId      string
	}{
		{
			name:        "create snapshot name unittest-1",
			snapName:    "unittest-1",
			sourceVolId: volumeId1,
			result:      true,
			snapId:      "",
		},
		{
			name:        "create snapshot name unittest-1 repeatedly",
			snapName:    "unittest-1",
			sourceVolId: volumeId1,
			result:      true,
			snapId:      "",
		},
		{
			name:        "create volume name unittest-2",
			snapName:    "unittest-2",
			sourceVolId: volumeId1,
			result:      true,
			snapId:      "",
		},
	}
	for i, v := range testcases {
		snapId, err := sm.CreateSnapshot(v.snapName, v.sourceVolId)
		if err != nil {
			t.Errorf("test %s: %s", v.name, err.Error())
		} else {
			snap, _ := sm.FindSnapshot(snapId)
			testcases[i].snapId = *snap.SnapshotID
			if *snap.SnapshotName != v.snapName {
				t.Errorf("test %s: expect %s but actually %s", v.name, v.snapName, *snap.SnapshotName)
			}
		}
	}
}

func TestDeleteVolume(t *testing.T) {
	sm := getsm()
	// testcase
	testcases := []struct {
		name    string
		id      string
		isError bool
	}{
		{
			name:    "delete first volume",
			id:      snapshotId1,
			isError: false,
		},
		{
			name:    "delete first volume repeatedly",
			id:      snapshotId1,
			isError: true,
		},
		{
			name:    "delete not exist volume",
			id:      "ss-1234567",
			isError: true,
		},
	}
	for _, v := range testcases {
		err := sm.DeleteSnapshot(v.id)
		if err != nil && !v.isError {
			t.Errorf("error name %s: %s", v.name, err.Error())
		}
	}
}
