package hls

import (
	"reflect"
	"testing"
)

func TestPlaylist_SimplePlaylist(t *testing.T) {
	// Create playlist
	playlist := &Playlist{}

	// Add playlist lines
	lines := []string{
		BuildPlainTagLine(EXTM3U),
		BuildNumberTagLine(EXT_X_TARGETDURATION, 10),
		BuildSegmentLine(Segment{URI: "http://media.example.com/first.ts", Duration: 9.009}),
		BuildSegmentLine(Segment{URI: "http://media.example.com/second.ts", Duration: 9.009}),
		BuildSegmentLine(Segment{URI: "http://media.example.com/third.ts", Duration: 3.003}),
	}
	for _, line := range lines {
		playlist.AddLine(line)
	}

	got := playlist.File()
	want := []byte(`#EXTM3U
#EXT-X-TARGETDURATION:10
#EXTINF:9.009,
http://media.example.com/first.ts
#EXTINF:9.009,
http://media.example.com/second.ts
#EXTINF:3.003,
http://media.example.com/third.ts`)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Playlist file content is incorrect.\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}
