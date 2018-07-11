package block

import (
	"reflect"
	"testing"
)

func TestNewQingStorageClassFromMap(t *testing.T) {
	testcases := []struct {
		name  string
		mp    map[string]string
		sc    qingStorageClass
		isErr bool
	}{
		{
			name: "normal",
			mp: map[string]string{
				"type":    "0",
				"maxSize": "1000",
				"minSize": "10",
			},
			sc: qingStorageClass{
				VolumeType:    0,
				VolumeMaxSize: 1000,
				VolumeMinSize: 10,
			},
			isErr: false,
		},
		{
			name: "type is string",
			mp: map[string]string{
				"type":    "k",
				"maxSize": "1000",
				"minSize": "10",
			},
			sc:    qingStorageClass{},
			isErr: true,
		},
		{
			name: "size is string",
			mp: map[string]string{
				"type":    "0",
				"maxSize": "s",
				"minSize": "10",
			},
			sc:    qingStorageClass{},
			isErr: true,
		},
		{
			name: "max size less than min size",
			mp: map[string]string{
				"type":    "0",
				"maxSize": "1000",
				"minSize": "1001",
			},
			sc:    qingStorageClass{},
			isErr: true,
		},
		{
			name: "max size equal to min size",
			mp: map[string]string{
				"type":    "0",
				"maxSize": "1000",
				"minSize": "1000",
			},
			sc: qingStorageClass{
				VolumeType:    0,
				VolumeMaxSize: 1000,
				VolumeMinSize: 1000,
			},
			isErr: false,
		},
		{
			name: "size less than zero",
			mp: map[string]string{
				"type":    "0",
				"maxSize": "1000",
				"minSize": "-2",
			},
			sc:    qingStorageClass{},
			isErr: true,
		},
	}
	for _, v := range testcases {
		res, err := NewQingStorageClassFromMap(v.mp)
		if err != nil {
			if !v.isErr {
				t.Errorf("name %s raise error: %s", v.name, err.Error())
			}
		} else if v.isErr && err == nil {
			t.Errorf("name %s: expect error occur %t, but actually %t", v.name, v.isErr, !v.isErr)
		} else if !v.isErr && !reflect.DeepEqual(*res, v.sc) {
			t.Errorf("name %s: expect [%+v], but actually [%+v]", v.name, v.sc, res)
		}
	}
}

func TestFormatVolumeSize(t *testing.T) {
	testcases := []struct {
		name   string
		sc     qingStorageClass
		size   int
		result int
	}{
		{
			name: "normal sc, normal size",
			sc: qingStorageClass{
				VolumeMinSize: 10,
				VolumeMaxSize: 500,
			},
			size:   24,
			result: 30,
		},
		{
			name: "normal sc, size less than zero",
			sc: qingStorageClass{
				VolumeMinSize: 10,
				VolumeMaxSize: 500,
			},
			size:   -1,
			result: 10,
		},
		{
			name: "normal sc, size less than min size",
			sc: qingStorageClass{
				VolumeMinSize: 10,
				VolumeMaxSize: 500,
			},
			size:   -1,
			result: 10,
		},
		{
			name: "normal sc, size equal to max size",
			sc: qingStorageClass{
				VolumeMinSize: 10,
				VolumeMaxSize: 500,
			},
			size:   500,
			result: 500,
		},
		{
			name: "normal sc, size greater than max size",
			sc: qingStorageClass{
				VolumeMinSize: 10,
				VolumeMaxSize: 500,
			},
			size:   501,
			result: 500,
		},
		{
			name: "narrow sc, ceil size 1",
			sc: qingStorageClass{
				VolumeMinSize: 74,
				VolumeMaxSize: 83,
			},
			size:   77,
			result: 80,
		},
		{
			name: "narrow sc, ceil size 2",
			sc: qingStorageClass{
				VolumeMinSize: 74,
				VolumeMaxSize: 83,
			},
			size:   82,
			result: 83,
		},
		{
			name: "narrow sc, size greater than max size",
			sc: qingStorageClass{
				VolumeMinSize: 74,
				VolumeMaxSize: 83,
			},
			size:   89,
			result: 83,
		},
		{
			name: "narrow sc, size less than max size",
			sc: qingStorageClass{
				VolumeMinSize: 74,
				VolumeMaxSize: 83,
			},
			size:   71,
			result: 74,
		},
		{
			name: "narrow sc, size less than max size",
			sc: qingStorageClass{
				VolumeMinSize: 74,
				VolumeMaxSize: 83,
			},
			size:   71,
			result: 74,
		},
		{
			name: "equal sc, size less than min size 1",
			sc: qingStorageClass{
				VolumeMinSize: 502,
				VolumeMaxSize: 502,
			},
			size:   23,
			result: 502,
		},
		{
			name: "equal sc, size less than min size 2",
			sc: qingStorageClass{
				VolumeMinSize: 502,
				VolumeMaxSize: 502,
			},
			size:   501,
			result: 502,
		},
		{
			name: "equal sc, size equal to max size",
			sc: qingStorageClass{
				VolumeMinSize: 502,
				VolumeMaxSize: 502,
			},
			size:   502,
			result: 502,
		},
		{
			name: "equal sc, size greater than max size 1",
			sc: qingStorageClass{
				VolumeMinSize: 502,
				VolumeMaxSize: 502,
			},
			size:   505,
			result: 502,
		},
		{
			name: "equal sc, size greater than max size 2",
			sc: qingStorageClass{
				VolumeMinSize: 502,
				VolumeMaxSize: 502,
			},
			size:   643,
			result: 502,
		},
	}
	for _, v := range testcases {
		res := v.sc.FormatVolumeSize(v.size)
		if res != v.result {
			t.Errorf("name %s, expect %d, but actually %d", v.name, v.result, res)
		}
	}
}
