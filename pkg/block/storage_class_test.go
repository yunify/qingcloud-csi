package block

import (
	"testing"
	"io/ioutil"
	"encoding/json"
)

var winfilepath = "C:\\Users\\wangx\\Documents\\config.json"

func Test_getConfigFromQingStorageClass(t *testing.T){
	content, err := ioutil.ReadFile(winfilepath)
	if err != nil{
		t.Error("Open file error: ", err.Error())
	}
	sc := qingStorageClass{}
	err = json.Unmarshal(content, &sc)
	if err != nil{
		t.Error("get storage class error: ", err.Error())
	}
	// print storage class
	bytes, _:=json.Marshal(sc)
	t.Log("storage class:", string(bytes[:]))
	// get config
	config := getConfigFromStorageClass(&sc)
	// print config
	t.Log("config:", config)
}
