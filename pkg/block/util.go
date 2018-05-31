package block

func FormatVolumeSize(size int)int{
	if size % 10 != 0 {
		size = (size / 10 + 1) * 10
	}
	if size <= 10{
		return 10
	}else if size >=500{
		return 500
	}else{
		return size
	}
}
