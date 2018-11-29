package zone

import (
	"github.com/golang/glog"
	"github.com/yunify/qingcloud-csi/pkg/server"
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
)

const (
	ZoneStatusActive  = "active"
	ZoneStatusFaulty  = "faulty"
	ZoneStatusDefunct = "defunct"
)

type ZoneManager interface {
	GetZoneList() ([]string, error)
}

type zoneManager struct {
	zoneService *qcservice.QingCloudService
}

// NewZoneManagerFromConfig
// Create zone manager from config
func NewZoneManagerFromConfig(config *qcconfig.Config) (ZoneManager, error) {
	// initial qingcloud iaas service
	qs, err := qcservice.Init(config)
	if err != nil {
		return nil, err
	}
	// initial zone provisioner
	zm := zoneManager{
		zoneService: qs,
	}
	glog.Infof("Finished initial zone manager")
	return &zm, nil
}

// NewZoneManagerFromFile
// Create zone manager from file
func NewZoneManagerFromFile(filePath string) (ZoneManager, error) {
	config, err := server.ReadConfigFromFile(filePath)
	if err != nil {
		glog.Errorf("Failed read config file [%s], error: [%s]", filePath, err.Error())
		return nil, err
	}
	glog.Infof("Succeed read config file [%s]", filePath)
	return NewZoneManagerFromConfig(config)
}

// GetZoneList gets active zone list
func (zm *zoneManager) GetZoneList() (zones []string, err error) {
	output, err := zm.zoneService.DescribeZones(&qcservice.DescribeZonesInput{})
	// Error:
	// 1. Error is not equal to nil.
	if err != nil {
		return nil, err
	}
	if output == nil {
		glog.Errorf("should not response [%#v]", output)
	}
	for i := range output.ZoneSet {
		if *output.ZoneSet[i].Status == ZoneStatusActive {
			zones = append(zones, *output.ZoneSet[i].ZoneID)
		}
	}
	return zones, nil
}
