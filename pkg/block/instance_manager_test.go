package block

import (
	"runtime"
	"testing"
)

var getim = func() InstanceManager {
	// get storage class
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = "../../ut-config.yaml"
	}
	if runtime.GOOS == "linux" {
		filepath = "../../ut-config.yaml"
	}
	if runtime.GOOS == "darwin" {
		filepath = "../../ut-config.yaml"
	}
	qcConfig, err := ReadConfigFromFile(filepath)
	if err != nil {
		return nil
	}
	im, err := NewInstanceManagerWithConfig(qcConfig)
	if err != nil {
		return nil
	}

	return im
}

func TestFindInstance(t *testing.T) {
	im := getim()
	testcases := []struct {
		name  string
		id    string
		found bool
	}{
		{
			name:  "Avaiable",
			id:    instanceId1,
			found: true,
		},
		{
			name:  "Not found",
			id:    "instance-1234",
			found: false,
		},
		{
			name:  "By name",
			id:    "neonsan-test",
			found: false,
		},
	}
	for _, v := range testcases {
		ins, err := im.FindInstance(v.id)
		if err != nil {
			t.Errorf("name %s error: %s", v.name, err.Error())
		}
		if v.found && *ins.InstanceID != v.id {
			t.Errorf("name %s: find id error")
		}
	}
}
