package hls

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf16"
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
	textInfoEncodingSize       = 1
	textInfoUserDefinedType    = "TXXX"
	textInfoUTF8Terminated     = 0x00
	textInfoUTF16BOMTerminated = 0x00
	textInfoUTF8Encoding       = 0x03
	textInfoUTF16BOMEncoding   = 0x01
	textInfoISO88591           = 0x00
)

const (
	utf16BOM = 0xFEFF
)

type ID3Tag struct {
	Version        string
	Size           int
	TextInfoFrames []TextInfoFrame
}

type FrameHeader struct {
	ID         string
	Size       int
	StatusFlag FrameStatusFlag
	FormatFlag FrameFormatFlag
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
	Header      FrameHeader
	Encoding    byte
	Description string
	Value       string
}

type AudioContext struct {
	ID3Tag ID3Tag
}

func readTextInfoUTF8Value(data []byte, m1 int) (string, int) {
	m2 := m1

	for data[m2] != textInfoUTF8Terminated {
		m2 += 1
	}

	s := string(data[m1:m2])
	m2 += 1

	return s, m2
}

func readTextInfoUTF16BOMValue(data []byte, m1 int) (string, int) {
	m2 := m1

	var isBigEndian bool
	if binary.BigEndian.Uint16(data[m2:m2+2]) == utf16BOM {
		isBigEndian = true
	}
	m2 += 2

	var chars []uint16
	if isBigEndian {
		char := binary.BigEndian.Uint16(data[m2 : m2+2])
		for char != textInfoUTF16BOMTerminated {
			chars = append(chars, char)
			m2 += 2
			char = binary.BigEndian.Uint16(data[m2 : m2+2])
		}
	} else {
		char := binary.LittleEndian.Uint16(data[m2 : m2+2])
		for char != textInfoUTF16BOMTerminated {
			chars = append(chars, char)
			m2 += 2
			char = binary.LittleEndian.Uint16(data[m2 : m2+2])
		}
	}

	s := string(utf16.Decode(chars))
	m2 += 2

	return s, m2
}

func readTextInfoValue(data []byte, m1 int, encoding byte) (string, int) {
	s := ""
	m2 := m1

	switch encoding {
	case textInfoUTF8Encoding, textInfoISO88591:
		s, m2 = readTextInfoUTF8Value(data, m1)
	case textInfoUTF16BOMEncoding:
		s, m2 = readTextInfoUTF16BOMValue(data, m1)
	}

	return s, m2
}

func readFrameHeader(data []byte, m1 int) (FrameHeader, int) {
	var header FrameHeader

	m2 := m1
	header.ID = string(data[m2 : m2+id3FrameIDSize])
	m2 += id3FrameIDSize
	header.Size = int(data[m2])<<21 +
		int(data[m2+1])<<14 +
		int(data[m2+2])<<7 +
		int(data[m2+3])
	m2 += id3FrameSizeSize

	header.StatusFlag = FrameStatusFlag{
		TagAlterPreserved:  (data[m2] & 0x40) != 0x00,
		FileAlterPreserved: (data[m2] & 0x20) != 0x00,
		ReadOnly:           (data[m2] & 0x10) != 0x00,
	}
	m2 += id3FrameFlagSize
	header.FormatFlag = FrameFormatFlag{
		Grouped:             (data[m2] & 0x40) != 0x00,
		Compressed:          (data[m2] & 0x08) != 0x00,
		Encrypted:           (data[m2] & 0x04) != 0x00,
		Unsynchronized:      (data[m2] & 0x02) != 0x00,
		DataLengthIndicated: (data[m2] & 0x01) != 0x00,
	}
	m2 += id3FrameFlagSize

	return header, m2
}

func readTextInfoFrame(data []byte, m1 int) (TextInfoFrame, int) {
	var frame TextInfoFrame

	header, m2 := readFrameHeader(data, m1)
	frame.Header = header

	frame.Encoding = data[m2]
	m2 += textInfoEncodingSize

	if frame.Header.ID == textInfoUserDefinedType {
		frame.Description, m2 = readTextInfoValue(data, m2, frame.Encoding)
		frame.Value, m2 = readTextInfoValue(data, m2, frame.Encoding)
	} else {
		frame.Value, m2 = readTextInfoValue(data, m2, frame.Encoding)
	}

	return frame, m2
}

func parseID3Tag(audio io.Reader, audioContext *AudioContext) error {
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
	audioContext.ID3Tag = ID3Tag{
		Version: version,
	}

	// Tag size
	audioContext.ID3Tag.Size = int(buf[6])<<21 +
		int(buf[7])<<14 +
		int(buf[8])<<7 +
		int(buf[9])

	// Read tag frames
	buf = make([]byte, audioContext.ID3Tag.Size)
	n, err = audio.Read(buf)
	if err != nil {
		return err
	}
	if n != audioContext.ID3Tag.Size {
		return nil
	}

	// Parse tag frames
	m1 := 0
	for m1 < audioContext.ID3Tag.Size-id3HeaderSize {
		switch buf[m1] {
		case id3FrameTextInfoType:
			frame, m2 := readTextInfoFrame(buf, m1)
			audioContext.ID3Tag.TextInfoFrames = append(audioContext.ID3Tag.TextInfoFrames, frame)
			m1 = m2
		default:
			header, m2 := readFrameHeader(buf, m1)
			m1 = m2 + header.Size
		}
	}

	return nil
}

func ParseAudio(audio io.Reader, audioContext *AudioContext) error {
	if err := parseID3Tag(audio, audioContext); err != nil {
		return err
	}

	return nil
}
