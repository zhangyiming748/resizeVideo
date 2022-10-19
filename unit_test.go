package resizeVideo

import "testing"

func TestResizeVideo(t *testing.T) {
	src := "/Users/zen/Github/resizeVideo/DB"
	dst := "/Users/zen/Github/resizeVideo/DB/1080"
	pattern := "mp4"
	threads := "4"
	isDelete := true
	ResizeVideo(src, dst, pattern, threads, isDelete)
}
