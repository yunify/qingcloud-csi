package mock

import (
	"github.com/yunify/qingcloud-csi/pkg/cloud"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"reflect"
	"testing"
)

// 1. create volume
// 2. search volume id and name
// 3. attach volume, search volume
// 4. try to delete volume
// 5. detach volume
// 6. delete volume
func TestMockVolumeDB_Create(t *testing.T) {
	var db = NewMockVolumeDB()
	test := struct {
		name    string
		volName string
		volSize int
		volType driver.VolumeType
		volRepl string
		errStr  string
		volId   string
	}{
		name:    "Normal",
		volName: "normal-vol",
		volSize: 10,
		volType: driver.HighPerformanceDiskType,
		volRepl: cloud.DiskReplicaTypeName[cloud.DiskMultiReplicaType],
		errStr:  "",
	}

	// 1 CreateVolume
	volType := test.volType.Int()
	volId, err := db.Create(&qcservice.Volume{
		VolumeName: &test.volName,
		Size:       &test.volSize,
		VolumeType: &volType,
		Repl:       &test.volRepl,
	})
	if err != nil {
		t.Logf("failed to create volume")
	}
	test.volId = volId
	// 2 SearchVolumeId
	volInfoCreated := db.Search(test.volId)
	if *volInfoCreated.VolumeName != test.volName {
		t.Logf("name %s: expect %s but actually %s", test.name, test.volName, *volInfoCreated.VolumeName)
	}
	// SearchVolumeName
	vols := db.SearchName(test.volName)
	if !reflect.DeepEqual(vols[0], volInfoCreated) {
		t.Logf("search volume not equal")
	}
	// 3 AttachVolume
	err = db.Attach(test.volId, "i-12345667")
	if err != nil {
		t.Errorf(err.Error())
	}
	// SearchVolume
	volInfoAttached := db.Search(test.volId)
	if volInfoAttached != nil && *volInfoAttached.Status != cloud.DiskStatusInuse {
		t.Errorf("name %s: expect %s but actually %s", test.name, cloud.DiskStatusInuse,
			*volInfoAttached.Status)
	}
	// 4 Try to DeleteVolume
	err = db.Delete(test.volId)
	if err == nil {
		t.Error("volume should detached before delete")
	}
	// 5 DetachVolume
	err = db.Detach(test.volId, "i-12345667")
	if err != nil {
		t.Error(err.Error())
	}
	// SearchVolumeId
	volInfoDetached := db.Search(test.volId)
	if *volInfoDetached.Status != cloud.DiskStatusAvailable {
		t.Errorf("name %s: expect %s but actually %s", test.name, cloud.DiskStatusAvailable, *volInfoDetached.Status)
	}
	// 6 DeleteVolume
	err = db.Delete(test.volId)
	if err != nil {
		t.Logf(err.Error())
	}
	// SearchVolumeId
	volInfoDeleted := db.Search(test.volId)
	if *volInfoDeleted.Status != cloud.DiskStatusDeleted {
		t.Logf("name %s: expect %s but actually %s", test.name, cloud.DiskStatusDeleted, *volInfoDeleted.Status)
	}

}
