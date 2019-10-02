package gosrt

import (
	"testing"

	"github.com/pkg/errors"
)

func TestReadFile(t *testing.T) {
	var data = []struct {
		in            string
		subtitleCount int
		err           error
	}{
		{
			"./testdata/one_subtitle.srt",
			1,
			nil,
		},
		{
			"./testdata/three_subtitle.srt",
			3,
			nil,
		},
		{
			"./testdata/empty.srt",
			0,
			nil,
		},
		{
			"./testdata/invalid_subtitle.srt",
			0,
			errors.New("invalid time format: 00:00:03,200 -> 00:00:05,500"),
		},
	}

	for _, d := range data {
		t.Run(d.in, func(t *testing.T) {
			// given & when
			actual, err := ReadFile(d.in)
			// then
			if err != nil {
				if d.err.Error() != err.Error() {
					t.Errorf("In:%s, Expected:%v, Actual:%v", d.in, d.err, err)
				}
			} else {
				if d.subtitleCount != len(actual) {
					t.Errorf("In:%s, Expected:%d, Actual:%d", d.in, d.subtitleCount, len(actual))
				}
			}
		})
	}
}
