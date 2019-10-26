/*
Copyright (C) 2018 Yunify, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this work except in compliance with the License.
You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"reflect"
	"testing"
)

func TestNewDefaultQingSnapshotClass(t *testing.T) {
	tests := []struct {
		name    string
		opt     map[string]string
		sc      *QingSnapshotClass
		isError bool
	}{
		{
			name: "normal",
			opt: map[string]string{
				"tags": "tag-glozcqzd, tag-y7uu1q2a",
			},
			sc: &QingSnapshotClass{
				tags: []string{"tag-glozcqzd", "tag-y7uu1q2a"},
			},
			isError: false,
		},
		{
			name: "upper case1",
			opt: map[string]string{
				"Tags": "tag-glozcqzd, tag-y7uu1q2a",
			},
			sc: &QingSnapshotClass{
				tags: []string{"tag-glozcqzd", "tag-y7uu1q2a"},
			},
			isError: false,
		},
		{
			name: "upper case2",
			opt: map[string]string{
				"TAGS": "tag-glozcqzd, tag-y7uu1q2a",
			},
			sc: &QingSnapshotClass{
				tags: []string{"tag-glozcqzd", "tag-y7uu1q2a"},
			},
			isError: false,
		},
		{
			name: "normal",
			opt: map[string]string{
				"tags": "tag-glozcqzd, tag-y7uu1q2a,tag-y7uuweea",
			},
			sc: &QingSnapshotClass{
				tags: []string{"tag-glozcqzd", "tag-y7uu1q2a", "tag-y7uuweea"},
			},
			isError: false,
		},
		{
			name: "empty",
			opt: map[string]string{
				"tags": "",
			},
			sc: &QingSnapshotClass{
				tags: nil,
			},
			isError: false,
		},
		{
			name: "unset tags",
			opt:  map[string]string{},
			sc: &QingSnapshotClass{
				tags: nil,
			},
			isError: false,
		},
	}
	for _, test := range tests {
		res, err := NewQingSnapshotClassFromMap(test.opt)
		if (err != nil) != test.isError {
			t.Errorf("name %s: expect %t but actually %t", test.name, test.isError, err != nil)
		}
		if !reflect.DeepEqual(test.sc, res) {
			t.Errorf("name %s: expect %v but actually %v", test.name, test.sc, res)
		}
	}
}
