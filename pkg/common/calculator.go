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
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"math"
)

// GibToByte
// Convert GiB to Byte
func GibToByte(num int) int64 {
	return int64(num) * Gib
}

// ByteCeilToGib
// Convert Byte to Gib
func ByteCeilToGib(num int64) int {
	if num <= 0 {
		return 0
	}
	res := num / Gib
	if res*Gib < num {
		res += 1
	}
	return int(res)
}

// Valid capacity bytes in capacity range
func IsValidCapacityBytes(cur int64, capRanges *csi.CapacityRange) bool {
	if capRanges == nil {
		return true
	}
	if capRanges.GetRequiredBytes() > 0 && cur < capRanges.GetRequiredBytes() {
		return false
	}
	if capRanges.GetLimitBytes() > 0 && cur > capRanges.GetLimitBytes() {
		return false
	}
	return true
}

func GetRequestSizeBytes(capRange *csi.CapacityRange) (int64, error) {
	if capRange == nil {
		return 0, nil
	}

	requiredBytes := capRange.GetRequiredBytes()

	limitBytes := capRange.GetLimitBytes()
	if limitBytes == 0 {
		limitBytes = math.MaxInt64
	}

	if requiredBytes > limitBytes {
		return -1, fmt.Errorf("volume required bytes %d greater than limit bytes %d", requiredBytes, limitBytes)
	}
	return requiredBytes, nil
}
