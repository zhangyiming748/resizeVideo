package resizeVideo

import "testing"

func TestResizeVideo(t *testing.T) {
	src := "/Users/zen/Downloads/整理/ff/Aerith"
	pattern := "mp4;mkv"
	threads := "10"
	ResizeVideos(src, pattern, threads)
}
func TestResizeAllVideos(t *testing.T) {
	src := "/Users/zen/Downloads/整理"
	pattern := "mp4;mkv"
	threads := "4"
	ResizeAllVideos(src, pattern, threads)
}
