package resizeVideo

import "testing"

func TestResizeVideo(t *testing.T) {
	src := "/Users/zen/Downloads/DHAI-152/QOC/honoka"
	dst := "/Users/zen/Downloads/DHAI-152/QOC/honoka/h265"
	pattern := "mp4"
	threads := "2"
	isDelete := true
	ResizeVideo(src, dst, pattern, threads, isDelete)
}
