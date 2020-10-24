package utils

import "testing"

func TestOutputMp4(t *testing.T) {
	if err := ConcatVideo(
		"/Users/Downloads/test/filelist.txt",
		"/Users/Downloads/test.mp4",
	); err != nil {
		t.Error(err)
	}
}
