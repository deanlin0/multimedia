package hls

import (
	"fmt"
)

type Segment struct {
	URI      string
	Duration float64
}

func BuildSegmentLine(segment Segment) string {
	line := fmt.Sprintf("#%s:%g,", EXTINF, segment.Duration)
	line = fmt.Sprintf("%s\n%s", line, segment.URI)

	return line
}
