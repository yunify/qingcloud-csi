package block

import(
	"testing"
	"io/ioutil"
	"encoding/json"
	"os"
	"fmt"
)

var getvp = func() *volumeProvisioner{
	// get storage class
	var winfilepath = "C:\\Users\\wangx\\Documents\\config.json"
	content, err := ioutil.ReadFile(winfilepath)
	if err != nil{
		fmt.Errorf("Open file error: %s", err.Error())
		os.Exit(-1)
	}
	sc := qingStorageClass{}
	err = json.Unmarshal(content, &sc)
	if err != nil{
		fmt.Errorf("get storage class error: %s", err.Error())
		os.Exit(-1)
	}

	// get volume provisioner
	vp, err := newVolumeProvisioner(&sc)
	if err != nil{
		fmt.Errorf("new volume provisioner error: %s", err.Error())
		os.Exit(-1)
	}
	return vp
}

func TestFindVolume(t *testing.T){
	// testcase
	testcase := []struct{
		id string
		exist bool
	}{
		{"vol-fhlkhxpr", true},
		{"vol-vol-fhlkhxpw",false},
	}

	vp := getvp()
	// test findVolume
	for _, v:= range testcase{
		flag, err := vp.findVolume(v.id)
		if err != nil{
			t.Error("find volume error: ", err.Error())
		}
		if (flag != nil) == v.exist{
			t.Logf("volume id %s, expect %t, actually %t", v.id, v.exist, flag != nil)
		}else{
			t.Errorf("volume id %s, expect %t, actually %t", v.id, v.exist, flag != nil)
		}
	}
}