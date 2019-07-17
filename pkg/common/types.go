package common

const Int64Max = int64(^uint64(0) >> 1)

const (
	Kib    int64 = 1024
	Mib    int64 = Kib * 1024
	Gib    int64 = Mib * 1024
	Gib100 int64 = Gib * 100
	Tib    int64 = Gib * 1024
	Tib100 int64 = Tib * 100
)

const (
	FileSystemExt3    string = "ext3"
	FileSystemExt4    string = "ext4"
	FileSystemXfs     string = "xfs"
	DefaultFileSystem string = FileSystemExt4
)
