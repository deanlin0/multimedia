package hls

import (
	"fmt"
	"io"
)

const (
	id3Magic      = "ID3"
	id3HeaderSize = 10
)

type ID3Header struct {
	Version string
}

type AudioContext struct {
	ID3Header ID3Header
}

func parseID3Header(audio io.Reader, audioContext *AudioContext) error {
	buf := make([]byte, id3HeaderSize)
	n, err := audio.Read(buf)
	if err != nil {
		return err
	}

	if n == id3HeaderSize && string(buf[0:3]) == id3Magic {
		version := fmt.Sprintf("2.%d.%d", buf[3], buf[4])
		audioContext.ID3Header = ID3Header{
			Version: version,
		}
	}

	return nil
}

func ParseAudio(audio io.Reader) (*AudioContext, error) {
	var audioContext AudioContext

	if err := parseID3Header(audio, &audioContext); err != nil {
		return nil, err
	}

	return &audioContext, nil
}
