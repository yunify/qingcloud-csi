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
	"github.com/container-storage-interface/spec/lib/go/csi"
	"testing"
)

func TestGibToByte(t *testing.T) {
	tests := []struct {
		name  string
		gib   int
		bytes int64
	}{
		{
			name:  "normal",
			gib:   23,
			bytes: 23 * Gib,
		},
		{
			name:  "large number",
			gib:   65536000,
			bytes: 65536000 * Gib,
		},
		{
			name:  "zero Gib",
			gib:   0,
			bytes: 0,
		},
		{
			name:  "minus Gib",
			gib:   -24,
			bytes: -24 * Gib,
		},
	}
	for _, test := range tests {
		res := GibToByte(test.gib)
		if test.bytes != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.bytes, res)
		}
	}
}

func TestByteCeilToGib(t *testing.T) {
	tests := []struct {
		name  string
		nByte int64
		nGib  int
	}{
		{
			name:  "normal",
			nByte: 3 * Gib,
			nGib:  3,
		},
		{
			name:  "ceil to gib",
			nByte: 3*Gib + 3,
			nGib:  4,
		},
		{
			name:  "zero bytes",
			nByte: 0,
			nGib:  0,
		},
		{
			name:  "minus value",
			nByte: -1,
			nGib:  0,
		},
	}
	for _, test := range tests {
		res := ByteCeilToGib(test.nByte)
		if test.nGib != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.nGib, res)
		}
	}
}

func TestIsValidCapacityBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		capRange *csi.CapacityRange
		isValid  bool
	}{
		{
			name:  "normal",
			bytes: 10 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 10 * Gib,
				LimitBytes:    10 * Gib,
			},
			isValid: true,
		},
		{
			name:  "invalid range",
			bytes: 10 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 11 * Gib,
				LimitBytes:    10 * Gib,
			},
			isValid: false,
		},
		{
			name:     "empty range",
			bytes:    10 * Gib,
			capRange: &csi.CapacityRange{},
			isValid:  true,
		},
		{
			name:     "nil range",
			bytes:    10 * Gib,
			capRange: nil,
			isValid:  true,
		},
		{
			name:  "without floor",
			bytes: 10 * Gib,
			capRange: &csi.CapacityRange{
				LimitBytes: 10*Gib + 1,
			},
			isValid: true,
		},
		{
			name:  "invalid floor",
			bytes: 11 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 11*Gib + 1,
			},
			isValid: false,
		},
		{
			name:  "without ceil",
			bytes: 14 * Gib,
			capRange: &csi.CapacityRange{
				RequiredBytes: 14 * Gib,
			},
			isValid: true,
		},
		{
			name:  "invalid ceil",
			bytes: 14 * Gib,
			capRange: &csi.CapacityRange{
				LimitBytes: 14*Gib - 1,
			},
			isValid: false,
		},
	}
	for _, test := range tests {
		res := IsValidCapacityBytes(test.bytes, test.capRange)
		if test.isValid != res {
			t.Errorf("name %s: expect %t, but actually %t", test.name, test.isValid, res)
		}
	}
}

func TestGetRequestSizeBytes(t *testing.T) {
	tests := []struct {
		name     string
		capRange *csi.CapacityRange
		nBytes   int64
	}{
		{
			name: "normal",
			capRange: &csi.CapacityRange{
				RequiredBytes: 10 * Gib,
				LimitBytes:    10 * Gib,
			},
			nBytes: 10 * Gib,
		},
		{
			name:     "empty range",
			capRange: &csi.CapacityRange{},
			nBytes:   0,
		},
		{
			name:     "nil range",
			capRange: nil,
			nBytes:   0,
		},
		{
			name: "normal range 2",
			capRange: &csi.CapacityRange{
				RequiredBytes: 23 * Gib,
				LimitBytes:    25 * Gib,
			},
			nBytes: 23 * Gib,
		},
		{
			name: "invalid range",
			capRange: &csi.CapacityRange{
				RequiredBytes: 23 * Gib,
				LimitBytes:    21 * Gib,
			},
			nBytes: -1,
		},
		{
			name: "less than zero",
			capRange: &csi.CapacityRange{
				RequiredBytes: -23 * Gib,
				LimitBytes:    -21 * Gib,
			},
			nBytes: -1,
		},
	}
	for _, test := range tests {
		res, _ := GetRequestSizeBytes(test.capRange)
		if test.nBytes != res {
			t.Errorf("name %s: expect %d, but actually %d", test.name, test.nBytes, res)
		}
	}
}
