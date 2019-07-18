package common

import "k8s.io/kubernetes/pkg/util/mount"

func NewSafeMounter() *mount.SafeFormatAndMount {
	realMounter := mount.New("")
	realExec := mount.NewOsExec()
	return &mount.SafeFormatAndMount{
		Interface: realMounter,
		Exec:      realExec,
	}
}
