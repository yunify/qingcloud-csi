package block

import (
	"runtime"
	"testing"
)

var getim = func() *instanceManager {
	// get storage class
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = "C:\\Users\\wangx\\Documents\\config.yaml"
	}
	if runtime.GOOS == "linux" {
		filepath = "/root/config.yaml"
	}
	if runtime.GOOS == "darwin" {
		filepath = "./config.yaml"
	}
	config, err := ReadConfigFromFile(filepath)
	if err != nil {
		return nil
	}
	im, err := NewInstanceManagerWithConfig(config)
	if err != nil {
		return nil
	}
	return im
}

func TestFindInstance(t *testing.T) {
	// testcase
	testcase := []struct {
		id    string
		exist bool
	}{
		{"i-hgz8mri2", true},
		{"i-hgz8mri3", false},
	}

	im := getim()
	// test findVolume
	for _, v := range testcase {
		flag, err := im.FindInstance(v.id)
		if err != nil {
			t.Error("find instance error: ", err.Error())
		}
		if (flag != nil) == v.exist {
			t.Logf("instance id %s, expect %t, actually %t", v.id, v.exist, flag != nil)
		} else {
			t.Errorf("instance id %s, expect %t, actually %t", v.id, v.exist, flag != nil)
		}
	}
}
