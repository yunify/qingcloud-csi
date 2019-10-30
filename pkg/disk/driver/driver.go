/*
Copyright (C) 2018 Yunify, Inc.

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

package driver

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"k8s.io/klog"
)

type DiskDriver struct {
	name          string
	version       string
	nodeId        string
	maxVolume     int64
	volumeCap     []*csi.VolumeCapability_AccessMode
	controllerCap []*csi.ControllerServiceCapability
	nodeCap       []*csi.NodeServiceCapability
	pluginCap     []*csi.PluginCapability
}

type InitDiskDriverInput struct {
	Name          string
	Version       string
	NodeId        string
	MaxVolume     int64
	VolumeCap     []csi.VolumeCapability_AccessMode_Mode
	ControllerCap []csi.ControllerServiceCapability_RPC_Type
	NodeCap       []csi.NodeServiceCapability_RPC_Type
	PluginCap     []*csi.PluginCapability
}

// GetDiskDriver
// Create disk driver
func GetDiskDriver() *DiskDriver {
	return &DiskDriver{}
}

func (d *DiskDriver) InitDiskDriver(input *InitDiskDriverInput) {
	if input == nil {
		return
	}
	d.name = input.Name
	d.version = input.Version
	// Setup Node Id
	d.nodeId = input.NodeId
	// Setup max volume
	d.maxVolume = input.MaxVolume
	// Setup cap
	d.addVolumeCapabilityAccessModes(input.VolumeCap)
	d.addControllerServiceCapabilities(input.ControllerCap)
	d.addNodeServiceCapabilities(input.NodeCap)
	d.addPluginCapabilities(input.PluginCap)
}

func (d *DiskDriver) addVolumeCapabilityAccessModes(vc []csi.VolumeCapability_AccessMode_Mode) {
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vc {
		klog.V(4).Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}
	d.volumeCap = vca
}

func (d *DiskDriver) addControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	var csc []*csi.ControllerServiceCapability
	for _, c := range cl {
		klog.V(4).Infof("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}
	d.controllerCap = csc
}

func (d *DiskDriver) addNodeServiceCapabilities(nl []csi.NodeServiceCapability_RPC_Type) {
	var nsc []*csi.NodeServiceCapability
	for _, n := range nl {
		klog.V(4).Infof("Enabling node service capability: %v", n.String())
		nsc = append(nsc, NewNodeServiceCapability(n))
	}
	d.nodeCap = nsc
}

func (d *DiskDriver) addPluginCapabilities(cap []*csi.PluginCapability) {
	d.pluginCap = cap
}

func (d *DiskDriver) ValidateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) bool {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return true
	}

	for _, cap := range d.controllerCap {
		if c == cap.GetRpc().Type {
			return true
		}
	}
	return false
}

func (d *DiskDriver) ValidateNodeServiceRequest(c csi.NodeServiceCapability_RPC_Type) bool {
	if c == csi.NodeServiceCapability_RPC_UNKNOWN {
		return true
	}
	for _, cap := range d.nodeCap {
		if c == cap.GetRpc().Type {
			return true
		}
	}
	return false

}

func (d *DiskDriver) ValidateVolumeCapability(cap *csi.VolumeCapability) bool {
	if !d.ValidateVolumeAccessMode(cap.GetAccessMode().GetMode()) {
		return false
	}
	return true
}

func (d *DiskDriver) ValidateVolumeCapabilities(caps []*csi.VolumeCapability) bool {
	for _, cap := range caps {
		if !d.ValidateVolumeAccessMode(cap.GetAccessMode().GetMode()) {
			return false
		}
	}
	return true
}

func (d *DiskDriver) ValidateVolumeAccessMode(c csi.VolumeCapability_AccessMode_Mode) bool {
	for _, mode := range d.volumeCap {
		if c == mode.GetMode() {
			return true
		}
	}
	return false
}

func (d *DiskDriver) ValidatePluginCapabilityService(cap csi.PluginCapability_Service_Type) bool {
	for _, v := range d.GetPluginCapability() {
		if v.GetService() != nil && v.GetService().GetType() == cap {
			return true
		}
	}
	return false
}

func (d *DiskDriver) GetName() string {
	return d.name
}

func (d *DiskDriver) GetVersion() string {
	return d.version
}

func (d *DiskDriver) GetInstanceId() string {
	return d.nodeId
}

func (d *DiskDriver) GetMaxVolumePerNode() int64 {
	return d.maxVolume
}

func (d *DiskDriver) GetControllerCapability() []*csi.ControllerServiceCapability {
	return d.controllerCap
}

func (d *DiskDriver) GetNodeCapability() []*csi.NodeServiceCapability {
	return d.nodeCap
}

func (d *DiskDriver) GetPluginCapability() []*csi.PluginCapability {
	return d.pluginCap
}

func (d *DiskDriver) GetVolumeCapability() []*csi.VolumeCapability_AccessMode {
	return d.volumeCap
}

func (d *DiskDriver) GetTopologyZoneKey() string {
	return fmt.Sprintf("topology.%s/zone", d.GetName())
}

func (d *DiskDriver) GetTopologyInstanceTypeKey() string {
	return fmt.Sprintf("topology.%s/instance-type", d.GetName())
}
