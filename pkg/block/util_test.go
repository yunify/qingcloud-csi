package block

import (
	"testing"
)

func TestReadServerConfigFromFile(t *testing.T){
	filePath := "config.yaml"
	config, err := ReadConfigFromFile(filePath)
	if err != nil{
		t.Error(err.Error())
	}
	t.Logf("%s", config.AccessKeyID)
}
