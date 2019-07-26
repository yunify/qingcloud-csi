package mock

import (
	"fmt"
	"github.com/yunify/qingcloud-csi/pkg/cloudprovider"
	"github.com/yunify/qingcloud-csi/pkg/common"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"k8s.io/klog"
	"sync"
	"time"
)

type mockStorageSystem struct {
}

type mockVolumeDB struct {
	volume map[string]*qcservice.Volume
	sync.Mutex
}

func NewMockVolumeDB() *mockVolumeDB {
	volMp := make(map[string]*qcservice.Volume)
	return &mockVolumeDB{
		volume: volMp,
	}
}

func (m mockVolumeDB) Create(item *qcservice.Volume) (string, error) {
	if item == nil {
		return "", fmt.Errorf("try to add nil item")
	}
	if item.VolumeName == nil || item.Repl == nil || item.Size == nil || item.VolumeType == nil {
		return "", fmt.Errorf("create volume error: lack of input args")
	}
	volId := "vol-" + common.GenerateHashInEightBytes(*item.VolumeName+time.Now().UTC().String())
	status := cloudprovider.DiskStatusAvailable
	item.Status = &status
	m.Lock()
	defer m.Unlock()
	m.volume[volId] = item
	klog.Infof("succeed to update volume %s", volId)
	return volId, nil
}

func (m mockVolumeDB) Delete(volId string) error {
	if m.Search(volId) == nil {
		return fmt.Errorf("delete volume error: volume %s does not exist", volId)
	}
	if *m.volume[volId].Status != cloudprovider.DiskStatusAvailable {
		return fmt.Errorf("delete volume error: volume %s does not available", volId)
	}
	status := cloudprovider.DiskStatusDeleted
	m.Lock()
	m.volume[volId].Status = &status
	defer m.Unlock()
	klog.Infof("succeed to delete volume %s", volId)
	return nil
}

func (m mockVolumeDB) Search(volId string) *qcservice.Volume {
	return m.volume[volId]
}

func (m mockVolumeDB) SearchName(name string) []*qcservice.Volume {
	res := []*qcservice.Volume{}
	for k, v := range m.volume {
		if *v.VolumeName == name {
			res = append(res, m.volume[k])
		}
	}
	return res
}

func (m mockVolumeDB) Attach(volumeId string, instanceId string) error {
	volInfo := m.Search(volumeId)
	if volInfo == nil {
		return fmt.Errorf("attach volume error: volume %s does not exist", volumeId)
	}
	if *volInfo.Status == cloudprovider.DiskStatusInuse {
		return fmt.Errorf("attach volume error: volume %s already attached", volumeId)
	}
	status := cloudprovider.DiskStatusInuse
	m.Lock()
	volInfo.Instance = &qcservice.Instance{
		InstanceID: &instanceId,
	}
	volInfo.Status = &status
	defer m.Unlock()
	klog.Infof("succeed to attach volume %s", volumeId)
	return nil
}

func (m mockVolumeDB) Detach(volumeId string, instanceId string) error {
	volInfo := m.Search(volumeId)
	if volInfo == nil {
		return fmt.Errorf("attach volume error: volume %s does not exist", volumeId)
	}
	if *volInfo.Status != cloudprovider.DiskStatusInuse {
		return fmt.Errorf("attach volume error: volume %s does not attached", volumeId)
	}
	if *volInfo.Instance.InstanceID != instanceId {
		return fmt.Errorf("attach volume error: volume %s has been attached to another instance", volumeId)
	}
	status := cloudprovider.DiskStatusAvailable
	m.Lock()
	volInfo.Instance.InstanceID = nil
	volInfo.Status = &status
	defer m.Unlock()
	klog.Infof("succeed to attach volume %s", volumeId)
	return nil
}
