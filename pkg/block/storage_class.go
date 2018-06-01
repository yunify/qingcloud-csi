package block

import (
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
)

type qingStorageClass struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Zone            string `json:"zone"`
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Protocol        string `json:"protocol"`
	VolumeType      int    `json:"type"`
	VolumeMaxSize   int    `json:"maxSize"`
	VolumeMinSize   int    `json:"minSize"`
}

func NewDefaultQingStorageClass() qingStorageClass {
	return qingStorageClass{
		AccessKeyId:     "KEY_ID",
		AccessKeySecret: "KEY_SECRET",
		Zone:            "sh1a",
		Host:            "api.qingcloud.com",
		Port:            443,
		Protocol:        "https",
		VolumeType:      0,
		VolumeMaxSize:   500,
		VolumeMinSize:   10,
	}
}

func (sc qingStorageClass) formatVolumeSize(size int) int {
	if size <= sc.VolumeMinSize {
		return sc.VolumeMinSize
	} else if size >= sc.VolumeMaxSize {
		return sc.VolumeMaxSize
	}
	if size%10 != 0 {
		size = (size/10 + 1) * 10
	}
	return size
}

func (sc qingStorageClass) getConfig() (config *qcconfig.Config) {
	config, err := qcconfig.NewDefault()
	if err != nil {
		glog.Error(err)
		return nil
	}
	config.AccessKeyID = sc.AccessKeyId
	config.SecretAccessKey = sc.AccessKeySecret
	config.Zone = sc.Zone
	config.Host = sc.Host
	config.Port = sc.Port
	config.Protocol = sc.Protocol
	return config
}
