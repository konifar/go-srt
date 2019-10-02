package gosrt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	timeSeparator = " --> "
)

// SubtitleScanner is the wrapper which contains next subtitle and error
type SubtitleScanner struct {
	scanner *bufio.Scanner
	nextSub Subtitle
	err     error
}

// NewScanner creates new SubtitleScanner
func NewScanner(r io.Reader) *SubtitleScanner {
	s := bufio.NewScanner(r)
	s.Split(scanDoubleNewline)
	return &SubtitleScanner{s, Subtitle{}, nil}
}

// Scan and parse subtitle
func (s *SubtitleScanner) Scan() (wasRead bool) {
	if s.scanner.Scan() {
		subtitle, err := parseSubtitle(s.scanner.Text())
		if err != nil {
			s.err = err
			return false
		}

		s.nextSub = *subtitle
		return true
	}

	return false
}

// Err gets error in scanner
func (s *SubtitleScanner) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.scanner.Err()
}

// Subtitle gets subtitle in scanner
func (s *SubtitleScanner) Subtitle() Subtitle {
	return s.nextSub
}

/*
 * 1
 * 00:00:00,000 --> 00:00:00,000
 * Don-don donuts! Let's go nuts!
**/
func parseSubtitle(text string) (*Subtitle, error) {
	var subtitle = &Subtitle{}

	elements := strings.Split(text, "\n")

	for i := 0; i < len(elements); i++ {
		text := strings.TrimRight(elements[i], "\r")
		switch i {
		case 0: // Number
			nextNum, err := strconv.Atoi(text)
			if err != nil {
				return nil, err
			}
			subtitle.Number = nextNum
		case 1: // Time
			times := strings.Split(text, timeSeparator)
			if len(times) < 2 {
				return nil, fmt.Errorf("invalid time format: %s", text)
			}

			startTime, err := parseTime(times[0])
			if err != nil {
				return nil, err
			}
			endTime, err := parseTime(times[1])
			if err != nil {
				return nil, err
			}
			subtitle.Start = startTime
			subtitle.End = endTime
		default: // Subtitle text
			if len(subtitle.Text) > 0 {
				subtitle.Text += "\n"
			}
			subtitle.Text += text
		}
	}

	return subtitle, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// Custom split func to for a double newline (one empty line)
func scanDoubleNewline(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\n', '\n'}); i >= 0 {
		return i + 2, dropCR(data[0:i]), nil
	} else if i := bytes.Index(data, []byte{'\n', '\r', '\n'}); i >= 0 {
		return i + 3, dropCR(data[0:i]), nil
	}
	if atEOF {
		return len(data), dropCR(data), nil
	}
	return 0, nil, nil
}

// Parse a formatted time like hours:minutes:seconds,milliseconds
// Ex) 00:00:00,000
func parseTime(input string) (time.Duration, error) {
	regex := regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2}),(\d{3})`)
	matches := regex.FindStringSubmatch(input)

	if len(matches) < 4 {
		return time.Duration(0), fmt.Errorf("invalid time format:%s", input)
	}

	hour, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Duration(0), err
	}
	minute, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Duration(0), err
	}
	second, err := strconv.Atoi(matches[3])
	if err != nil {
		return time.Duration(0), err
	}
	millisecond, err := strconv.Atoi(matches[4])
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(second)*time.Second + time.Duration(millisecond)*time.Millisecond), nil
}
