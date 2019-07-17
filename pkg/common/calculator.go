package common

import "github.com/container-storage-interface/spec/lib/go/csi"

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
func IsValidCapacityBytes(cur int64, capRanges csi.CapacityRange) bool {
	if cur < capRanges.GetRequiredBytes() || cur > capRanges.GetLimitBytes() {
		return false
	}
	return true
}
