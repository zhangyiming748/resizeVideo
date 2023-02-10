package resizeVideo

import "testing"

func TestResizeVideo(t *testing.T) {
	src := "/Users/zen/Downloads/整理/dance/梓/Left"
	pattern := "mp4"
	threads := "10"
	isDelete := false
	ResizeVideo(src, pattern, threads, isDelete)
}
