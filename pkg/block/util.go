package block

import (
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"github.com/container-storage-interface/spec/lib/go/csi/v0"
)

const (
	InstanceFilePath     = "/etc/qingcloud/instance-id"
	ConfigFilePath       = "/root/config.yaml"
	Int64_Max            = int64(^uint64(0) >> 1)
	WaitInterval         = 10 * time.Second
	OperationWaitTimeout = 180 * time.Second
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

func ReadCurrentInstanceId() {
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
	if instanceIdFromFile == "" {
		ReadCurrentInstanceId()
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

func HasSameAccessMode(accessMode []*csi.VolumeCapability_AccessMode, cap []*csi.VolumeCapability)bool{
	for _, c := range cap {
		found := false
		for _, c1 := range accessMode{
			if c1.GetMode() == c.GetAccessMode().GetMode() {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func GbToByte(num int) int64{
	if num <0{
		return 0
	}
	return int64(num)*gib
}

func ByteCeilToGb(num int64) int{
	if num <= 0{
		return 0
	}
	res := num/gib
	if res *gib < num{
		res +=1
	}
	return int(res)
}

func GbGreatThanByte(gb int, byte int64) int{
	gb_int64 := GbToByte(gb)
	if gb_int64 < byte{
		return -1
	}else if gb_int64 == byte{
		return 0
	}else{
		return 1
	}
}