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

package mock

import (
	"fmt"
	"github.com/yunify/qingcloud-csi/pkg/cloud"
	"github.com/yunify/qingcloud-csi/pkg/common"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"time"
)

var _ cloud.CloudManager = &mockCloudManager{}

type mockCloudManager struct {
	qcconfig  *qcconfig.Config
	snapshots map[string]*qcservice.Snapshot
	volumes   map[string]*qcservice.Volume
}

func NewMockCloudManagerFromConfig(config *qcconfig.Config) (cloud.CloudManager, error) {
	return &mockCloudManager{
		qcconfig: config,
	}, nil
}

func (m *mockCloudManager) FindSnapshot(snapId string) (snapInfo *qcservice.Snapshot, err error) {
	for _, v := range m.snapshots {
		if *v.SnapshotID == snapId {
			return v, nil
		}
	}
	return nil, nil
}

func (m *mockCloudManager) FindSnapshotByName(snapName string) (snapInfo *qcservice.Snapshot, err error) {
	for _, v := range m.snapshots {
		if *v.SnapshotName == snapName {
			return v, nil
		}
	}
	return nil, nil
}
func (m *mockCloudManager) CreateSnapshot(snapName string, volId string) (snapId string, err error) {
	volInfo, err := m.FindVolume(volId)
	if err != nil {
		return "", err
	}
	if volInfo == nil {
		return "", fmt.Errorf("create snapshot %s error: volume %s does not exist", snapName, volId)
	}
	snapStatus := string(cloud.SnapshotStatusAvailable)
	snapId = common.GenerateHashInEightBytes(snapName + volId + time.Now().UTC().String())
	snapEntity := &qcservice.Snapshot{
		SnapshotID:   &snapId,
		SnapshotName: &snapName,
		SnapshotResource: &qcservice.SnapshotResource{
			VolumeID:   &volId,
			VolumeType: volInfo.VolumeType,
			Size:       volInfo.Size,
		},
		Status: &snapStatus,
	}
	m.snapshots[snapId] = snapEntity
	return snapId, nil
}

func (m *mockCloudManager) DeleteSnapshot(snapId string) (err error) {
	snapInfo, err := m.FindSnapshot(snapId)
	if err != nil {
		return err
	}
	if snapInfo == nil || *snapInfo.Status == cloud.SnapshotStatusDeleted {
		return fmt.Errorf("delete snapshot %s error: snapshot has been deleted", snapId)
	}

	return nil
}

func (m *mockCloudManager) CreateVolumeFromSnapshot(volName string, snapId string, zone string) (volId string,
	err error) {
	exVol, err := m.FindVolumeByName(volName)
	if err != nil {
		return "", err
	}
	if exVol != nil {
		return "", fmt.Errorf("create volume error: volume %s already exist", volName)
	}

	return "", nil
}

// Volume Method
func (m *mockCloudManager) FindVolume(volId string) (volInfo *qcservice.Volume, err error) {
	info, ok := m.volumes[volId]
	if !ok {
		return nil, nil
	}
	switch *info.Status {
	case cloud.DiskStatusDeleted:
		fallthrough
	case cloud.DiskStatusCeased:
		return nil, nil
	default:
		return info, nil
	}
}

func (m *mockCloudManager) FindVolumeByName(volName string) (volInfo *qcservice.Volume, err error) {
	for _, v := range m.volumes {
		if *v.VolumeName == volName {
			switch *v.Status {
			case cloud.DiskStatusDeleted:
				fallthrough
			case cloud.DiskStatusCeased:
				continue
			default:
				return v, nil
			}
		}
	}
	return nil, nil
}

func (m *mockCloudManager) CreateVolume(volName string, requestSize int, replicas int, volType int, zone string) (
	volId string, err error) {
	exVol, err := m.FindVolumeByName(volName)
	if err != nil {
		return "", err
	}
	if exVol != nil {
		return "", fmt.Errorf("create volume error: volume %s already exist", volName)
	}
	volId = "vol-" + common.GenerateHashInEightBytes(volName+time.Now().UTC().String())
	replStr := cloud.DiskReplicaTypeName[replicas]
	status := cloud.DiskStatusAvailable
	vol := &qcservice.Volume{
		VolumeID:   &volId,
		VolumeName: &volName,
		VolumeType: &volType,
		Size:       &requestSize,
		Repl:       &replStr,
		Status:     &status,
	}
	m.volumes[volId] = vol
	return volId, nil
}

func (m *mockCloudManager) DeleteVolume(volId string) (err error) {
	exVol, err := m.FindVolume(volId)
	if err != nil {
		return err
	}
	if exVol == nil {
		return fmt.Errorf("delete volume error: volume %s does not exist", volId)
	}
	status := cloud.DiskStatusDeleted
	exVol.Status = &status
	m.volumes[volId] = exVol
	return nil
}

func (m *mockCloudManager) AttachVolume(volId string, instanceId string) (err error) {

	return nil
}

func (m *mockCloudManager) DetachVolume(volId string, instanceId string) (err error) {
	return nil
}

func (m *mockCloudManager) ResizeVolume(volId string, requestSize int) (err error) {
	return nil
}

func (m *mockCloudManager) CloneVolume(volName string, volType int, srcVolId string, zone string) (newVolId string, err error) {
	return "", nil
}

// Util Method
func (m *mockCloudManager) FindInstance(instanceId string) (instanceInfo *qcservice.Instance, err error) {
	return nil, nil
}
func (m *mockCloudManager) GetZone() (zoneName string) {
	return ""
}
func (m *mockCloudManager) GetZoneList() (zoneNameList []string, err error) {
	return nil, nil
}
func (m *mockCloudManager) waitJob(jobId string) (err error) {
	return nil
}

// FindTags finds and gets tags information
func (m *mockCloudManager) FindTag(tagId string) (tagInfo *qcservice.Tag, err error) {
	return nil, nil
}

// IsValidTags checks tags exists.
func (m *mockCloudManager) IsValidTags(tagsId []string) bool {
	return false
}

// AttachTags add a slice of tags on a object
func (m *mockCloudManager) AttachTags(tagsId []string, resourceId string, resourceType string) (err error) {
	return nil
}
