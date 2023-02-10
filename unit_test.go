package resizeVideo

import "testing"

func TestResizeVideo(t *testing.T) {
	src := "/Users/zen/Downloads/整理/dance/梦梦"
	pattern := "mp4"
	threads := "10"
	isDelete := true
	ResizeVideo(src, pattern, threads, isDelete)
}
