package hls

import (
	"strings"
)

type Playlist struct {
	lines []string
}

func (p *Playlist) AddLine(line string) *Playlist {
	p.lines = append(p.lines, line)
	return p
}

func (p *Playlist) File() []byte {
	s := strings.Join(p.lines, "\n")
	return []byte(s)
}
