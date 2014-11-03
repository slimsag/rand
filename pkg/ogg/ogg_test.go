package ogg

import (
	"os"
	"testing"
)

func TestOgg(t *testing.T) {
	f, err := os.Open("testdata/tune_stereo_44100hz_vorbis.ogg")
	if err != nil {
		t.Fatal(err)
	}
	dec := NewDecoder(f)

	for {
		// Decode a single packet.
		packet, err := dec.Decode()
		if packet == nil {
			// We've encountered the end of the stream and the packet is nil.
			break
		}
		if err != nil {
			// We've encountered a fatal error while decoding the stream.
			t.Fatal(err)
		}

		t.Log("Packet spans across", len(packet.Pages), "pages.")
		for _, page := range packet.Pages {
			t.Log(string(page.CapturePattern[:]))
			t.Logf("%+v\n", page)
		}
	}
}
