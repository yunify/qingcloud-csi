package volume

import "testing"

func TestFormatVolumeSize(t *testing.T) {
	testcase := []struct {
		name    string
		inType  int
		inSize  int
		outSize int
	}{
		{
			name:    "normal size",
			inType:  0,
			inSize:  20,
			outSize: 20,
		},
		{
			name:    "format size 1",
			inType:  200,
			inSize:  123,
			outSize: 130,
		},
		{
			name:    "format size 2",
			inType:  2,
			inSize:  123,
			outSize: 150,
		},
		{
			name:    "format size 3",
			inType:  5,
			inSize:  123,
			outSize: 200,
		},
		{
			name:    "less than min size",
			inType:  5,
			inSize:  20,
			outSize: VolumeTypeToMinSize[5],
		},
		{
			name:    "more than max size",
			inType:  2,
			inSize:  9999,
			outSize: VolumeTypeToMaxSize[2],
		},
		{
			name:    "type not found",
			inType:  1,
			inSize:  30,
			outSize: -1,
		},
	}
	for _, o := range testcase {
		resSize := FormatVolumeSize(o.inType, o.inSize)
		if resSize != o.outSize {
			t.Errorf("name %s: expect %d, but actually %d", o.name, o.outSize, resSize)
		}

	}
}
