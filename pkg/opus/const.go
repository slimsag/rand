package opus

/*
#include <opus.h>
*/
import "C"

import (
	"errors"
	"fmt"
)

var (
	ErrBadArg         = errors.New("opus: one or more invalid / out of range arguments")
	ErrBufferTooSmall = errors.New("opus: the mode struct passed is invalid")
	ErrInternalError  = errors.New("opus: internal error")
	ErrInvalidPacket  = errors.New("opus: the compressed data passed is corrupted")
	ErrUnimplemented  = errors.New("opus: invalid / unsupported request number")
	ErrInvalidState   = errors.New("opus: an encoder or decoder structure is invalid")
	ErrAllocFail      = errors.New("opus: memory allocation has failed")
)

func opusError(code C.int) error {
	switch code {
	case C.OPUS_OK:
		return nil
	case C.OPUS_BAD_ARG:
		return ErrBadArg
	case C.OPUS_BUFFER_TOO_SMALL:
		return ErrBufferTooSmall
	case C.OPUS_INTERNAL_ERROR:
		return ErrInternalError
	case C.OPUS_INVALID_PACKET:
		return ErrInvalidPacket
	case C.OPUS_UNIMPLEMENTED:
		return ErrUnimplemented
	case C.OPUS_INVALID_STATE:
		return ErrInvalidState
	case C.OPUS_ALLOC_FAIL:
		return ErrAllocFail
	default:
		panic(fmt.Sprintf("unkown error code (error %d / 0x%X)", code, code))
	}
}

// Pre-defined values for CTL interface
const (
	// Values for the various encoder CTLs

	// Auto/default setting
	AUTO = C.OPUS_AUTO

	// Maximum bitrate
	BITRATE_MAX = C.OPUS_BITRATE_MAX

	// Best for most VoIP/videoconference applications where listening quality
	// and intelligibility matter most.
	APPLICATION_VOIP = C.OPUS_APPLICATION_VOIP

	// Best for broadcast/high-fidelity application where the decoded audio
	// should be as close as possible to the input.
	APPLICATION_AUDIO = C.OPUS_APPLICATION_AUDIO

	// Only use when lowest-achievable latency is what matters most.
	// Voice-optimized modes cannot be used.
	APPLICATION_RESTRICTED_LOWDELAY = C.OPUS_APPLICATION_RESTRICTED_LOWDELAY

	// Signal being encoded is voice
	SIGNAL_VOICE = C.OPUS_SIGNAL_VOICE

	// Signal being encoded is music
	SIGNAL_MUSIC = C.OPUS_SIGNAL_MUSIC

	// 4 kHz bandpass
	BANDWIDTH_NARROWBAND = C.OPUS_BANDWIDTH_NARROWBAND

	// 6 kHz bandpass
	BANDWIDTH_MEDIUMBAND = C.OPUS_BANDWIDTH_MEDIUMBAND

	// 8 kHz bandpass
	BANDWIDTH_WIDEBAND = C.OPUS_BANDWIDTH_WIDEBAND

	// 12 kHz bandpass
	BANDWIDTH_SUPERWIDEBAND = C.OPUS_BANDWIDTH_SUPERWIDEBAND

	// 20 kHz bandpass
	BANDWIDTH_FULLBAND = C.OPUS_BANDWIDTH_FULLBAND

	// Select frame size from the argument (default)
	FRAMESIZE_ARG = C.OPUS_FRAMESIZE_ARG

	// Use 2.5 ms frames
	FRAMESIZE_2_5_MS = C.OPUS_FRAMESIZE_2_5_MS

	// Use 5 ms frames
	FRAMESIZE_5_MS = C.OPUS_FRAMESIZE_5_MS

	// Use 10 ms frames
	FRAMESIZE_10_MS = C.OPUS_FRAMESIZE_10_MS

	// Use 20 ms frames
	FRAMESIZE_20_MS = C.OPUS_FRAMESIZE_20_MS

	// Use 40 ms frames
	FRAMESIZE_40_MS = C.OPUS_FRAMESIZE_40_MS

	// Use 60 ms frames
	FRAMESIZE_60_MS = C.OPUS_FRAMESIZE_60_MS
)
