package block

import(
	qcconfig "github.com/yunify/qingcloud-sdk-go/config"
	"testing"
	"github.com/golang/glog"
)

func createConfig()(config *qcconfig.Config, err error){
	config = &qcconfig.Config{}
	config.LoadConfigFromFilepath("C:\\Users\\wangx\\Documents\\config.yaml")
	return config,err
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
		flag, err := vm.IsVolumeIdExist(v.id)
		if err != nil{
			t.Errorf("test in %s: error: %v", v.id, err)
		}
		if flag != v.ret{
			t.Errorf("testcase failed in %s, expected %t, actually %t",
				v.id, v.ret, flag)
		}else{
			t.Logf("testcase success in %s, result %t",v.id, flag)
		}
	}
}

func TestVolumeCreate(t *testing.T){
	// create volume manager
	config, err := createConfig()
	if err != nil{
		glog.Error(err)
	}
	vm, err := newVolumeManager(config)
	if err != nil{
		glog.Error(err)
	}
	testcases := []struct{
		claim volumeClaim
	}{
		{volumeClaim{VolName: "pvc-hp-0001", VolSizeRequest: 12,VolType: 0 }},
	//	{volumeClaim{VolName:"pvc-hpp-0001", VolSizeRequest: 2, VolType:3}},
	}
	for i, _:=range testcases{
		t.Logf("\tcreate volume... request = %d\n", testcases[i].claim.VolSizeRequest)
		err:= vm.CreateVolume(&testcases[i].claim)
		if err != nil{
			t.Errorf("Error: %v", err.Error())
			continue
		}
		t.Logf("\tsuccess, volume id = %d, capacity = %d\n",
			testcases[i].claim.VolID, testcases[i].claim.VolSizeCapacity)
	}
}

func TestVolumeDelete(t *testing.T){
	// create volume manager
	config, err := createConfig()
	if err != nil{
		glog.Error(err)
	}
	vm, err := newVolumeManager(config)
	if err != nil{
		glog.Error(err)
	}
	testcases := []struct{
		id string
	}{
		{"vol-rslqhbsa"},
	}
	for i, _:=range testcases{
		err := vm.DeleteVolume(testcases[i].id)
		if err != nil{
			t.Error(err)
		}else{
			t.Logf("Delete volume %s success", testcases[i].id)
		}
	}
}