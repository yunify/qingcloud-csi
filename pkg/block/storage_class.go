package block

import (
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
)

type qingStorageClass struct{
	AccessKeyId string  `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Zone string `json:"zone"`
	Type string `json:"type"`
	Host string `json:"host"`
	Port int `json:"port"`
	Protocol string `json:"protocol"`
}

func getConfigFromStorageClass(sc *qingStorageClass)(config *qcconfig.Config){
	config,err := qcconfig.NewDefault()
	if err != nil{
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