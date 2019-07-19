package cloudprovider

import (
	"errors"
	"fmt"
	qcclient "github.com/yunify/qingcloud-sdk-go/client"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"k8s.io/klog"
)

var _ CloudManager = &cloudManager{}

type CloudManager interface {
	// Snapshot Method
	FindSnapshot(snapId string) (snapInfo *qcservice.Snapshot, err error)
	FindSnapshotByName(snapName string) (snapInfo *qcservice.Snapshot, err error)
	CreateSnapshot(snapName string, volId string) (snapId string, err error)
	DeleteSnapshot(snapId string) (err error)
	// Volume Method
	FindVolume(volId string) (volInfo *qcservice.Volume, err error)
	FindVolumeByName(volName string) (volInfo *qcservice.Volume, err error)
	CreateVolume(volName string, requestSize int, repl int, volType int) (volId string, err error)
	CreateVolumeFromSnapshot(volName string, snapId string) (volId string, err error)
	DeleteVolume(volId string) (err error)
	AttachVolume(volId string, instanceId string) (err error)
	DetachVolume(volId string, instanceId string) (err error)
	ResizeVolume(volId string, requestSize int) (err error)
	// Util Method
	FindInstance(instanceId string) (instanceInfo *qcservice.Instance, err error)
	GetZone() (zoneName string)
	GetZoneList() (zoneNameList []string, err error)
	waitJob(jobId string) (err error)
}

type cloudManager struct {
	instanceService *qcservice.InstanceService
	snapshotService *qcservice.SnapshotService
	volumeService   *qcservice.VolumeService
	jobService      *qcservice.JobService
	cloudService    *qcservice.QingCloudService
}

func NewCloudManager(config *qcconfig.Config) (CloudManager, error) {
	// initial qingcloud iaas service
	qs, err := qcservice.Init(config)
	if err != nil {
		return nil, err
	}
	// create services
	is, _ := qs.Instance(config.Zone)
	ss, _ := qs.Snapshot(config.Zone)
	vs, _ := qs.Volume(config.Zone)
	js, _ := qs.Job(config.Zone)

	// initial cloud manager
	cm := cloudManager{
		instanceService: is,
		snapshotService: ss,
		volumeService:   vs,
		jobService:      js,
		cloudService:    qs,
	}
	klog.Infof("Succeed to initial cloud manager")
	return &cm, nil
}

// NewCloudManagerFromFile
// Create cloud manager from file
func NewCloudManagerFromFile(filePath string) (CloudManager, error) {
	// create config
	config, err := ReadConfigFromFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewCloudManager(config)
}

// Find snapshot by snapshot id
// Return: 	nil,	nil: 	not found snapshot
//			snapshot, nil: 	found snapshot
//			nil, 	error:	internal error
func (cm *cloudManager) FindSnapshot(id string) (snapshot *qcservice.Snapshot, err error) {
	verboseMode := EnableDescribeSnapshotVerboseMode
	// Set DescribeSnapshot input
	input := &qcservice.DescribeSnapshotsInput{
		Snapshots: []*string{&id},
		Verbose:   &verboseMode,
	}
	// Call describe snapshot
	output, err := cm.snapshotService.DescribeSnapshots(input)
	// 1. Error is not equal to nil.
	if err != nil {
		return nil, err
	}
	// 2. Return code is not equal to 0.
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf("call IaaS DescribeSnapshot err: snapshot id %s in %s",
			id, cm.snapshotService.Config.Zone)
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
			fmt.Errorf("call IaaS DescribeSnapshot err: find duplicate snapshot, snapshot id %s in %s",
				id, cm.snapshotService.Config.Zone)
	}
}

// Find snapshot by snapshot name
// In Qingcloud IaaS platform, it is possible that two snapshots have the same name.
// In Kubernetes, the CO will set a unique PV name.
// CSI driver take the PV name as a snapshot name.
// Return: 	nil, 		nil: 	not found snapshots
//			snapshots,	nil:	found snapshot
//			nil,		error:	internal error
func (cm *cloudManager) FindSnapshotByName(name string) (snapshot *qcservice.Snapshot, err error) {
	if len(name) == 0 {
		return nil, nil
	}
	verboseMode := EnableDescribeSnapshotVerboseMode
	// Set input arguments
	input := &qcservice.DescribeSnapshotsInput{
		SearchWord: &name,
		Verbose:    &verboseMode,
	}
	// Call DescribeSnapshot
	output, err := cm.snapshotService.DescribeSnapshots(input)
	// Handle error
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf("call IaaS DescribeSnapshots err: snapshot name %s in %s",
			name, cm.snapshotService.Config.Zone)
	}
	// Not found snapshots
	for _, v := range output.SnapshotSet {
		if *v.SnapshotName == name && *v.Status != SnapshotStatusCeased && *v.Status != SnapshotStatusDeleted {
			return v, nil
		}
	}
	return nil, nil
}

// CreateSnapshot
// 1. format snapshot size
// 2. create snapshot
// 3. wait job
func (cm *cloudManager) CreateSnapshot(snapshotName string, resourceId string) (snapshotId string, err error) {
	// 0. Set CreateSnapshot args
	isFull := int(SnapshotFull)
	// set input value
	input := &qcservice.CreateSnapshotsInput{
		SnapshotName: &snapshotName,
		IsFull:       &isFull,
		Resources:    []*string{&resourceId},
	}

	// 1. Create snapshot
	klog.Infof("Call IaaS CreateSnapshot request snapshot name: %s, zone: %s, resource id %s, is full snapshot %T",
		*input.SnapshotName, cm.GetZone(), *input.Resources[0], *input.IsFull == SnapshotFull)
	output, err := cm.snapshotService.CreateSnapshots(input)
	if err != nil {
		return "", err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return "", fmt.Errorf("call IaaS CreateSnapshot error: %s", *output.Message)
	}
	snapshotId = *output.Snapshots[0]
	klog.Infof("Call IaaS CreateSnapshots snapshot name %s snapshot id %s succeed", snapshotName, snapshotId)
	return snapshotId, nil
}

// DeleteSnapshot
// 1. delete snapshot by id
// 2. wait job
func (sm *cloudManager) DeleteSnapshot(snapshotId string) error {
	// set input value
	input := &qcservice.DeleteSnapshotsInput{
		Snapshots: []*string{&snapshotId},
	}
	// delete snapshot
	klog.Infof("Call IaaS DeleteSnapshot request id: %s, zone: %s",
		snapshotId, *sm.snapshotService.Properties.Zone)
	output, err := sm.snapshotService.DeleteSnapshots(input)
	if err != nil {
		return err
	}
	// wait job
	klog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := sm.waitJob(*output.JobID); err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	klog.Infof("Call IaaS DeleteSnapshot %s succeed", snapshotId)
	return nil
}

// Find volume by volume ID
// Return: 	nil,	nil: 	not found volumes
//			volume, nil: 	found volume
//			nil, 	error:	internal error
func (cm *cloudManager) FindVolume(id string) (volInfo *qcservice.Volume, err error) {
	// Set DescribeVolumes input
	input := &qcservice.DescribeVolumesInput{
		Volumes: []*string{&id},
	}
	// Call describe volume
	output, err := cm.volumeService.DescribeVolumes(input)
	// Error:
	// 1. Error is not equal to nil.
	if err != nil {
		return nil, err
	}
	// 2. Return code is not equal to 0.
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil,
			fmt.Errorf("call IaaS DescribeVolumes err: volume id %s in %s", id, cm.volumeService.Config.Zone)
	}
	switch *output.TotalCount {
	// Not found volumes
	case 0:
		return nil, nil
	// Found one volume
	case 1:
		if *output.VolumeSet[0].Status == DiskStatusCeased || *output.VolumeSet[0].
			Status == DiskStatusDeleted {
			return nil, nil
		}
		return output.VolumeSet[0], nil
	// Found duplicate volumes
	default:
		return nil,
			fmt.Errorf("call IaaS DescribeVolumes err: find duplicate volumes, volume id %s in %s",
				id, cm.volumeService.Config.Zone)
	}
}

// Find volume by volume name
// In Qingcloud IaaS platform, it is possible that two volumes have the same name.
// In Kubernetes, the CO will set a unique PV name.
// CSI driver take the PV name as a volume name.
// Return: 	nil, 		nil: 	not found volumes
//			volumes,	nil:	found volume
//			nil,		error:	internal error
func (cm *cloudManager) FindVolumeByName(name string) (volume *qcservice.Volume, err error) {
	if len(name) == 0 {
		return nil, nil
	}
	// Set input arguments
	input := &qcservice.DescribeVolumesInput{
		SearchWord: &name,
	}
	// Call DescribeVolumes
	output, err := cm.volumeService.DescribeVolumes(input)
	// Handle error
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return nil, fmt.Errorf("call IaaS DescribeVolumes err: volume name %s in %s",
			name, cm.volumeService.Config.Zone)
	}
	// Not found volumes
	for _, v := range output.VolumeSet {
		if *v.VolumeName != name {
			continue
		}
		if *v.Status == DiskStatusCeased || *v.Status == DiskStatusDeleted {
			continue
		}
		return v, nil
	}
	return nil, nil
}

// CreateVolume
// 1. format volume size
// 2. create volume
// 3. wait job
func (cm *cloudManager) CreateVolume(volumeName string, requestSize int, replica int, volType int) (volumeId string,
	err error) {
	// 0. Set CreateVolume args
	// create volume count
	count := 1
	// volume replicas
	replStr := DiskReplicaTypeName[replica]
	// set input value
	input := &qcservice.CreateVolumesInput{
		Count:      &count,
		Repl:       &replStr,
		Size:       &requestSize,
		VolumeName: &volumeName,
		VolumeType: &volType,
	}
	// 1. Create volume
	klog.Infof("Call IaaS CreateVolume request size: %d GB, zone: %s, type: %d, count: %d, replica: %s, name: %s",
		*input.Size, *cm.volumeService.Properties.Zone, *input.VolumeType, *input.Count, *input.Repl, *input.VolumeName)
	output, err := cm.volumeService.CreateVolumes(input)
	if err != nil {
		return "", err
	}
	// wait job
	klog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := cm.waitJob(*output.JobID); err != nil {
		return "", err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return "", fmt.Errorf(*output.Message)
	}
	volumeId = *output.Volumes[0]
	klog.Infof("Call IaaS CreateVolume name %s id %s succeed", volumeName, volumeId)
	return *output.Volumes[0], nil
}

// CreateVolumeFromSnapshot
// In QingCloud, the volume size created from snapshot is equal to original volume.
func (cm *cloudManager) CreateVolumeFromSnapshot(volumeName string, snapshotId string) (volumeId string, err error) {
	input := &qcservice.CreateVolumeFromSnapshotInput{
		VolumeName: &volumeName,
		Snapshot:   &snapshotId,
	}
	klog.Infof("Call IaaS CreateVolumeFromSnapshot request volume name: %s, snapshot id: %s\n",
		*input.VolumeName, *input.Snapshot)
	output, err := cm.snapshotService.CreateVolumeFromSnapshot(input)
	if err != nil {
		return "", err
	}
	// wait job
	klog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := cm.waitJob(*output.JobID); err != nil {
		return "", err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return "", fmt.Errorf(*output.Message)
	}
	klog.Infof("Call IaaS CreateVolumeFromSnapshot succeed, volume id %s", *output.VolumeID)
	return *output.VolumeID, nil
}

// DeleteVolume
// 1. delete volume by id
// 2. wait job
func (cm *cloudManager) DeleteVolume(id string) error {
	// set input value
	input := &qcservice.DeleteVolumesInput{
		Volumes: []*string{&id},
	}
	// delete volume
	klog.Infof("Call IaaS DeleteVolume request id: %s, zone: %s",
		id, *cm.volumeService.Properties.Zone)
	output, err := cm.volumeService.DeleteVolumes(input)
	if err != nil {
		return err
	}
	// wait job
	klog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := cm.waitJob(*output.JobID); err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	klog.Infof("Call IaaS DeleteVolume %s succeed", id)
	return nil
}

// AttachVolume
// 1. attach volume on instance
// 2. wait job
func (cm *cloudManager) AttachVolume(volumeId string, instanceId string) error {
	// set input parameter
	input := &qcservice.AttachVolumesInput{
		Volumes:  []*string{&volumeId},
		Instance: &instanceId,
	}
	// attach volume
	klog.Infof("Call IaaS AttachVolume request volume id: %s, instance id: %s, zone: %s", volumeId, instanceId,
		cm.GetZone())
	output, err := cm.volumeService.AttachVolumes(input)
	if err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	// wait job
	klog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := cm.waitJob(*output.JobID); err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	klog.Infof("Call IaaS AttachVolume %s on instance %s succeed", volumeId, instanceId)
	return nil
}

// detach volume
// 1. detach volume
// 2. wait job
func (cm *cloudManager) DetachVolume(volumeId string, instanceId string) error {
	// set input parameter
	input := &qcservice.DetachVolumesInput{
		Volumes:  []*string{&volumeId},
		Instance: &instanceId,
	}
	// detach volume
	klog.Infof("Call IaaS DetachVolume request volume id: %s, instance id: %s, zone: %s", volumeId,
		instanceId, cm.GetZone())
	output, err := cm.volumeService.DetachVolumes(input)
	if err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	// wait job
	klog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := cm.waitJob(*output.JobID); err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return fmt.Errorf(*output.Message)
	}
	klog.Infof("Call IaaS DetachVolume %s succeed", volumeId)
	return nil
}

// ResizeVolume can expand the size of a volume offline
// requestSize: GB
func (cm *cloudManager) ResizeVolume(volumeId string, requestSize int) error {
	// resize
	klog.Infof("Call IaaS ResizeVolume request volume %s size %d Gib in zone [%s]",
		volumeId, requestSize, cm.GetZone())
	input := &qcservice.ResizeVolumesInput{
		Size:    &requestSize,
		Volumes: []*string{&volumeId},
	}
	output, err := cm.volumeService.ResizeVolumes(input)
	if err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return errors.New(*output.Message)
	}
	// wait job
	klog.Infof("Call IaaS WaitJob %s", *output.JobID)
	if err := cm.waitJob(*output.JobID); err != nil {
		return err
	}
	// check output
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
		return errors.New(*output.Message)
	}
	klog.Infof("Call IaaS ResizeVolume id %s size %d succeed", volumeId, requestSize)
	return nil
}

// Find instance by instance ID
// Return: 	nil,	nil: 	not found instance
//			instance, nil: 	found instance
//			nil, 	error:	internal error
func (cm *cloudManager) FindInstance(id string) (instance *qcservice.Instance, err error) {
	seeCluster := EnableDescribeInstanceAppCluster
	// set describe instance input
	input := qcservice.DescribeInstancesInput{
		Instances:     []*string{&id},
		IsClusterNode: &seeCluster,
	}
	// call describe instance
	output, err := cm.instanceService.DescribeInstances(&input)
	// error
	if err != nil {
		return nil, err
	}
	if *output.RetCode != 0 {
		klog.Errorf("Ret code: %d, message: %s", *output.RetCode, *output.Message)
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
		return nil, fmt.Errorf("find duplicate instances id %s in %s", id, cm.instanceService.Config.Zone)
	}
}

// GetZone
// Get current zone in Qingcloud IaaS
func (cm *cloudManager) GetZone() string {
	if cm == nil {
		return ""
	}
	return cm.cloudService.Config.Zone
}

// GetZoneList gets active zone list
func (zm *cloudManager) GetZoneList() (zones []string, err error) {
	output, err := zm.cloudService.DescribeZones(&qcservice.DescribeZonesInput{})
	// Error:
	// 1. Error is not equal to nil.
	if err != nil {
		return nil, err
	}
	if output == nil {
		klog.Error("should not response nil")
	}
	for i := range output.ZoneSet {
		if *output.ZoneSet[i].Status == ZoneStatusActive {
			zones = append(zones, *output.ZoneSet[i].ZoneID)
		}
	}
	return zones, nil
}

func (cm *cloudManager) waitJob(jobId string) error {
	err := qcclient.WaitJob(cm.jobService, jobId, WaitJobTimeout, WaitJobInterval)
	if err != nil {
		return fmt.Errorf("call IaaS WaitJob id %s, error: ", err)
	}
	return nil
}
