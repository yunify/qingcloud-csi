package block

type blockVolume struct {
	VolName string
	VolID   string
	VolSize int
	Zone    string
	Sc      qingStorageClass
}
