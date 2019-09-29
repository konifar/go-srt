package go_srt

import (
	"bufio"
	"bytes"
	"errors"
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

// SubtitleScanner contains the next subtitle
type SubtitleScanner struct {
	scanner *bufio.Scanner
	nextSub Subtitle
}

// NewScanner creates a new SubtitleScanner from the given io.Reader.
func NewScanner(r io.Reader) *SubtitleScanner {
	s := bufio.NewScanner(r)
	s.Split(scanDoubleNewline)
	return &SubtitleScanner{s, Subtitle{}}
}

// Parse a formatted time like hours:minutes:seconds,milliseconds
// Ex) 00:00:00,000
func parseTime(input string) (time.Duration, error) {
	regex := regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2}),(\d{3})`)
	matches := regex.FindStringSubmatch(input)

	if len(matches) < 4 {
		return time.Duration(0), errors.New("invalid time format")
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

// Scan advances the SubtitleScanner-state, reading a new Subtitle
/*
 * 1
 * 00:00:00,000 --> 00:00:00,000
 * Don-don donuts! Let's go nuts!
**/
func (s *SubtitleScanner) Scan() error {
	if s.scanner.Scan() {
		var (
			nextnum      int
			start        time.Duration
			end          time.Duration
			subtitletext string
		)

		str := strings.Split(s.scanner.Text(), "\n")

		for i := 0; i < len(str); i++ {
			text := strings.TrimRight(str[i], "\r")
			switch i {
			case 0: // Number
				num, err := strconv.Atoi(text)
				if err != nil {
					return err
				}
				nextnum = num
			case 1: // Time
				times := strings.Split(text, timeSeparator)
				if len(times) < 2 {
					return fmt.Errorf("invalid time format: %s", text)
				}

				startTime, err := parseTime(times[0])
				if err != nil {
					return err
				}
				endTime, err := parseTime(times[1])
				if err != nil {
					return err
				}
				start = startTime
				end = endTime
			default: // Subtitle text
				if len(subtitletext) > 0 {
					subtitletext += "\n"
				}
				subtitletext += text
			}
		}

		s.nextSub = Subtitle{nextnum, start, end, subtitletext}
	}

	return nil
}

func (s *SubtitleScanner) Subtitle() Subtitle {
	return s.nextSub
}
