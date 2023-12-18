package gosrt

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// ReadSubtitles reads from io.Reader
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

// ReadFile reads from filename
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

func WriteFile(fileName string, subtitles []Subtitle) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, subtitle := range subtitles {
		_, err := fmt.Fprintf(file, "%d\n%s --> %s\n%s\n\n", subtitle.Number, formatTime(subtitle.Start), formatTime(subtitle.End), subtitle.Text)
		if err != nil {
			return err
		}
	}

	return nil
}

func formatTime(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)
}
