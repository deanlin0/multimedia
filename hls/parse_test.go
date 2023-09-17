package hls

import (
	"os"
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
}
