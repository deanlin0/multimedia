package id3

import (
	"os"
	"reflect"
	"testing"
)

func TestParseAudio_readID3Tag(t *testing.T) {
	testCases := []struct {
		filePath string
		want     ID3Tag
	}{
		{
			filePath: "./testdata/ramp_jazz.mp3",
			want: ID3Tag{
				Header: ID3TagHeader{
					Version: "2.4.0",
					TagSize: 85,
				},
				TextInfoFrames: []ID3TextInfoFrame{
					{
						Header:      ID3FrameHeader{"TDRC", 12, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF8Encoding,
						Description: "",
						Value:       "2022-11-16",
					},
					{
						Header:      ID3FrameHeader{"TXXX", 18, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF8Encoding,
						Description: "time_reference",
						Value:       "0",
					},
					{
						Header:      ID3FrameHeader{"TSSE", 15, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF8Encoding,
						Description: "",
						Value:       "Lavf59.27.100",
					},
				},
			},
		},
		{
			filePath: "./testdata/funky_weekend_night.mp3",
			want: ID3Tag{
				Header: ID3TagHeader{
					Version: "2.3.0",
					TagSize: 640733,
				},
				TextInfoFrames: []ID3TextInfoFrame{
					{
						Header:      ID3FrameHeader{"TALB", 15, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF16BOMEncoding,
						Description: "",
						Value:       "DOVA用",
					},
					{
						Header:      ID3FrameHeader{"TPE1", 15, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF16BOMEncoding,
						Description: "",
						Value:       "蒲鉾さちこ",
					},
					{
						Header:      ID3FrameHeader{"TPE2", 15, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF16BOMEncoding,
						Description: "",
						Value:       "蒲鉾さちこ",
					},
					{
						Header:      ID3FrameHeader{"TCOM", 51, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF16BOMEncoding,
						Description: "",
						Value:       "蒲鉾さちこ(Kamaboko Sachiko)",
					},
					{
						Header:      ID3FrameHeader{"TCON", 23, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF16BOMEncoding,
						Description: "",
						Value:       "Jazz+Funk",
					},
					{
						Header:      ID3FrameHeader{"TIT2", 43, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoUTF16BOMEncoding,
						Description: "",
						Value:       "Funky weekend night",
					},
					{
						Header:      ID3FrameHeader{"TRCK", 3, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoISO88591,
						Description: "",
						Value:       "1",
					},
					{
						Header:      ID3FrameHeader{"TYER", 6, ID3FrameStatusFlag{false, false, false}, ID3FrameFormatFlag{false, false, false, false, false}},
						Encoding:    textInfoISO88591,
						Description: "",
						Value:       "2022",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		f, err := os.Open(tc.filePath)
		if err != nil {
			t.Fatalf("Cannot open test audio file. err: %s\n", err.Error())
		}

		data := make([]byte, id3HeaderSize+tc.want.Header.TagSize)
		if _, err := f.Read(data); err != nil {
			t.Fatalf("Cannot read test audio file. err: %s\n", err.Error())
		}

		got, _ := readID3Tag(data, 0)
		want := tc.want

		if got.Header.Version != want.Header.Version {
			t.Errorf("ID3 version is incorrect.\ngot: %s\nwant: %s\n", got.Header.Version, want.Header.Version)
		}
		if got.Header.TagSize != want.Header.TagSize {
			t.Errorf("ID3 tag size is incorrect.\ngot: %d\nwant: %d\n", got.Header.TagSize, want.Header.TagSize)
		}
		if !reflect.DeepEqual(got.TextInfoFrames, want.TextInfoFrames) {
			t.Errorf("ID3 text info frames are incorrect.\ngot: %#v\nwant: %#v\n", got.TextInfoFrames, want.TextInfoFrames)
		}
	}
}

func TestParseAudio_readVBRHeader(t *testing.T) {
	f, err := os.Open("./testdata/ramp_jazz.mp3")
	if err != nil {
		t.Fatalf("Cannot open test audio file. err: %s\n", err.Error())
	}

	var vbrHeaderOffset int64 = 131
	if _, err := f.Seek(vbrHeaderOffset, 0); err != nil {
		t.Fatalf("Cannot skip to the first frame sync. err: %s\n", err.Error())
	}
	data := make([]byte, 120)
	if _, err := f.Read(data); err != nil {
		t.Fatalf("Cannot read test audio file. err: %s\n", err.Error())
	}

	got, _ := readVBRHeader(data, 0)
	wantID := vbrInfoMagic
	wantNumOfFrames, wantFileSize, wantQuality := 7572, 7913011, 0
	wantTOC := []int{
		0, 3, 5, 8, 10, 13, 16, 18, 20, 23,
		25, 28, 31, 33, 36, 38, 41, 44, 46, 48,
		51, 54, 56, 59, 61, 64, 67, 69, 72, 74,
		76, 80, 82, 84, 87, 89, 91, 95, 97, 100,
		102, 104, 108, 110, 112, 115, 117, 120, 123, 125,
		128, 130, 133, 136, 138, 140, 143, 146, 148, 151,
		153, 155, 159, 161, 164, 166, 168, 172, 174, 176,
		179, 181, 183, 187, 189, 192, 194, 196, 200, 202,
		204, 207, 209, 212, 215, 217, 219, 222, 225, 228,
		230, 232, 235, 238, 240, 243, 245, 247, 250, 253,
	}

	if got.NumOfFrames == nil {
		t.Fatalf("VBR number of frames is not read.\n")
	}
	if got.FileSize == nil {
		t.Fatalf("VBR file size is not read.\n")
	}
	if got.Quality == nil {
		t.Fatalf("VBR quality is not read.\n")
	}

	if got.ID != wantID {
		t.Errorf("VBR ID is incorrect.\ngot: %s\nwant: %s\n", got.ID, wantID)
	}
	if *got.NumOfFrames != wantNumOfFrames {
		t.Errorf("VBR number of frames is incorrect.\ngot: %v\nwant: %v\n", *got.NumOfFrames, wantNumOfFrames)
	}
	if *got.FileSize != wantFileSize {
		t.Errorf("VBR file size is incorrect.\ngot: %v\nwant: %v\n", *got.FileSize, wantFileSize)
	}
	if !reflect.DeepEqual(got.TOC, wantTOC) {
		t.Errorf("VBR TOC are incorrect.\ngot: %v\nwant: %v\n", got.TOC, wantTOC)
	}
	if *got.Quality != wantQuality {
		t.Errorf("VBR quality is incorrect.\ngot: %v\nwant: %v\n", *got.Quality, wantQuality)
	}
}

func TestParseAudio_readMPEGAudioFrameHeader(t *testing.T) {
	f, err := os.Open("./testdata/ramp_jazz.mp3")
	if err != nil {
		t.Fatalf("Cannot open test audio file. err: %s\n", err.Error())
	}

	var firstFrameSyncOffset int64 = 1139
	if _, err := f.Seek(firstFrameSyncOffset, 0); err != nil {
		t.Fatalf("Cannot skip to the first frame sync. err: %s\n", err.Error())
	}
	data := make([]byte, 4)
	if _, err := f.Read(data); err != nil {
		t.Fatalf("Cannot read test audio file. err: %s\n", err.Error())
	}

	got, _ := readMPEGAudioFrameHeader(data, 0)
	want := MPEGAudioFrameHeader{
		MPEGAudioVersion: "1",
		Layer:            3,
		Protected:        false,
		Bitrate:          320,
		SampleRate:       44100,
	}

	if got.MPEGAudioVersion != want.MPEGAudioVersion {
		t.Errorf("MPEG audio version is incorrect.\ngot: %s\nwant: %s\n", got.MPEGAudioVersion, want.MPEGAudioVersion)
	}
	if got.Layer != want.Layer {
		t.Errorf("MPEG audio layer is incorrect.\ngot: %d\nwant: %d\n", got.Layer, want.Layer)
	}
	if got.Protected != want.Protected {
		t.Errorf("MPEG audio CRC protection is incorrect.\ngot: %v\nwant: %v\n", got.Protected, want.Protected)
	}
	if got.Bitrate != want.Bitrate {
		t.Errorf("MPEG audio bitrate is incorrect.\ngot: %d\nwant: %d\n", got.Bitrate, want.Bitrate)
	}
	if got.SampleRate != want.SampleRate {
		t.Errorf("MPEG audio sample rate is incorrect.\ngot: %d\nwant: %d\n", got.SampleRate, want.SampleRate)
	}
}
