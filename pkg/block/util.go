package block

var VOLUME_TYPE_MAP = map[string]map[string]int{
	"sh1a": {"hp": 0, "hc": 1, "hpp": 3},
	"sh2a": {"hp": 0, "hc": 1, "hpp": 3},
}

func FormatVolumeSize(size int)int{
	if size <= 0{
		return 0
	}
	if size % 10 != 0 {
		size = (size / 10 + 1) * 10
	}
	switch{
	case size <=10:
		return 10
	case size >=500:
		return 500
	default:
		return size
	}
}
