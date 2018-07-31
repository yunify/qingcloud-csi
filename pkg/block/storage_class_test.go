// +-------------------------------------------------------------------------
// | Copyright (C) 2018 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package block

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewQingStorageClassFromMap(t *testing.T) {
	testcases := []struct {
		name     string
		mp       map[string]string
		sc       qingStorageClass
		isError  bool
		strError string
	}{
		{
			name: "normal",
			mp: map[string]string{
				"type":     "0",
				"maxSize":  "1000",
				"minSize":  "10",
				"stepSize": "10",
				"fsType":   "ext4",
			},
			sc: qingStorageClass{
				VolumeType:     0,
				VolumeMaxSize:  1000,
				VolumeMinSize:  10,
				VolumeStepSize: 10,
				VolumeFsType:   FileSystemExt4,
			},
			isError:  false,
			strError: "",
		},
		{
			name: "default storageclass",
			mp:   map[string]string{},
			sc: qingStorageClass{
				VolumeType:     0,
				VolumeMaxSize:  500,
				VolumeMinSize:  10,
				VolumeStepSize: 10,
				VolumeFsType:   FileSystemExt4,
			},
			isError:  false,
			strError: "",
		},
		{
			name: "type is string",
			mp: map[string]string{
				"type":     "k",
				"maxSize":  "1000",
				"minSize":  "10",
				"stepSize": "10",
				"fsType":   "xfs",
			},
			isError:  true,
			strError: "strconv.Atoi: parsing",
		},
		{
			name: "size is string",
			mp: map[string]string{
				"type":     "0",
				"maxSize":  "s",
				"minSize":  "10",
				"stepSize": "10",
				"fsType":   "xfs",
			},
			sc:       qingStorageClass{},
			isError:  true,
			strError: "strconv.Atoi: parsing",
		},
		{
			name: "max size less than min size",
			mp: map[string]string{
				"type":     "0",
				"maxSize":  "1000",
				"minSize":  "1001",
				"stepSize": "10",
				"fsType":   "ext3",
			},
			sc:       qingStorageClass{},
			isError:  true,
			strError: "Volume maxSize must greater than or equal to volume minSize",
		},
		{
			name: "max size equal to min size",
			mp: map[string]string{
				"type":     "0",
				"maxSize":  "1000",
				"minSize":  "1000",
				"stepSize": "10",
				"fsType":   "ext4",
			},
			sc: qingStorageClass{
				VolumeType:     0,
				VolumeMaxSize:  1000,
				VolumeMinSize:  1000,
				VolumeStepSize: 10,
				VolumeFsType:   FileSystemExt4,
			},
			isError:  false,
			strError: "",
		},
		{
			name: "size less than zero",
			mp: map[string]string{
				"type":     "0",
				"maxSize":  "1000",
				"minSize":  "-2",
				"stepSize": "10",
				"fsType":   "ext4",
			},
			sc:       qingStorageClass{},
			isError:  true,
			strError: "MinSize must not less than zero",
		},
		{
			name: "step size equal to zero",
			mp: map[string]string{
				"type":     "0",
				"maxSize":  "1000",
				"minSize":  "200",
				"stepSize": "0",
				"fsType":   "ext4",
			},
			sc:       qingStorageClass{},
			isError:  true,
			strError: "StepSize must greate than zero",
		},
		{
			name: "input empty fsType",
			mp: map[string]string{
				"type":     "0",
				"maxSize":  "1000",
				"minSize":  "1000",
				"stepSize": "2",
				"fsType":   "",
			},
			sc: qingStorageClass{
				VolumeType:     0,
				VolumeMaxSize:  1000,
				VolumeMinSize:  1000,
				VolumeStepSize: 2,
				VolumeFsType:   "",
			},
			isError:  true,
			strError: "Does not support fsType",
		},
		{
			name: "input wrong fsType",
			mp: map[string]string{
				"type":    "0",
				"maxSize": "1000",
				"minSize": "1000",
				"fsType":  "wrong",
			},
			sc:       qingStorageClass{},
			isError:  true,
			strError: "Does not support fsType",
		},
		{
			name: "not input fsType",
			mp: map[string]string{
				"type":    "0",
				"maxSize": "1000",
				"minSize": "1000",
			},
			sc: qingStorageClass{
				VolumeType:     0,
				VolumeMaxSize:  1000,
				VolumeMinSize:  1000,
				VolumeStepSize: 10,
				VolumeFsType:   FileSystemExt4,
			},
			isError:  false,
			strError: "",
		},
	}
	for _, v := range testcases {
		res, err := NewQingStorageClassFromMap(v.mp)
		if err != nil {
			if v.isError == false {
				t.Errorf("name %s: expect %t, actually false [%s]", v.name, v.isError, err.Error())
			} else if v.isError == true && !strings.Contains(err.Error(), v.strError) {
				t.Errorf("name %s: expect [%s], actually [%s]", v.name, v.strError, err.Error())
			}
		} else if !reflect.DeepEqual(*res, v.sc) {
			t.Errorf("name %s: sc does not equal", v.name)
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
				VolumeMinSize:  10,
				VolumeMaxSize:  500,
				VolumeStepSize: 10,
			},
			size:   24,
			result: 30,
		},
		{
			name: "normal sc, size less than zero",
			sc: qingStorageClass{
				VolumeMinSize:  10,
				VolumeMaxSize:  500,
				VolumeStepSize: 10,
			},
			size:   -1,
			result: 10,
		},
		{
			name: "normal sc, size less than min size",
			sc: qingStorageClass{
				VolumeMinSize:  10,
				VolumeMaxSize:  500,
				VolumeStepSize: 10,
			},
			size:   8,
			result: 10,
		},
		{
			name: "normal sc, size equal to max size",
			sc: qingStorageClass{
				VolumeMinSize:  10,
				VolumeMaxSize:  500,
				VolumeStepSize: 10,
			},
			size:   500,
			result: 500,
		},
		{
			name: "normal sc, size greater than max size",
			sc: qingStorageClass{
				VolumeMinSize:  10,
				VolumeMaxSize:  500,
				VolumeStepSize: 10,
			},
			size:   501,
			result: 500,
		},
		{
			name: "equal sc, size less than min size 1",
			sc: qingStorageClass{
				VolumeMinSize:  502,
				VolumeMaxSize:  502,
				VolumeStepSize: 10,
			},
			size:   23,
			result: 502,
		},
		{
			name: "step size is 100",
			sc: qingStorageClass{
				VolumeMinSize:  100,
				VolumeMaxSize:  6000,
				VolumeStepSize: 100,
			},
			size:   443,
			result: 500,
		},
		{
			name: "step size is 50",
			sc: qingStorageClass{
				VolumeMinSize:  100,
				VolumeMaxSize:  6000,
				VolumeStepSize: 50,
			},
			size:   433,
			result: 450,
		},
	}
	for _, v := range testcases {
		res := v.sc.FormatVolumeSize(v.size, v.sc.VolumeStepSize)
		if res != v.result {
			t.Errorf("name %s, expect %d, but actually %d", v.name, v.result, res)
		}
	}
}
