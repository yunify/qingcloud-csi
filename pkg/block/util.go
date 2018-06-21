package block

import (
	"github.com/golang/glog"
	"io/ioutil"
	"os"
)

const (
	InstanceFilepath = "/etc/qingcloud/instance-id"
)

var instanceID string

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
	bytes, err := ioutil.ReadFile(InstanceFilepath)
	if err != nil {
		glog.Errorf("Get instance id error: %s", err.Error())
		os.Exit(1)
	}
	instanceID = string(bytes[:])
	glog.Infof("Current instance id is \"%s\"", instanceID)
}

func GetCurrentInstanceId() string {
	return instanceID
}
