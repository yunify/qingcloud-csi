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
