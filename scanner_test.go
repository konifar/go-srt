package gosrt

import (
	"github.com/pkg/errors"
	"testing"
	"time"
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

func TestParseSubtitle(t *testing.T) {
	var data = []struct {
		in       string
		expected *Subtitle
		err      error
	}{
		{
			`1
00:00:01,000 --> 00:00:05,000
Don-don donuts! Let's go nuts!`,
			&Subtitle{Number: 1, Start: time.Duration(1000000000), End: time.Duration(5000000000), Text: "Don-don donuts! Let's go nuts!"},
			nil,
		},
		{
			`2
00:00:01,000 --> 00:00:05,000
Don-don donuts!
Let's go nuts!`,
			&Subtitle{Number: 2, Start: time.Duration(1000000000), End: time.Duration(5000000000), Text: "Don-don donuts!\nLet's go nuts!"},
			nil,
		},
		{
			`2
00:00:01,000 -> 00:00:05,000
Don-don donuts!
Let's go nuts!`,
			nil,
			errors.New("invalid time format: 00:00:01,000 -> 00:00:05,000"),
		},
	}

	for _, d := range data {
		t.Run(d.in, func(t *testing.T) {
			// given & when
			actual, err := parseSubtitle(d.in)
			// then
			if err != nil {
				if d.err.Error() != err.Error() {
					t.Errorf("In:%s, Expected:%v, Actual:%v", d.in, d.err, err)
				}
			} else {
				if d.expected.String() != actual.String() {
					t.Errorf("In:%s, Expected:%s, Actual:%s", d.in, d.expected.String(), actual.String())
				}
			}
		})
	}
}
