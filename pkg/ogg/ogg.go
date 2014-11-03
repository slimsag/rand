package ogg

import (
	"azul3d.org/v1/audio"
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
)

const (
	// If set the page contains data of a packet continued from the previous
	// page. If unset then the page contains a fresh packet.
	Continuation uint8 = 0x01

	// If set the page is the first page of a logical bitstream (beginning of
	// stream). If unset the page is not a first page.
	BOS uint8 = 0x02

	// If set the page is the last page of a logical bitstream (end of stream).
	// If unset the page is not a last page.
	EOS uint8 = 0x04
)

// PageHeader describes the header of a single page in an ogg stream.
type PageHeader struct {
	// The capture pattern (magic) used to signal the start of a page. Is is
	// always the ASCII characters "OggS". It can be used to resynchronize the
	// stream when data has been lost (the parser can scan for the next page).
	CapturePattern [4]byte

	// 1 Byte signifying the version number of the Ogg file format used in this
	// stream (only version zero exists).
	Version uint8

	// The bits in this field identify the specific type of this page. Flags
	// are Continuation, BOS, and EOS.
	HeaderType uint8

	// Contains position information. For example, for an audio stream, it may
	// contain the total number of PCM samples encoded after including all
	// frames finished on this page. For a video stream it may contain the
	// total number of video frames encoded after this page. This is a hint for
	// the decoder and gives it some timing and position information. Its
	// meaning is dependent on the codec for that logical bitstream and
	// specified in a specific media mapping. A special value of -1 indicates
	// that no packets finish on this page.
	GranulePosition int64

	// The unique serial number by which the logical bitstream is identified.
	// This field is a serial number that identifies a page as belonging to a
	// particular logical bitstream. Each logical bitstream in a file has a
	// unique value, and this field allows implementations to deliver the pages
	// to the appropriate decoder. In a typical Vorbis and Theora file, one
	// stream is the audio (Vorbis), and the other is the video (Theora).
	BitstreamSerialNumber int32

	// The sequence number of the page so the decoder can identify page loss.
	// This sequence number is increasing on each logical bitstream seperately.
	PageSequenceNumber int32

	// 32-bit CRC checksum of the page (including header with zero CRC field
	// and page content). The generator polynomial is 0x04c11db7.
	Checksum int32

	// The number of segment entries encoded in the segment table.
	NumberPageSegments uint8

	//SegmentTable
}

// Page represents a single page in a ogg stream.
type Page struct {
	*PageHeader
	//	*SegmentTable
}

/*
   9. segment_table: number_page_segments Bytes containing the lacing
      values of all segments in this page.  Each Byte contains one
      lacing value.

   The total header size in bytes is given by:
   header_size = number_page_segments + 27 [Byte]

   The total page size in Bytes is given by:
   page_size = header_size + sum(lacing_values: 1..number_page_segments)
   [Byte]
*/

// Validate returns an error if any of the PageHeader fields are invalid with
// accordinace to RFC 3533. Otherwise (if no errors exist), nil is returned.
func (p *PageHeader) Validate() error {
	// Check for an invalid (non-"OggS") capture pattern.
	if !bytes.Equal(p.CapturePattern[:], oggsBytes) {
		return ErrInvalidCapturePattern
	}

	// Check for an invalid version
	if p.Version != 0 {
		return ErrInvalidVersion
	}

	// Check for invalid header type flags.
	ht := p.HeaderType
	ht &^= Continuation
	ht &^= BOS
	ht &^= EOS
	if ht > 0 {
		return ErrInvalidHeaderType
	}
	return nil
}

var (
	ErrInvalidCapturePattern = errors.New("ogg: invalid capture pattern in page header")
	ErrInvalidVersion        = errors.New("ogg: invalid version in page header")
	ErrInvalidHeaderType     = errors.New("ogg: invalid header type in page header")
	ErrExpectedPage          = errors.New("ogg: expected another page")
)

var (
	// ASCII "OggS" bytes used as the capture pattern of an ogg page header.
	oggsBytes = []byte{79, 103, 103, 83}
)

type Packet struct {
	Pages []*Page
}

// Decoder represents a decoder for a single ogg bitstream.
type Decoder struct {
	r   *bufio.Reader
	eos bool // Whether we've hit the EOS yet.
}

// resync peeks (and then reads one byte) until the magic "OggS" ASCII string
// is found (which signifies the start of an Ogg page header). Because resync
// only peeks to find "OggS", reading data from d.r will still return those
// bytes. This method is strictly to ensure synchronization in the event of
// data loss or co.
func (d *Decoder) resync() error {
	var (
		capturePattern []byte
		err            error
	)
	for {
		// Peek four bytes ahead and see if we find the start of a page.
		capturePattern, err = d.r.Peek(4)
		if err != nil {
			return err
		}

		// See if we got the start of a page or not.
		if bytes.Equal(capturePattern, oggsBytes) {
			// We did! We are synchronized!
			return nil
		}

		// Discard a byte from the stream and continue attempting to synchronize.
		_, err := d.r.ReadByte()
		if err != nil {
			return err
		}
	}
	panic("never here")
}

func (d *Decoder) decodePage() (*Page, error) {
	var (
		err    error
		header PageHeader
	)

	// Resync the stream, if needed.
	err = d.resync()
	if err != nil {
		return nil, err
	}

	// Decode the page header.
	err = binary.Read(d.r, binary.LittleEndian, &header)
	if err != nil {
		return nil, err
	}

	// Validate the page header.
	err = header.Validate()
	if err != nil {
		return nil, err
	}

	log.Println(header.NumberPageSegments)

	page := &Page{
		PageHeader: &header,
	}
	return page, nil
}

// Decode decodes the next packet from the ogg bitstream and returns it. A
// successful decode will always return error == nil, except when the end of
// the stream is reached in which case it may return a non-nil packet (which
// contains data) and error == audio.EOS. It is encouraged to write code like:
//  packet, err := d.Decode()
//  if packet == nil {
//      // packet is only ever nil when err == audio.EOS, so this is the end
//      // of the stream and we can stop decoding.
//  }
//  if err != nil {
//      // Decoding error occured (but not end of stream).
//  }
func (d *Decoder) Decode() (*Packet, error) {
	if d.eos {
		return nil, audio.EOS
	}
	packet := &Packet{}

	for {
		// Decode the next page.
		page, err := d.decodePage()
		if err != nil {
			return nil, ErrExpectedPage
		}
		packet.Pages = append(packet.Pages, page)

		// Check if this is the end of the stream, if so return the packet if
		// it has any data, and return the audio.EOS error.
		if page.HeaderType&EOS > 0 {
			d.eos = true
			if len(packet.Pages) > 0 {
				return packet, audio.EOS
			}
			return nil, audio.EOS
		}

		// Check if the page is not a continuation page, this means the page
		// literally contains a whole packet and we've done all we need to do
		// for this packet.
		if page.HeaderType&Continuation == 0 {
			break
		}
	}

	return packet, nil
}

// NewDecoder creates and initializes a ogg bitstream decoder and returns it.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: bufio.NewReader(r),
	}
}
