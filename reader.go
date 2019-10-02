package gosrt

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Read from io.Reader
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

// Read from filename
func ReadFile(fileName string) ([]Subtitle, error) {
	var f *os.File
	f, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		err = errors.Wrapf(err, "failed to open file:%s", fileName)
		return []Subtitle{}, err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	return ReadSubtitles(f)
}
