package go_srt

import "time"

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