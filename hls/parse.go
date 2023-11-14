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
	vbrXingMagic       = "Xing"
	vbrInfoMagic       = "Info"
	vbrFlagSize        = 4
	vbrNumOfFramesSize = 4
	vbrFileSizeSize    = 4
	vbrTOCSize         = 100
	vbrQualitySize     = 4
)

const (
	mpegAudioPacketSize = 1024
)

const (
	mpegAudioFrameSync              = 0x7FF
	mpegAudioFrameHeaderSize        = 4
	mpegAudioFrameSyncBitSize       = 11
	mpegAudioFrameVersionBitSize    = 2
	mpegAudioFrameLayerBitSize      = 2
	mpegAudioFrameProtectionBitSize = 1
	mpegAudioFrameBitrateBitSize    = 4
	mpegAudioFrameSampleRateBitSize = 2
)

const (
	mpegAudioVersion2_5Bits             = 0
	mpegAudioVersionReservedBits        = 1
	mpegAudioVersion2Bits               = 2
	mpegAudioVersion1Bits               = 3
	mpegAudioLayerReservedBits          = 0
	mpegAudioLayer1Bits                 = 3
	mpegAudioLayer2Bits                 = 2
	mpegAudioLayer3Bits                 = 1
	mpegAudioProtectionBitsProtected    = 0
	mpegAudioProtectionBitsNotProtected = 1
)

var (
	mpegBitrateMap [2][3][15]int = [2][3][15]int{
		{
			{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448},
			{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384},
			{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320},
		},
		{
			{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256},
			{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
			{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
		},
	}
	mpegSampleRateMap [3][3]int = [3][3]int{
		{44100, 48000, 32000},
		{22050, 24000, 16000},
		{11025, 12000, 8000},
	}
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

type VBRHeader struct {
	ID          string
	NumOfFrames *int
	FileSize    *int
	TOC         []int
	Quality     *int
}

type MPEGAudioFrameHeader struct {
	MPEGAudioVersion string
	Layer            int
	Protected        bool
	Bitrate          int
	SampleRate       int
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
	var bitMask uint32
	headerBits := binary.BigEndian.Uint32(data[0:mpegAudioFrameHeaderSize])

	// MP3 frame sync
	bitOffset -= mpegAudioFrameSyncBitSize
	bitMask = 1<<mpegAudioFrameSyncBitSize - 1
	mpegAudioSyncBits := (headerBits >> bitOffset) & bitMask
	if mpegAudioSyncBits != mpegAudioFrameSync {
		return MPEGAudioFrameHeader{}, -1
	}

	// MPEG audio version
	bitOffset -= mpegAudioFrameVersionBitSize
	bitMask = 1<<mpegAudioFrameVersionBitSize - 1
	mpegAudioVersionBits := (headerBits >> bitOffset) & bitMask
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
	bitMask = 1<<mpegAudioFrameLayerBitSize - 1
	mpegAudioLayerBits := (headerBits >> bitOffset) & bitMask
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

	// Protection bit
	bitOffset -= mpegAudioFrameProtectionBitSize
	bitMask = 1<<mpegAudioFrameProtectionBitSize - 1
	mpegProtectionBits := (headerBits >> bitOffset) & bitMask
	switch mpegProtectionBits {
	case mpegAudioProtectionBitsProtected:
		header.Protected = true
	case mpegAudioProtectionBitsNotProtected:
		header.Protected = false
	}

	// Bitrate
	bitOffset -= mpegAudioFrameBitrateBitSize
	bitMask = 1<<mpegAudioFrameBitrateBitSize - 1
	mpegBitrateBits := (headerBits >> bitOffset) & bitMask
	switch header.MPEGAudioVersion {
	case "1":
		header.Bitrate = mpegBitrateMap[0][header.Layer-1][mpegBitrateBits]
	case "2.5", "2":
		header.Bitrate = mpegBitrateMap[1][header.Layer-1][mpegBitrateBits]
	}

	// Sample rate
	bitOffset -= mpegAudioFrameSampleRateBitSize
	bitMask = 1<<mpegAudioFrameSampleRateBitSize - 1
	mpegSampleRateBits := (headerBits >> bitOffset) & bitMask
	switch header.MPEGAudioVersion {
	case "1":
		header.SampleRate = mpegSampleRateMap[0][mpegSampleRateBits]
	case "2":
		header.SampleRate = mpegSampleRateMap[1][mpegSampleRateBits]
	case "2.5":
		header.SampleRate = mpegSampleRateMap[2][mpegSampleRateBits]
	}

	m2 += mpegAudioFrameHeaderSize

	return header, m2
}

func readVBRTOC(data []byte, m1 int) ([]int, int) {
	var toc []int
	m2 := m1

	if len(data) < m2+vbrTOCSize {
		return nil, -1
	}

	for i := m2; i < m2+vbrTOCSize; i++ {
		pos := int(data[i])
		toc = append(toc, pos)
	}
	m2 += vbrTOCSize

	return toc, m2
}

func readVBRHeader(data []byte, m1 int) (VBRHeader, int) {
	var header VBRHeader
	m2 := m1

	// VBR ID
	id := string(data[m2 : m2+4])
	if id != vbrXingMagic && id != vbrInfoMagic {
		return VBRHeader{}, -1
	}
	header.ID = id
	m2 += 4

	// VBR flags
	flagBits := binary.BigEndian.Uint32(data[m2 : m2+vbrFlagSize])
	m2 += vbrFlagSize

	// Read number of frames, file size, toc, and quality if the flag is set
	if flagBits&0x01 == 0x01 {
		numOfFrames := int(binary.BigEndian.Uint32(data[m2 : m2+vbrNumOfFramesSize]))
		header.NumOfFrames = &numOfFrames
		m2 += vbrNumOfFramesSize
	}
	if flagBits&0x02 == 0x02 {
		fileSize := int(binary.BigEndian.Uint32(data[m2 : m2+vbrFileSizeSize]))
		header.FileSize = &fileSize
		m2 += vbrFileSizeSize
	}
	if flagBits&0x04 == 0x04 {
		toc, _m2 := readVBRTOC(data, m2)
		header.TOC = toc
		m2 = _m2
	}
	if flagBits&0x08 == 0x08 {
		quality := int(binary.BigEndian.Uint32(data[m2 : m2+vbrQualitySize]))
		header.Quality = &quality
		m2 += vbrQualitySize
	}

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
