// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import "io"

// Reader is a generic interface which describes any type who can have audio
// samples read from it into an audio buffer.
type Reader interface {
	// Read tries to read into the audio buffer, b, filling it with at max
	// b.Len() audio samples.
	//
	// Returned is the number of samples that where read into the buffer, and
	// an error if any occured.
	//
	// It is possible for the number of samples read to be non-zero; and for an
	// error to be returned at the same time (E.g. read 300 audio samples, but
	// also encountered io.EOF).
	Read(b Buffer) (read int, e error)
}

// ReadSeeker is the generic seekable audio reader interface.
type ReadSeeker interface {
	Reader

	// Seek seeks to the specified sample number, relative to the start of the
	// stream. As such, subsequent Read() calls on the Reader, begin reading at
	// the specified sample.
	//
	// If any error is returned, it means it was impossible to seek to the
	// specified audio sample for some reason, and that the current playhead is
	// unchanged.
	Seek(sample uint64) error
}

// Writer is a generic interface which describes any type who can have audio
// samples written from an audio buffer into it.
type Writer interface {
	// Write attempts to write all, b.Len(), samples in the buffer to the
	// writer.
	//
	// Returned is the number of samples from the buffer that where wrote to
	// the writer, and an error if any occured.
	//
	// The number of samples wrote may be less than buf.Len(), in which case
	// you should subsequently write b.Slice(wrote, b.Len()) until you have
	// finished sending all data or an error occurs.
	//
	// If any error is returned, it should be considered as fatal to the
	// writer, no more data can subsequently be wrote to the writer.
	Write(b Buffer) (wrote int, err error)
}

// WriterTo is the interface that wraps the WriteTo method.
//
// WriteTo writes data to w until there's no more data to write or when an
// error occurs. The return value n is the number of samples written. Any error
// encountered during the write is also returned.
//
// The Copy function uses WriterTo if available.
type WriterTo interface {
	WriteTo(w Writer) (n int64, err error)
}

// ReaderFrom is the interface that wraps the ReadFrom method.
//
// ReadFrom reads data from r until io.EOF or error. The return value n is the
// number of bytes read. Any error except io.EOF encountered during the read is
// also returned.
//
// The Copy function uses ReaderFrom if available.
type ReaderFrom interface {
	ReadFrom(r Reader) (n int64, err error)
}

// Copy copies from src to dst until either io.EOF is reached on src or an
// error occurs.  It returns the number of samples copied and the first error
// encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == io.EOF. Because Copy is
// defined to read from src until EOF, it does not treat an io.EOF from Read as
// an error to be reported.
//
// If src implements the WriterTo interface, the copy is implemented by calling
// src.WriteTo(dst). Otherwise, if dst implements the ReaderFrom interface, the
// copy is implemented by calling dst.ReadFrom(src).
func Copy(dst Writer, src Reader) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy. Avoids an
	// allocation and a copy.
	if wt, ok := src.(WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	buf := make(F64Samples, (32*1024)/8)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
