// Package opus implements a Go wrapper to libopus.
//
// This is a Go wrapper (via CGO) to libopus. More information about libopus
// and the Opus audio codec can be found at:
//
// http://www.opus-codec.org/
//
// The Opus codec is designed for interactive speech and audio transmission
// over the Internet. It is designed by the IETF Codec Working Group and
// incorporates technology from Skype's SILK codec and Xiph.Org's CELT codec.
//
// The Opus codec is designed to handle a wide range of interactive audio
// applications, including Voice over IP, videoconferencing, in-game chat, and
// even remote live music performances. It can scale from low bit-rate
// narrowband speech to very high quality stereo music. Its main features are:
//  Sampling rates from 8 to 48 kHz
//  Bit-rates from 6 kb/s to 510 kb/s
//  Support for both constant bit-rate (CBR) and variable bit-rate (VBR)
//  Audio bandwidth from narrowband to full-band
//  Support for speech and music
//  Support for mono and stereo
//  Support for multichannel (up to 255 channels)
//  Frame sizes from 2.5 ms to 60 ms
//  Good loss robustness and packet loss concealment (PLC)
//  Floating point and fixed-point implementation
package opus

/*
#cgo linux LDFLAGS: -lm
#cgo CFLAGS: -Iinclude/
*/
import "C"

func boolToInt32(x bool) int32 {
	if x {
		return 1
	}
	return 0
}

func int32ToBool(x int32) bool {
	if x == 1 {
		return true
	}
	return false
}
