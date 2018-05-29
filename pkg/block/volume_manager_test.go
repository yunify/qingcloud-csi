package block

import(
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"testing"
	"github.com/golang/glog"
)

func createConfig()(config *qcconfig.Config, err error){
	accessKeyId := ""
	accessKeySecret := ""
	config,err = qcconfig.New(accessKeyId,accessKeySecret)
	config.Zone = "sh1a"
	if err != nil{
		return nil, err
	}else{
		return config,nil
	}
}

func TestVolumeIdExist(t *testing.T){
	// create volume manager
	config, err := createConfig()
	if err != nil{
		glog.Error(err)
	}
	vm, err := newVolumeManager(config)
	if err != nil{
		glog.Error(err)
	}

	testcase := []struct{
		id string
		ret bool
	}{
		{"vol-57sm6cas", true},
		{"vol-aseereww", false},
	}
	for _, v:=range testcase{
		vm.persistentVolume.VolID = v.id
		flag, err := vm.IsVolumeIdExist()
		if err != nil{
			t.Errorf("test in %s: error: %v", v.id, err)
		}
		if flag != v.ret{
			t.Errorf("testcase failed in %s, expected %t, actually %t",
				v.id, v.ret, flag)
		}
	}


}