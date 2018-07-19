package block

import (
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	InstanceFilePath     = "/etc/qingcloud/instance-id"
	ConfigFilePath       = "/root/config.yaml"
	RetryString          = "please try later"
	Int64_Max            = int64(^uint64(0) >> 1)
	WaitInterval         = 10 * time.Second
	OperationWaitTimeout = 180 * time.Second
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100
)

var instanceIdFromFile string

func CreatePath(persistentStoragePath string) error {
	if _, err := os.Stat(persistentStoragePath); os.IsNotExist(err) {
		if err := os.MkdirAll(persistentStoragePath, os.FileMode(0755)); err != nil {
			return err
		}
	} else {
	}
	return nil
}

func readCurrentInstanceId() {
	bytes, err := ioutil.ReadFile(InstanceFilePath)
	if err != nil {
		glog.Errorf("Getting current instance-id error: %s", err.Error())
		os.Exit(1)
	}
	instanceIdFromFile = string(bytes[:])
	instanceIdFromFile = strings.Replace(instanceIdFromFile, "\n", "", -1)
	glog.Infof("Getting current instance-id: \"%s\"", instanceIdFromFile)
}

func GetCurrentInstanceId() string {
	if len(instanceIdFromFile) == 0 {
		readCurrentInstanceId()
	}
	return instanceIdFromFile
}

func ReadConfigFromFile(filePath string) (*qcconfig.Config, error) {
	config, err := qcconfig.NewDefault()
	if err != nil {
		return nil, err
	}
	if err = config.LoadConfigFromFilepath(filePath); err != nil {
		return nil, err
	}
	return config, nil
}

func ContainsVolumeCapability(accessModes []*csi.VolumeCapability_AccessMode, subCaps *csi.VolumeCapability) bool {
	for _, cap := range accessModes {
		if cap.GetMode() == subCaps.GetAccessMode().GetMode() {
			return true
		}
	}
	return false
}

func ContainsVolumeCapabilities(accessModes []*csi.VolumeCapability_AccessMode, subCaps []*csi.VolumeCapability) bool {
	for _, v := range subCaps {
		if !ContainsVolumeCapability(accessModes, v) {
			return false
		}
	}
	return true
}

func ContainsNodeServiceCapability(nodeCaps []*csi.NodeServiceCapability, subCap csi.NodeServiceCapability_RPC_Type) bool {
	for _, v := range nodeCaps {
		if strings.Contains(v.String(), subCap.String()) {
			return true
		}
	}
	return false
}

func GbToByte(num int) int64 {
	if num < 0 {
		return 0
	}
	return int64(num) * gib
}

func ByteCeilToGb(num int64) int {
	if num <= 0 {
		return 0
	}
	res := num / gib
	if res*gib < num {
		res += 1
	}
	return int(res)
}
