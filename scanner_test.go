package go_srt

import (
	"github.com/Kyash/platform-api/helpers/errors"
	"testing"
)

func TestParseTime(t *testing.T) {
	var data = []struct {
		in       string
		expected *string
		err      error
	}{
		{
			"01:23:45,678",
			String("1h23m45.678s"),
			nil,
		},
		{
			"1:23:45,678",
			nil,
			errors.New("invalid time format:1:23:45,678"),
		},
		{
			"01:x3:45,678",
			nil,
			errors.New("invalid time format:01:x3:45,678"),
		},
		{
			"01:23:45,67",
			nil,
			errors.New("invalid time format:01:23:45,67"),
		},
	}

	for _, d := range data {
		t.Run(d.in, func(t *testing.T) {
			// given & when
			actual, err := parseTime(d.in)
			// then
			if err != nil {
				if d.err.Error() != err.Error() {
					t.Errorf("In:%s, Expected:%v, Actual:%v", d.in, d.err, err)
				}
			} else {
				if *d.expected != actual.String() {
					t.Errorf("In:%s, Expected:%s, Actual:%s", d.in, d.expected, actual.String())
				}
			}
		})
	}
}
