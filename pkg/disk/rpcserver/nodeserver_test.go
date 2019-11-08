/*
Copyright (C) 2019 Yunify, Inc.

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

package rpcserver

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/yunify/qingcloud-csi/pkg/cloud"
	"github.com/yunify/qingcloud-csi/pkg/cloud/mock"
	"github.com/yunify/qingcloud-csi/pkg/common"
	"github.com/yunify/qingcloud-csi/pkg/disk/driver"
	"github.com/yunify/qingcloud-sdk-go/config"
	qcservice "github.com/yunify/qingcloud-sdk-go/service"
	"testing"
)

var ns csi.NodeServer

func init() {
	initNodeServer()
}

func getInstance(insId string, insType int, zone string) *qcservice.Instance {
	insRepl := cloud.DiskReplicaTypeName[2]
	insStatus := cloud.InstanceStatusRunning
	return &qcservice.Instance{
		InstanceID:    &insId,
		InstanceClass: &insType,
		Repl:          &insRepl,
		Status:        &insStatus,
		ZoneID:        &zone,
	}
}

func initNodeServer() {
	insId := "i-123456"
	input := &driver.InitDiskDriverInput{
		Name:      "test.csi.qingcloud.com",
		Version:   "v1.99.0",
		NodeId:    insId,
		MaxVolume: 10,
	}
	driver := &driver.DiskDriver{}
	driver.InitDiskDriver(input)

	cfg := &config.Config{Zone: "fake-zone"}
	cloudManager := &mock.MockCloudManager{}
	cloudManager.SetConfig(cfg)

	cloudManager.SetInstances(map[string]*qcservice.Instance{
		"i-123456": getInstance(insId, 201, cfg.Zone),
	})

	mounter := common.NewSafeMounter()

	ns = NewNodeServer(driver, cloudManager, mounter)
}

func TestNodeServer_NodeGetInfo(t *testing.T) {
	req := &csi.NodeGetInfoRequest{}
	resp, err := ns.NodeGetInfo(context.Background(), req)
	if err != nil {
		t.Errorf("Error %s", err)
	}

	t.Logf("%v", resp.GetAccessibleTopology())
}
