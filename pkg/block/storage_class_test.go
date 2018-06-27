package block

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
	"testing"
)

func Test_getConfigFromQingStorageClass(t *testing.T) {
	// new default storageclass
	sc := NewDefaultQingStorageClass()
	// get storageclass
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = "C:\\Users\\wangx\\Documents\\config.json"
	}
	if runtime.GOOS == "linux" {
		filepath = "/root/config.json"
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Error("Open file error: ", err.Error())
	}
	err = json.Unmarshal(content, &sc)
	if err != nil {
		t.Error("get storage class error: ", err.Error())
	}

	// print storage class
	bytes, _ := json.Marshal(sc)
	t.Log("storage class:", string(bytes[:]))
}

func TestFormatVolumeSize(t *testing.T) {
	// new default storageclass
	sc := NewDefaultQingStorageClass()
	// get storageclass
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = "C:\\Users\\wangx\\Documents\\config.json"
	}
	if runtime.GOOS == "linux" {
		filepath = "/root/config.json"
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Error("Open file error: ", err.Error())
	}
	err = json.Unmarshal(content, &sc)
	if err != nil {
		t.Error("get storage class error: ", err.Error())
	}
	// testcase
	testcase := []struct {
		size   int
		result int
	}{
		{-1, sc.VolumeMinSize},
		{0, sc.VolumeMinSize},
		{10, sc.VolumeMinSize},
		{9, sc.VolumeMinSize},
		{34, 40},
		{258, 260},
		{1091, sc.VolumeMaxSize},
	}
	for _, v := range testcase {
		actual := sc.formatVolumeSize(v.size)
		if actual == v.result {
			t.Logf("testcase success")
		} else {
			t.Errorf("testcase size: %d, expect %d, actually %d",
				v.size, v.result, actual)
		}
	}

}
