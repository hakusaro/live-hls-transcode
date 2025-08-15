package main

import "testing"

func TestEscapeFFmpegPath(t *testing.T) {
	input := "/path/with space/file's.mp4"
	expected := `/path/with\ space/file\'s.mp4`
	if got := escapeFFmpegPath(input); got != expected {
		t.Errorf("escapeFFmpegPath(%q) = %q, want %q", input, got, expected)
	}
}
