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

package common

import (
	"testing"
)

func TestEntryFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
	}{
		{
			name:     "normal",
			funcName: "CreateVolume",
		},
		{
			name:     "without function name",
			funcName: "",
		},
	}
	for _, v := range tests {
		info, hash := EntryFunction(v.funcName)
		t.Logf("name %s: info %s, hash %s", v.name, info, hash)
	}
}

func TestExitFunction(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		hash     string
	}{
		{
			name:     "normal",
			funcName: "CreateVolume",
			hash:     "a7b7f7f2",
		},
		{
			name:     "without function name",
			funcName: "",
			hash:     "a7b7f7f2",
		},
	}
	for _, v := range tests {
		info := ExitFunction(v.funcName, v.hash)
		t.Logf("name %s: info %s", v.name, info)
	}
}

func TestGenerateHashInEightBytes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		hash  string
	}{
		{
			name:  "normal",
			input: "snapshot",
			hash:  "2aa38b8d",
		},
		{
			name:  "empty input",
			input: "",
			hash:  "811c9dc5",
		},
	}
	for _, v := range tests {
		res := GenerateHashInEightBytes(v.input)
		if v.hash != res {
			t.Errorf("name %s: expect %s but actually %s", v.name, v.hash, res)
		}
	}
}

func TestRetryLimiter(t *testing.T) {
	maxRetry := 5
	r := NewRetryLimiter(maxRetry)
	r.Add("2")
	r.Add("2")
	if r.Try("1") != true {
		t.Errorf("expect true but actually false")
	}
	r.Add("2")
	r.Add("3")
	r.Add("2")
	r.Add("2")
	if r.Try("2") != true {
		t.Errorf("expect true but actually false")
	}
	r.Add("2")
	if r.Try("2") != false {
		t.Errorf("expect false but actually true")
	}
}
