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

const (
	mpegAudioPacketSize = 1024
)

const (
	mpegAudioFrameSync              = 0xFFE
	mpegAudioFrameHeaderSize        = 4
	mpegAudioFrameSyncBitSize       = 11
	mpegAudioFrameVersionBitSize    = 2
	mpegAudioFrameLayerBitSize      = 2
	mpegAudioFrameProtectionBitSize = 2
)

const (
	mpegAudioVersion2_5Bits             = 0
	mpegAudioVersionReservedBits        = 1
	mpegAudioVersion2Bits               = 2
	mpegAudioVersion1Bits               = 3
	mpegAudioLayerReservedBits          = 0
	mpegAudioLayer1Bits                 = 1
	mpegAudioLayer2Bits                 = 2
	mpegAudioLayer3Bits                 = 3
	mpegAudioProtectionBitsProtected    = 0
	mpegAudioProtectionBitsNotProtected = 1
)

type ID3Tag struct {
	Header         ID3TagHeader
	TextInfoFrames []ID3TextInfoFrame
}

type ID3TagHeader struct {
	Version string
	TagSize int
}

type ID3FrameHeader struct {
	ID         string
	Size       int
	StatusFlag ID3FrameStatusFlag
	FormatFlag ID3FrameFormatFlag
}

type ID3FrameStatusFlag struct {
	TagAlterPreserved  bool
	FileAlterPreserved bool
	ReadOnly           bool
}

type ID3FrameFormatFlag struct {
	Grouped             bool
	Compressed          bool
	Encrypted           bool
	Unsynchronized      bool
	DataLengthIndicated bool
}

type ID3TextInfoFrame struct {
	Header      ID3FrameHeader
	Encoding    byte
	Description string
	Value       string
}

type MPEGAudioFrameHeader struct {
	MPEGAudioVersion string
	Layer            int
	Protected        bool
}

type MPEGAudioFrame struct {
	Header MPEGAudioFrameHeader
}

type AudioContext struct {
	ID3Tag ID3Tag
}

func readID3TextInfoUTF8Value(data []byte, m1 int) (string, int) {
	var textInfoValue string
	m2 := m1

	for data[m2] != textInfoUTF8Terminated {
		m2 += 1
	}

	textInfoValue = string(data[m1:m2])
	m2 += 1

	return textInfoValue, m2
}

func readID3TextInfoUTF16BOMValue(data []byte, m1 int) (string, int) {
	var textInfoValue string
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

	textInfoValue = string(utf16.Decode(chars))
	m2 += 2

	return textInfoValue, m2
}

func readID3TextInfoValue(data []byte, m1 int, encoding byte) (string, int) {
	var textInfoValue string
	m2 := m1

	switch encoding {
	case textInfoUTF8Encoding, textInfoISO88591:
		textInfoValue, m2 = readID3TextInfoUTF8Value(data, m1)
	case textInfoUTF16BOMEncoding:
		textInfoValue, m2 = readID3TextInfoUTF16BOMValue(data, m1)
	}

	return textInfoValue, m2
}

func readID3FrameHeader(data []byte, m1 int) (ID3FrameHeader, int) {
	var header ID3FrameHeader
	m2 := m1

	header.ID = string(data[m2 : m2+id3FrameIDSize])
	m2 += id3FrameIDSize

	header.Size = int(data[m2])<<21 +
		int(data[m2+1])<<14 +
		int(data[m2+2])<<7 +
		int(data[m2+3])
	m2 += id3FrameSizeSize

	header.StatusFlag = ID3FrameStatusFlag{
		TagAlterPreserved:  (data[m2] & 0x40) != 0x00,
		FileAlterPreserved: (data[m2] & 0x20) != 0x00,
		ReadOnly:           (data[m2] & 0x10) != 0x00,
	}
	m2 += id3FrameFlagSize

	header.FormatFlag = ID3FrameFormatFlag{
		Grouped:             (data[m2] & 0x40) != 0x00,
		Compressed:          (data[m2] & 0x08) != 0x00,
		Encrypted:           (data[m2] & 0x04) != 0x00,
		Unsynchronized:      (data[m2] & 0x02) != 0x00,
		DataLengthIndicated: (data[m2] & 0x01) != 0x00,
	}
	m2 += id3FrameFlagSize

	return header, m2
}

func readID3TextInfoFrame(data []byte, m1 int) (ID3TextInfoFrame, int) {
	var frame ID3TextInfoFrame

	header, m2 := readID3FrameHeader(data, m1)
	frame.Header = header

	frame.Encoding = data[m2]
	m2 += textInfoEncodingSize

	if frame.Header.ID == textInfoUserDefinedType {
		frame.Description, m2 = readID3TextInfoValue(data, m2, frame.Encoding)
		frame.Value, m2 = readID3TextInfoValue(data, m2, frame.Encoding)
	} else {
		frame.Value, m2 = readID3TextInfoValue(data, m2, frame.Encoding)
	}

	return frame, m2
}

func readMPEGAudioFrameHeader(data []byte, m1 int) (MPEGAudioFrameHeader, int) {
	var header MPEGAudioFrameHeader
	m2 := m1

	// Read the header by shifting offset
	var bitOffset uint32 = mpegAudioFrameHeaderSize * 8
	headerBits := binary.LittleEndian.Uint32(data[0:mpegAudioFrameHeaderSize])

	// MP3 frame sync
	bitOffset -= mpegAudioFrameSyncBitSize
	mp3SyncBits := headerBits >> bitOffset & mpegAudioFrameSync
	if mp3SyncBits != mpegAudioFrameSync {
		return MPEGAudioFrameHeader{}, -1
	}

	// MPEG audio version
	bitOffset -= mpegAudioFrameVersionBitSize
	mpegAudioVersionBits := headerBits >> bitOffset
	switch mpegAudioVersionBits {
	case mpegAudioVersion2_5Bits:
		header.MPEGAudioVersion = "2.5"
	case mpegAudioVersionReservedBits:
		return MPEGAudioFrameHeader{}, -1
	case mpegAudioVersion2Bits:
		header.MPEGAudioVersion = "2"
	case mpegAudioVersion1Bits:
		header.MPEGAudioVersion = "1"
	}

	// Layer
	bitOffset -= mpegAudioFrameLayerBitSize
	mpegAudioLayerBits := headerBits >> bitOffset
	switch mpegAudioLayerBits {
	case mpegAudioLayerReservedBits:
		return MPEGAudioFrameHeader{}, -1
	case mpegAudioLayer3Bits:
		header.Layer = 3
	case mpegAudioLayer2Bits:
		header.Layer = 2

	case mpegAudioLayer1Bits:
		header.Layer = 1
	}

	// Protection bits
	bitOffset -= mpegAudioFrameProtectionBitSize
	mpegProtectionBits := headerBits >> bitOffset
	switch mpegProtectionBits {
	case mpegAudioProtectionBitsProtected:
		header.Protected = true
	case mpegAudioProtectionBitsNotProtected:
		header.Protected = false
	}

	m2 += mpegAudioFrameHeaderSize

	return header, m2
}

func readID3TagHeader(data []byte, m1 int) (ID3TagHeader, int) {
	var header ID3TagHeader
	m2 := m1

	// ID3 magic
	if string(data[m2:m2+3]) != id3Magic {
		return ID3TagHeader{}, -1
	}

	// ID3 version
	header.Version = fmt.Sprintf("2.%d.%d", data[m2+3], data[m2+4])

	// ID3 tag size
	header.TagSize = int(data[m2+6])<<21 +
		int(data[m2+7])<<14 +
		int(data[m2+8])<<7 +
		int(data[m2+9])

	m2 += id3HeaderSize

	return header, m2
}

func readID3Tag(data []byte, m1 int) (ID3Tag, int) {
	var tag ID3Tag

	// Read tag header
	header, m2 := readID3TagHeader(data, m1)
	if m2 == -1 {
		return ID3Tag{}, -1
	}
	if id3HeaderSize+header.TagSize > len(data) {
		return ID3Tag{}, -1
	}
	tag.Header = header

	// Read tag frames
	m3 := m2
	for tag.Header.TagSize > m3-m2 {
		switch data[m3] {
		case id3FrameTextInfoType:
			frame, _m3 := readID3TextInfoFrame(data, m3)
			tag.TextInfoFrames = append(tag.TextInfoFrames, frame)
			m3 = _m3
		default:
			header, _m3 := readID3FrameHeader(data, m3)
			m3 = _m3 + header.Size
		}
	}

	return tag, m3
}

func ParseAudio(audio io.Reader, audioContext *AudioContext) error {
	return nil
}
