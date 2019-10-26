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

import "strings"

const (
	SnapshotClassTagsName = "tags"
)

type QingSnapshotClass struct {
	tags []string
}

// NewDefaultQingSnapshotClass create default QingSnapshotClass object
func NewDefaultQingSnapshotClass() *QingSnapshotClass {
	return &QingSnapshotClass{}
}

func NewQingSnapshotClassFromMap(opt map[string]string) (*QingSnapshotClass, error) {
	var tags []string
	for k, v := range opt {
		switch strings.ToLower(k) {
		case SnapshotClassTagsName:
			if len(v) > 0 {
				tags = strings.Split(strings.ReplaceAll(v, " ", ""), ",")
			}
		}
	}
	sc := NewDefaultQingSnapshotClass()
	sc.SetTags(tags)
	return sc, nil
}

func (sc QingSnapshotClass) GetTags() []string {
	return sc.tags
}

func (sc *QingSnapshotClass) SetTags(tags []string) {
	sc.tags = tags
}
