package block

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

var getip = func() *instanceProvider {
	// get storage class
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = "C:\\Users\\wangx\\Documents\\config.json"
	}
	if runtime.GOOS == "linux" {
		filepath = "/root/config.json"
	}
	if runtime.GOOS == "darwin" {
		filepath = "./config.json"
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Errorf("Open file error: %s", err.Error())
		os.Exit(-1)
	}
	sc := qingStorageClass{}
	err = json.Unmarshal(content, &sc)
	if err != nil {
		fmt.Errorf("get storage class error: %s", err.Error())
		os.Exit(-1)
	}

	// get volume vendor
	ip, err := newInstanceProvider(&sc)
	if err != nil {
		fmt.Errorf("new volume provider error: %s", err.Error())
		os.Exit(-1)
	}
	return ip
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

	ip := getip()
	// test findVolume
	for _, v := range testcase {
		flag, err := ip.findInstance(v.id)
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
