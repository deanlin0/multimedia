package hls

import (
	"os"
	"reflect"
	"testing"
)

func TestParseAudio_parseID3Tag(t *testing.T) {
	testCases := []struct {
		filePath string
		want     AudioContext
	}{
		{
			filePath: "./testdata/ramp_jazz.mp3",
			want: AudioContext{
				ID3Tag: ID3Tag{
					Version: "2.4.0",
					Size:    85,
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
		},
		{
			filePath: "./testdata/funky_weekend_night.mp3",
			want: AudioContext{
				ID3Tag: ID3Tag{
					Version: "2.3.0",
					Size:    640733,
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
		},
	}

	for _, tc := range testCases {
		f, err := os.Open(tc.filePath)
		if err != nil {
			t.Fatalf("Cannot open test audio file. err: %s\n", err.Error())
		}

		got := AudioContext{}
		want := tc.want

		if err := parseID3Tag(f, &got); err != nil {
			t.Errorf("Failed to parse audio. err: %s\n", err.Error())
		}
		if got.ID3Tag.Version != want.ID3Tag.Version {
			t.Errorf("ID3 version is incorrect.\ngot: %s\nwant: %s\n", got.ID3Tag.Version, want.ID3Tag.Version)
		}
		if got.ID3Tag.Size != want.ID3Tag.Size {
			t.Errorf("ID3 tag size is incorrect.\ngot: %d\nwant: %d\n", got.ID3Tag.Size, want.ID3Tag.Size)
		}
		if !reflect.DeepEqual(got.ID3Tag.TextInfoFrames, want.ID3Tag.TextInfoFrames) {
			t.Errorf("ID3 text info frames are incorrect.\ngot: %#v\nwant: %#v\n", got.ID3Tag.TextInfoFrames, want.ID3Tag.TextInfoFrames)
		}
	}
}
