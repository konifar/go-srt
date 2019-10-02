package go_srt

import (
	"fmt"
	"time"
)

/*
 1
 00:00:00,000 --> 00:00:00,000
 Don-don donuts! Let's go nuts!
*/
type Subtitle struct {
	Number int
	Start  time.Duration
	End    time.Duration
	Text   string
}

func (s *Subtitle) String() string {
	return fmt.Sprintf("Number:%d, Start:%v, End:%v, Text:%s", s.Number, s.Start, s.End, s.Text)
}
