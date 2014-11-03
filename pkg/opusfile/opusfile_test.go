package opusfile

import (
	"os"
	"testing"
)

func TestOpus(t *testing.T) {
	f, err := os.Open("testdata/tune_stereo_44100hz_vorbis.ogg")
	if err != nil {
		t.Fatal(err)
	}

	opusFile, err := Open(f)
	if err != nil {
		t.Fatal(err)
	}
	_ = opusFile
}
