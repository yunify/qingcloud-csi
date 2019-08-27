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
func GibToByte(nGib int) int64 {
	return int64(nGib) * Gib
}

// ByteCeilToGib
// Convert Byte to Gib
func ByteCeilToGib(nByte int64) int {
	if nByte <= 0 {
		return 0
	}
	res := nByte / Gib
	if res*Gib < nByte {
		res += 1
	}
	return int(res)
}

// Valid capacity bytes in capacity range
func IsValidCapacityBytes(cur int64, capRange *csi.CapacityRange) bool {
	if capRange == nil {
		return true
	}
	if capRange.GetRequiredBytes() > 0 && cur < capRange.GetRequiredBytes() {
		return false
	}
	if capRange.GetLimitBytes() > 0 && cur > capRange.GetLimitBytes() {
		return false
	}
	return true
}

// GetRequestSizeBytes get minimal required bytes and not exceed limit bytes.
func GetRequestSizeBytes(capRange *csi.CapacityRange) (int64, error) {
	if capRange == nil {
		return 0, nil
	}

	requiredBytes := capRange.GetRequiredBytes()
	limitBytes := capRange.GetLimitBytes()
	if requiredBytes < 0 || limitBytes < 0 {
		return -1, fmt.Errorf("capacity range [%d,%d] should not less than zero", requiredBytes, limitBytes)
	}

	if limitBytes == 0 {
		limitBytes = math.MaxInt64
	}

	if requiredBytes > limitBytes {
		return -1, fmt.Errorf("volume required bytes %d greater than limit bytes %d", requiredBytes, limitBytes)
	}
	return requiredBytes, nil
}
