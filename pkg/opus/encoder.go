package opus

/*
#cgo linux LDFLAGS: -lm
#cgo CFLAGS: -Iinclude/

#include <opus.h>

int azul3d_opus_encoder_ctl_int32(OpusEncoder* st, int request, opus_int32 arg) {
	return opus_encoder_ctl(st, request, arg);
}

int azul3d_opus_encoder_ctl_int32_ptr(OpusEncoder* st, int request, opus_int32* arg) {
	return opus_encoder_ctl(st, request, arg);
}
*/
import "C"

import "unsafe"

// Opus encoder state. This contains the complete state of an Opus encoder. It
// is position independent and can be freely copied.
type Encoder C.OpusEncoder

func (e *Encoder) cptr() *C.OpusEncoder {
	return (*C.OpusEncoder)(unsafe.Pointer(e))
}

func (e *Encoder) ctl_int32(request, arg int32) int {
	return int(C.azul3d_opus_encoder_ctl_int32(
		e.cptr(),
		C.int(request),
		C.opus_int32(arg),
	))
}

func (e *Encoder) ctl_int32_ptr(request int, arg *int32) int {
	return int(C.azul3d_opus_encoder_ctl_int32_ptr(
		e.cptr(),
		C.int(request),
		(*C.opus_int32)(unsafe.Pointer(arg)),
	))
}

// Encode encodes an Opus frame.
//
// pcm: the 16-bit PCM input signal (interleaved if 2 channels). Whose length is
// frameSize*channels.
//
// frameSize: the number of samples per channel in the input signal. This must
// be an Opus frame size for the encoder's sampling rate. For example at 48kHz
// the permitted values are 120, 240, 480, 960, 1920, and 2880. Passing in a
// duration of less than 10 ms (480 samples at 48 kHz) will prevent the encoder
// from using the LPC or hybrid modes.
//
// data: the slice to store the output payload in, containing len(data) bytes
// of space. The length of the slice may be used to impose an upper limit on
// the instant bitrate, but should not be used as the only bitrate control. Use
// SET_BITRATE to control the bitrate.
//
// returns: a slice of the data slice representing the encoded packet, OR nil
// and an error.
func (e *Encoder) Encode(pcm []int16, frameSize int, data []byte) ([]byte, error) {
	n := C.opus_encode(
		e.cptr(),
		(*C.opus_int16)(unsafe.Pointer(&pcm[0])),
		C.int(frameSize),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.opus_int32(len(data)),
	)
	if n > 0 {
		return data[:n], nil
	}
	return nil, opusError(C.int(-n))
}

// EncodeFloat encodes an Opus frame from floating point input.
//
// pcm: the 32-bit float PCM input signal (interleaved if 2 channels). Whose
// length is frameSize*channels. Samples have a normal range of +/-1.0, but can
// have a range beyond that (but they will be clipped by decoders using the
// integer API and should only be used if it is known that the far end supports
// extended dynamic range).
//
// frameSize: the number of samples per channel in the input signal. This must
// be an Opus frame size for the encoder's sampling rate. For example at 48kHz
// the permitted values are 120, 240, 480, 960, 1920, and 2880. Passing in a
// duration of less than 10 ms (480 samples at 48 kHz) will prevent the encoder
// from using the LPC or hybrid modes.
//
// data: the slice to store the output payload in, containing len(data) bytes
// of space. The length of the slice may be used to impose an upper limit on
// the instant bitrate, but should not be used as the only bitrate control. Use
// SET_BITRATE to control the bitrate.
//
// returns: a slice of the data slice representing the encoded packet, OR nil
// and an error.
func (e *Encoder) EncodeFloat(pcm []float32, frameSize int, data []byte) ([]byte, error) {
	n := C.opus_encode_float(
		e.cptr(),
		(*C.float)(unsafe.Pointer(&pcm[0])),
		C.int(frameSize),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.opus_int32(len(data)),
	)
	if n > 0 {
		return data[:n], nil
	}
	return nil, opusError(C.int(-n))
}

// Allocates and initializes an encoder state. There are three coding modes:
//
//  1. APPLICATION_VOIP gives best quality at a given bitrate for voice
//  signals. It enhances the  input signal by high-pass filtering and
//  emphasizing formants and harmonics. Optionally  it includes in-band forward
//  error correction to protect against packet loss. Use this mode for typical
//  VoIP applications. Because of the enhancement, even at high bitrates the
//  output may sound different from the input.
//
//  2. APPLICATION_AUDIO gives best quality at a given bitrate for most
//  non-voice signals like music. Use this mode for music and mixed
//  (music/voice) content, broadcast, and applications requiring less than 15
//  ms of coding delay.
//
//  3. APPLICATION_RESTRICTED_LOWDELAY configures low-delay mode that disables
//  the speech-optimized mode in exchange for slightly reduced delay. This mode
//  can only be set on an newly initialized or freshly reset encoder because it
//  changes the codec delay. This is useful when the caller knows that the
//  speech-optimized modes will not be needed (use with caution).
//
// Fs: Sampling rate of input signal (Hz). This must be one of 8000, 12000,
// 16000, 24000, or 48000.
//
// channels: Number of channels (1 or 2) in input signal
//
// application: Coding mode (APPLICATION_VOIP/APPLICATION_AUDIO/APPLICATION_RESTRICTED_LOWDELAY)
//
// note: Regardless of the sampling rate and number channels selected, the Opus
// encoder can switch to a lower audio bandwidth or number of channels if the
// bitrate selected is too low. This also means that it is safe to always use
// 48 kHz stereo input and let the encoder optimize the encoding.
func NewEncoder(Fs, channels, application int) (*Encoder, error) {
	enc := new(Encoder)
	err := opusError(C.opus_encoder_init(
		enc.cptr(),
		C.opus_int32(Fs),
		C.int(channels),
		C.int(application),
	))
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Encoder CTL API's below this point.

// SetComplexity configures the encoder's computational complexity. The
// supported range is 0-10 inclusive with 10 representing the highest
// complexity.
func (e *Encoder) SetComplexity(x int) {
	e.ctl_int32(C.OPUS_SET_COMPLEXITY_REQUEST, int32(x))
}

// Complexity returns the encoder's complexity configuration.
func (e *Encoder) Complexity() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_COMPLEXITY_REQUEST, &x)
	return int(x)
}

// SetBitrate configures the bitrate of the encoder. Rates from 500 to 512000
// bits per second are meaningful, as well as the special values AUTO and
// BITRATE_MAX.
//
// The value BITRATE_MAX can be used to cause the codec to use as much rate as
// it can, which is useful for controlling the rate by adjusting the output
// buffer size.
//
// The default bitrate (in bits per second) is determined based on the number
// of channels and the input sampling rate.
func (e *Encoder) SetBitrate(bitrate int) {
	e.ctl_int32(C.OPUS_SET_BITRATE_REQUEST, int32(bitrate))
}

// Bitrate returns the encoder's bitrate configuration.
func (e *Encoder) Bitrate() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_BITRATE_REQUEST, &x)
	return int(x)
}

// SetVBR enables or disables variable bitrate (VBR) in the encoder. The
// configured bitrate may not be met exactly because frames must be an integer
// number of bytes in length.
//
// Warning: Only the MDCT mode of Opus can provide hard CBR behavior.
//
// If vbr is false: Hard CBR. For LPC/hybrid modes at very low bit-rate, this
// can cause noticable quality degradation.
//
// If vbr is true: VBR (default). The exact type of VBR is controlled by the
// SetVBRConstraint method.
func (e *Encoder) SetVBR(vbr bool) {
	e.ctl_int32(C.OPUS_SET_VBR_REQUEST, boolToInt32(vbr))
}

// VBR returns the encoder's VBR configuration.
func (e *Encoder) VBR() bool {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_VBR_REQUEST, &x)
	return int32ToBool(x)
}

// SetVBRConstraint enables or disables constrained VBR in the encoder. This
// setting is ignored when the encoder is in CBR mode.
//
// Warning: Only the MDCT mode of Opus currently heeds the constraint. Speech
// mode ignores it completely, hybrid mode may fail to obey it if the LPC layer
// uses more bitrate than the constraint would have permitted.
//
// If constrained is true: Constrained VBR (default). This creates a maximum of
// one frame of buffering delay assuming a transport with serialization speed
// of the nominal bitrate.
//
// If constrained is false: Unconstrained VBR.
func (e *Encoder) SetVBRConstraint(constrained bool) {
	e.ctl_int32(C.OPUS_SET_VBR_CONSTRAINT_REQUEST, boolToInt32(constrained))
}

// VBRConstraint returns the encoder's VBR constraint configuration.
func (e *Encoder) VBRConstraint() bool {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_VBR_CONSTRAINT_REQUEST, &x)
	return int32ToBool(x)
}

// SetForceChannels configures mono/stereo forcing in the encoder. This can
// force the encoder to produce packets encoded as either mono or stereo,
// regardless of the format of the input audio. This is useful when the caller
// knows that the input signal is currently a mono source embedded in a stereo
// stream.
//
// Allowed values are AUTO (default value: not forced), 1 (forced mono), and 2
// (forced stereo).
func (e *Encoder) SetForceChannels(channels int) {
	e.ctl_int32(C.OPUS_SET_FORCE_CHANNELS_REQUEST, int32(channels))
}

// ForceChannels returns the encoder's force channels configuration.
func (e *Encoder) ForceChannels() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_FORCE_CHANNELS_REQUEST, &x)
	return int(x)
}

// SetMaxBandwidth configures the maximum bandpass that the encoder will select
// automatically. Applications should normally use this instead of SetBandwidth
// (leaving that set to the default, AUTO). This allows the application to set
// an upper bound based on the type of input it is providing, but still gives
// the encoder the freedom to reduce the bandpass when the bitrate becomes too
// low, for better overall quality.
//
// Allowed values are:
//  BANDWIDTH_NARROWBAND     4 kHz passband
//  BANDWIDTH_MEDIUMBAND     6 kHz passband
//  BANDWIDTH_WIDEBAND       8 kHz passband
//  BANDWIDTH_SUPERWIDEBAND 12 kHz passband
//  BANDWIDTH_FULLBAND      20 kHz passband (default)
func (e *Encoder) SetMaxBandwidth(x int) {
	e.ctl_int32(C.OPUS_SET_MAX_BANDWIDTH_REQUEST, int32(x))
}

// MaxBandwidth returns the encoder's max bandwidth configuration.
func (e *Encoder) MaxBandwidth() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_MAX_BANDWIDTH_REQUEST, &x)
	return int(x)
}

// SetBandwidth sets the encoder's bandpass to a specific value. This prevents
// the encoder from automatically selecting the bandpass based on the available
// bitrate. If an application knows the bandpass of the input audio it is
// providing, it should normally use SetMaxBandwidth instead, which still gives
// the encoder the freedom to reduce the bandpass when the bitrate becomes too
// low, for better overall quality.
//
// Allowed values are:
//  AUTO (default)
//  BANDWIDTH_NARROWBAND     4 kHz passband
//  BANDWIDTH_MEDIUMBAND     6 kHz passband
//  BANDWIDTH_WIDEBAND       8 kHz passband
//  BANDWIDTH_SUPERWIDEBAND 12 kHz passband
//  BANDWIDTH_FULLBAND      20 kHz passband
func (e *Encoder) SetBandwidth(x int) {
	e.ctl_int32(C.OPUS_SET_BANDWIDTH_REQUEST, int32(x))
}

// Bandwidth returns the encoder's bandwidth configuration.
func (e *Encoder) Bandwidth() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_BANDWIDTH_REQUEST, &x)
	return int(x)
}

// SetSignal configures the type of signal being encoded. This is a hint which
// helps the encoder's mode selection.
//
// Allowed values are:
//  AUTO          (default)
//  SIGNAL_VOICE  Bias thresholds towards choosing LPC or Hybrid modes.
//  SIGNAL_MUSIC  Bias thresholds towards choosing MDCT modes.
func (e *Encoder) SetSignal(x int) {
	e.ctl_int32(C.OPUS_SET_SIGNAL_REQUEST, int32(x))
}

// Signal returns the encoder's signal configuration.
func (e *Encoder) Signal() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_SIGNAL_REQUEST, &x)
	return int(x)
}

// SetApplication configures the encoder's intended application. The initial
// value is a mandatory argument to the NewEncoder function.
//
// Allowed values are:
//  APPLICATION_VOIP - Process signal for improved speech intelligibility.
//
//  APPLICATION_AUDIO - Favor faithfulness to the original input.
//
//  APPLICATION_RESTRICTED_LOWDELAY - Configure the minimum possible coding
//  delay by disabling certain modes of operation.
func (e *Encoder) SetApplication(x int) {
	e.ctl_int32(C.OPUS_SET_APPLICATION_REQUEST, int32(x))
}

// Application returns the encoder's application configuration.
func (e *Encoder) Application() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_APPLICATION_REQUEST, &x)
	return int(x)
}

// SampleRate returns the sampling rate that the encoder was initialized with.
// This simply returns the Fs value passed to NewEncoder.
func (e *Encoder) SampleRate() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_SAMPLE_RATE_REQUEST, &x)
	return int(x)
}

// Lookahead returns the total samples of delay added by the entire codec. This
// can be queried by the encoder and then the provided number of samples can be
// skipped on from the start of the decoder's output to provide time aligned
// input and output. From the perspective of a decoding application the real
// data begins this many samples late.
//
// The decoder contribution to this delay is identical for all decoders, but
// the encoder portion of the delay may vary from implementation to
// implementation, version to version, or even depend on the encoder's initial
// configuration.
//
// Applications needing delay compensation should call this rather than hard
// coding a value.
func (e *Encoder) Lookahead() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_LOOKAHEAD_REQUEST, &x)
	return int(x)
}

// SetInbandFEC configures the encoder's use of inband forward error correction
// (FEC). This is only applicable to the LPC layer.
//
// Allowed values are:
//  false - Disables inband FEC (default).
//  true - Enable inband FEC.
func (e *Encoder) SetInbandFEC(inbandFEC bool) {
	e.ctl_int32(C.OPUS_SET_INBAND_FEC_REQUEST, boolToInt32(inbandFEC))
}

// InbandFEC returns the encoder's inband forward error correction (FEC)
// configuration.
func (e *Encoder) InbandFEC() bool {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_INBAND_FEC_REQUEST, &x)
	return int32ToBool(x)
}

// SetPacketLossPerc configures the encoder's expected packet loss percentage.
// Higher values with trigger progressively more loss resistant behavior in the
// encoder at the expense of quality at a given bitrate in the lossless case,
// but greater quality under loss.
//
// Value is percentage in the range 0-100, inclusive, the default is zero.
func (e *Encoder) SetPacketLossPerc(percent int) {
	e.ctl_int32(C.OPUS_SET_PACKET_LOSS_PERC_REQUEST, int32(percent))
}

// PacketLossPerc returns the encoder's packet loss percentage configuration.
func (e *Encoder) PacketLossPerc() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_PACKET_LOSS_PERC_REQUEST, &x)
	return int(x)
}

// SetDTX configures the encoder's use of discontinuous transmission (DTX).
// This is only applicable to the LPC layer. By default DTX is disabled
// (false).
func (e *Encoder) SetDTX(dtx bool) {
	e.ctl_int32(C.OPUS_SET_DTX_REQUEST, boolToInt32(dtx))
}

// DTX returns the encoder's discontinuous transmission (DTX) configuration.
func (e *Encoder) DTX() bool {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_DTX_REQUEST, &x)
	return int32ToBool(x)
}

// SetLSBDepth configures the depth of signal being encoded. This is a hint
// which helps the encoder identify silence and near-silence. Acceptable values
// are between 8 and 24, with the default being 24.
func (e *Encoder) SetLSBDepth(depth int) {
	e.ctl_int32(C.OPUS_SET_LSB_DEPTH_REQUEST, int32(depth))
}

// LSBDepth returns the encoder's signal depth configuration.
func (e *Encoder) LSBDepth() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_LSB_DEPTH_REQUEST, &x)
	return int(x)
}

// SetExpertFrameDuration configures the encoder's use of variable duration
// frames. When variable duration is enabled, the encoder is free to use a
// shorter frame size than the one requested in the Encode() methods. It is the
// users responsability to verify how much audio was encoded by checking the
// ToC byte of the encoded packet. The part of the audio that was not encoded
// needs to be resent to the encoder for the next call. Do not use this option
// unless you really know what you are doing.
//
// Allowed values are:
//  FRAMESIZE_ARG      - Select frame size from the argument (default).
//  FRAMESIZE_2_5_MS   - Use 2.5 ms frames.
//  FRAMESIZE_5_MS     - Use 5 ms frames.
//  FRAMESIZE_10_MS    - Use 10 ms frames.
//  FRAMESIZE_20_MS    - Use 20 ms frames.
//  FRAMESIZE_40_MS    - Use 40 ms frames.
//  FRAMESIZE_60_MS    - Use 60 ms frames.
//  FRAMESIZE_VARIABLE - Optimize the frame size dynamically.
func (e *Encoder) SetExpertFrameDuration(fs int) {
	e.ctl_int32(C.OPUS_SET_EXPERT_FRAME_DURATION_REQUEST, int32(fs))
}

// ExpertFrameDuration returns the encoder's variable duration frame
// configuration.
func (e *Encoder) ExpertFrameDuration() int {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_EXPERT_FRAME_DURATION_REQUEST, &x)
	return int(x)
}

// SetPredictionDisabled disables almost all use of prediction, making frames
// almost completely independent. This reduces quality. (default:
// false/enabled).
func (e *Encoder) SetPredictionDisabled(disabled bool) {
	e.ctl_int32(C.OPUS_SET_PREDICTION_DISABLED_REQUEST, boolToInt32(disabled))
}

// PredictionDisabled returns the encoder's prediction disabled status.
func (e *Encoder) PredictionDisabled() bool {
	var x int32
	e.ctl_int32_ptr(C.OPUS_GET_PREDICTION_DISABLED_REQUEST, &x)
	return int32ToBool(x)
}
