package resizeVideo

import "testing"

func TestResizeVideo(t *testing.T) {
	src := "/Users/zen/Downloads/整理/ff/Aerith"
	pattern := "mp4;mkv"
	threads := "10"
	ResizeVideos(src, pattern, threads)
}
func TestResizeAllVideos(t *testing.T) {
	src := "/Users/zen/Downloads/NecDaz"
	pattern := "mp4;mkv"
	threads := "8"
	ResizeAllVideos(src, pattern, threads)
}
func TestFixAll4x3s(t *testing.T) {
	src := "/Users/zen/Downloads/4x3"
	pattern := "mp4;mkv"
	threads := "8"
	FixAll4x3s(src, pattern, threads)
}
