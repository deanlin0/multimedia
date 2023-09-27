package hls

import (
	"os"
	"reflect"
	"testing"
)

func TestParseAudio_parseID3Header(t *testing.T) {
	f, err := os.Open("./testdata/ramp_jazz.mp3")
	if err != nil {
		t.Fatalf("Cannot open test audio file. err: %s\n", err.Error())
	}

	want := AudioContext{
		ID3Header: ID3Header{
			Version: "2.4.0",
			TagSize: 85,
			TextInfoFrames: []TextInfoFrame{
				{
					"TDRC", 12, textInfoUTF8Encoding, "", "2022-11-16",
					FrameStatusFlag{false, false, false},
					FrameFormatFlag{false, false, false, false, false},
				},
				{
					"TXXX", 18, textInfoUTF8Encoding, "time_reference", "0",
					FrameStatusFlag{false, false, false},
					FrameFormatFlag{false, false, false, false, false},
				},
				{
					"TSSE", 15, textInfoUTF8Encoding, "", "Lavf59.27.100",
					FrameStatusFlag{false, false, false},
					FrameFormatFlag{false, false, false, false, false},
				},
			},
		},
	}
	got := AudioContext{}

	if err := parseID3Header(f, &got); err != nil {
		t.Errorf("Failed to parse audio. err: %s\n", err.Error())
	}
	if got.ID3Header.Version != want.ID3Header.Version {
		t.Errorf("ID3 version is incorrect.\ngot: %s\nwant: %s\n", got.ID3Header.Version, want.ID3Header.Version)
	}
	if got.ID3Header.TagSize != want.ID3Header.TagSize {
		t.Errorf("ID3 tag size is incorrect.\ngot: %d\nwant: %d\n", got.ID3Header.TagSize, want.ID3Header.TagSize)
	}
	if !reflect.DeepEqual(got.ID3Header.TextInfoFrames, want.ID3Header.TextInfoFrames) {
		t.Errorf("ID3 text info frames are incorrect.\ngot: %#v\nwant: %#v\n", got.ID3Header.TextInfoFrames, want.ID3Header.TextInfoFrames)
	}
}
