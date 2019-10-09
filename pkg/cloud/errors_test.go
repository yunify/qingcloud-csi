package cloud

import (
	"errors"
	"testing"
)

func TestIsLeaseInfoNotReady(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		isValid bool
	}{
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-umyf1cy2] lease info not ready yet, please try later)"),
			isValid: true,
		},
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-ya2k1hpj] lease info not ready yet, please try later)"),
			isValid: true,
		},
		{
			name: "not available",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"snapshot [ss-sdm7psjv] is not available, can not create volume from it)"),
			isValid: false,
		},
	}
	for _, test := range tests {
		res := IsLeaseInfoNotReady(test.err)
		if test.isValid != res {
			t.Errorf("IsLeaseInfoNotReady(\"%v\") = %t, but want %t", test.err, res, test.isValid)
		}
	}
}

func TestIsSnapshotNotAvailable(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		isValid bool
	}{
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-umyf1cy2] lease info not ready yet, please try later)"),
			isValid: false,
		},
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-ya2k1hpj] lease info not ready yet, please try later)"),
			isValid: false,
		},
		{
			name: "not available",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"snapshot [ss-sdm7psjv] is not available, " +
				"can not create volume from it)"),
			isValid: true,
		},
	}
	for _, test := range tests {
		res := IsSnapshotNotAvailable(test.err)
		if test.isValid != res {
			t.Errorf("IsNotAvailable(%v) = %t, but want %t", test.err, res, test.isValid)
		}
	}
}

func TestIsTryLater(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		isValid bool
	}{
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-umyf1cy2] lease info not ready yet, please try later)"),
			isValid: true,
		},
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-ya2k1hpj] lease info not ready yet, please try later)"),
			isValid: true,
		},
		{
			name: "snapshot is creating, try later",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [ss-u90rxvl1] is [creating], please try later"),
			isValid: true,
		},
		{
			name: "not available",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"snapshot [ss-sdm7psjv] is not available, " +
				"can not create volume from it)"),
			isValid: false,
		},
	}
	for _, test := range tests {
		res := IsTryLater(test.err)
		if test.isValid != res {
			t.Errorf("IsTryLater(%v) = %t, but want %t", test.err, res, test.isValid)
		}
	}
}

func TestIsCannotFindDevicePath(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		isValid bool
	}{
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-umyf1cy2] lease info not ready yet, please try later)"),
			isValid: false,
		},
		{
			name: "volume lease info not ready",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [vol-ya2k1hpj] lease info not ready yet, please try later)"),
			isValid: false,
		},
		{
			name: "snapshot is creating, try later",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"resource [ss-u90rxvl1] is [creating], please try later"),
			isValid: false,
		},
		{
			name: "not available",
			err: errors.New("QingCloud Error: Code (1400), Message (PermissionDenied, " +
				"snapshot [ss-sdm7psjv] is not available, " +
				"can not create volume from it)"),
			isValid: false,
		},
		{
			name:    "cannot find device path",
			err:     NewCannotFindDevicePathError("vol-1", "ins-2", "sh1a"),
			isValid: true,
		},
	}
	for _, test := range tests {
		res := IsCannotFindDevicePath(test.err)
		if test.isValid != res {
			t.Errorf("IsCannotFindDevicePath(%v) = %t, but want %t", test.err, res, test.isValid)
		}
	}
}
