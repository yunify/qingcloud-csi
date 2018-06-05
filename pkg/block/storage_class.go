package block

import (
	"github.com/golang/glog"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"fmt"
	"strconv"
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

func NewDefaultQingStorageClass() *qingStorageClass {
	return &qingStorageClass{
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

func NewStorageClassFromMap(opt map[string]string)(*qingStorageClass, error){
	var ok bool
	sc := NewDefaultQingStorageClass()
	sc.AccessKeyId ,ok=opt["accessKeyId"]
	if !ok{
		return nil, fmt.Errorf("Missing required parameter accessKeyId")
	}
	sc.AccessKeySecret, ok=opt["accessKeySecret"]
	if !ok{
		return nil, fmt.Errorf("Missing required parameter accessKeySecret")
	}
	sc.Zone, ok = opt["zone"]
	if !ok{
		return nil, fmt.Errorf("Missing required parameter zone")
	}
	sc.Host, ok = opt["host"]
	if !ok{
		return nil, fmt.Errorf("Missing required parameter host")
	}
	// port
	sport, ok := opt["port"]
	if !ok{
		return nil, fmt.Errorf("Missing required parameter port")
	}
	iport, err := strconv.Atoi(sport)
	if err != nil{
		return nil,err
	}else{
		sc.Port = iport
	}
	// protocol
	sc.Protocol, ok = opt["protocol"]
	if !ok{
		return nil, fmt.Errorf("Missing required parameter protocol")
	}
	// volume type
	sVolType, ok := opt["type"]
	if !ok{
		return nil, fmt.Errorf("Missing required parameter type")
	}
	iVolType, err := strconv.Atoi(sVolType)
	if err != nil{
		return nil,err
	}else{
		sc.VolumeType = iVolType
	}
	// Get volume maxsize +optional
	sMaxSize, ok := opt["maxSize"]
	iMaxSize, err := strconv.Atoi(sMaxSize)
	if err != nil{
		return nil,err
	}else{
		sc.VolumeMaxSize = iMaxSize
	}
	// Get volume minsize +optional
	sMinSize, ok := opt["minSize"]
	iMinSize, err := strconv.Atoi(sMinSize)
	if err != nil{
		return nil,err
	}else{
		sc.VolumeMinSize = iMinSize
	}
	// Ensure volume minSize less than volume maxSize
	if sc.VolumeMinSize >= sc.VolumeMaxSize{
		return nil, fmt.Errorf("Volume minSize must less than volume maxSize")
	}
	return sc,nil
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