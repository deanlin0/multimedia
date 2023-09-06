package hls

import "fmt"

const (
	EXTM3U               = "EXTM3U"
	EXT_X_TARGETDURATION = "EXT-X-TARGETDURATION"
	EXTINF               = "EXTINF"
	EXT_X_VERSION        = "EXT-X-VERSION"
	EXT_X_KEY            = "EXT-X-KEY"
)

func BuildPlainTagLine(tag string) string {
	return fmt.Sprintf("#%s", tag)
}

func BuildNumberTagLine(tag string, number float64) string {
	return fmt.Sprintf("#%s:%g", tag, number)
}
