package block

import (
	"github.com/golang/glog"
	"io/ioutil"
	"os"
	"strings"
)

const (
	InstanceFilepath = "/etc/qingcloud/instance-id"
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
	bytes, err := ioutil.ReadFile(InstanceFilepath)
	if err != nil {
		glog.Errorf("Getting current instance-id error: %s", err.Error())
		os.Exit(1)
	}
	instanceIdFromFile = string(bytes[:])
	instanceIdFromFile = strings.Replace(instanceIdFromFile, "\n", "", -1)
	glog.Infof("Getting current instance-id: \"%s\"", instanceIdFromFile)
}

func GetCurrentInstanceId() string {
	if instanceIdFromFile == ""{
		ReadCurrentInstanceId()
	}
	return instanceIdFromFile
}
