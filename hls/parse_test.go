package hls

import (
	"os"
	"reflect"
	"testing"
)

func TestParseAudio_parseID3Tag(t *testing.T) {
	f, err := os.Open("./testdata/ramp_jazz.mp3")
	if err != nil {
		t.Fatalf("Cannot open test audio file. err: %s\n", err.Error())
	}

	want := AudioContext{
		ID3Tag: ID3Tag{
			Version: "2.4.0",
			Size:    85,
			TextInfoFrames: []TextInfoFrame{
				{
					Header:      FrameHeader{"TDRC", 12, FrameStatusFlag{false, false, false}, FrameFormatFlag{false, false, false, false, false}},
					Encoding:    textInfoUTF8Encoding,
					Description: "",
					Value:       "2022-11-16",
				},
				{
					Header:      FrameHeader{"TXXX", 18, FrameStatusFlag{false, false, false}, FrameFormatFlag{false, false, false, false, false}},
					Encoding:    textInfoUTF8Encoding,
					Description: "time_reference",
					Value:       "0",
				},
				{
					Header:      FrameHeader{"TSSE", 15, FrameStatusFlag{false, false, false}, FrameFormatFlag{false, false, false, false, false}},
					Encoding:    textInfoUTF8Encoding,
					Description: "",
					Value:       "Lavf59.27.100",
				},
			},
		},
	}
	got := AudioContext{}

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
