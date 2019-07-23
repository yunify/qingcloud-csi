package driver

type Topology struct {
	zone         string
	instanceType InstanceType
}

func NewTopology(zone string, instanceType InstanceType) *Topology {
	return &Topology{zone, instanceType}
}

func (t *Topology) GetZone() string {
	return t.zone
}

func (t *Topology) GetInstanceType() InstanceType {
	return t.instanceType
}

func (t *Topology) SetZone(zone string) {
	t.zone = zone
}

func (t *Topology) SetInstanceType(instanceType InstanceType) {
	t.instanceType = instanceType
}
