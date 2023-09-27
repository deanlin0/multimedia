package hls

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	id3Magic             = "ID3"
	id3HeaderSize        = 10
	id3FrameIDSize       = 4
	id3FrameSizeSize     = 4
	id3FrameFlagSize     = 1
	id3FrameTextInfoType = 'T'
)

const (
	textInfoEncodingSize    = 1
	textInfoUserDefinedType = "TXXX"
	textInfoTerminated      = '\x00'
	textInfoUTF8Encoding    = '\x03'
)

type ID3Header struct {
	Version        string
	TagSize        int
	TextInfoFrames []TextInfoFrame
}

type FrameStatusFlag struct {
	TagAlterPreserved  bool
	FileAlterPreserved bool
	ReadOnly           bool
}

type FrameFormatFlag struct {
	Grouped             bool
	Compressed          bool
	Encrypted           bool
	Unsynchronized      bool
	DataLengthIndicated bool
}

type TextInfoFrame struct {
	ID          string
	Size        int
	Encoding    byte
	Description string
	Value       string
	StatusFlag  FrameStatusFlag
	FormatFlag  FrameFormatFlag
}

type AudioContext struct {
	ID3Header ID3Header
}

func parseTextInfoValue(data []byte, m1 int, encoding byte) (string, int) {
	m2 := m1

	for data[m2] != textInfoTerminated {
		m2 += 1
	}

	var s string
	switch encoding {
	case textInfoUTF8Encoding:
		s = string(data[m1:m2])
	}
	m2 += 1

	return s, m2
}

func parseTextInfoFrame(data []byte, m1 int) (TextInfoFrame, int) {
	var frame TextInfoFrame

	m2 := m1
	frame.ID = string(data[m2 : m2+id3FrameIDSize])
	m2 += id3FrameIDSize
	frame.Size = int(binary.BigEndian.Uint32(data[m2 : m2+id3FrameSizeSize]))
	m2 += id3FrameSizeSize

	frame.StatusFlag = FrameStatusFlag{
		TagAlterPreserved:  (data[m2] & 0x40) != 0x00,
		FileAlterPreserved: (data[m2] & 0x20) != 0x00,
		ReadOnly:           (data[m2] & 0x10) != 0x00,
	}
	m2 += id3FrameFlagSize
	frame.FormatFlag = FrameFormatFlag{
		Grouped:             (data[m2] & 0x40) != 0x00,
		Compressed:          (data[m2] & 0x08) != 0x00,
		Encrypted:           (data[m2] & 0x04) != 0x00,
		Unsynchronized:      (data[m2] & 0x02) != 0x00,
		DataLengthIndicated: (data[m2] & 0x01) != 0x00,
	}
	m2 += id3FrameFlagSize

	frame.Encoding = data[m2]
	m2 += textInfoEncodingSize

	if frame.ID == textInfoUserDefinedType {
		frame.Description, m2 = parseTextInfoValue(data, m2, frame.Encoding)
		frame.Value, m2 = parseTextInfoValue(data, m2, frame.Encoding)
	} else {
		frame.Value, m2 = parseTextInfoValue(data, m2, frame.Encoding)
	}

	return frame, m2
}

func parseID3Header(audio io.Reader, audioContext *AudioContext) error {
	buf := make([]byte, id3HeaderSize)
	n, err := audio.Read(buf)
	if err != nil {
		return err
	}
	if n != id3HeaderSize {
		return nil
	}

	// ID3 Magic
	if string(buf[0:3]) != id3Magic {
		return nil
	}

	// ID3 version
	version := fmt.Sprintf("2.%d.%d", buf[3], buf[4])
	audioContext.ID3Header = ID3Header{
		Version: version,
	}

	// Tag
	audioContext.ID3Header.TagSize = int(binary.BigEndian.Uint32(buf[6:10]))

	buf = make([]byte, audioContext.ID3Header.TagSize)
	n, err = audio.Read(buf)
	if err != nil {
		return err
	}
	if n != audioContext.ID3Header.TagSize {
		return nil
	}

	// Tag Frames
	m1 := 0
	for m1 < audioContext.ID3Header.TagSize-id3HeaderSize {
		switch buf[m1] {
		case id3FrameTextInfoType:
			frame, m2 := parseTextInfoFrame(buf, m1)
			audioContext.ID3Header.TextInfoFrames = append(audioContext.ID3Header.TextInfoFrames, frame)
			m1 = m2
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
