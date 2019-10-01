package go_srt

import (
	"github.com/pkg/errors"
	"io"
	"os"
)

func ReadSubtitles(r io.Reader) (subtitles []Subtitle, err error) {
	scanner := NewScanner(r)
	for scanner.Scan() {
		subtitles = append(subtitles, scanner.Subtitle())
	}
	err = scanner.Err()
	if err != nil {
		subtitles = []Subtitle{}
	}
	return
}

func ReadFile(fileName string) ([]Subtitle, error) {
	var f *os.File
	f, err := os.Open(fileName)
	if err != nil {
		err = errors.Wrapf(err, "failed to open file:%s", fileName)
		return []Subtitle{}, err
	}
	defer f.Close()

	return ReadSubtitles(f)
}
