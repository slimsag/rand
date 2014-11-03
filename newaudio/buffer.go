// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import (
	"errors"
	"io"
)

// Buffer is a generic audio buffer, it can conceptually be thought of as a
// slice of some audio encoding type.
//
// Conversion between two encoded audio buffers is as simple as:
//  dst, ok := src.(MuLawSamples)
//  if !ok {
//      // Create a new buffer of the target encoding and copy the samples over
//      // because src is not MuLaw encoded.
//      dst = make(MuLawSamples, src.Len())
//      Copy(dst, src)
//  }
type Buffer interface {
	// Len returns the number of elements in the buffer.
	//
	// Equivilent slice syntax:
	//
	//  len(b)
	Len() int

	// Set sets the specified index in the buffer to the specified F64 encoded
	// audio sample, s.
	//
	// If the buffer's audio samples are not stored in F64 encoding, then the
	// sample should be converted to the buffer's internal format and then
	// stored.
	//
	// Just like slices, buffer indices must be non-negative; and no greater
	// than (Len() - 1), or else a panic may occur.
	//
	// Equivilent slice syntax:
	//
	//  b[index] = s
	//   -> b.Set(index, s)
	//
	Set(index int, s F64)

	// At returns the F64 encoded audio sample at the specified index in the
	// buffer.
	//
	// If the buffer's audio samples are not stored in F64 encoding, then the
	// sample should be converted to F64 encoding, and subsequently returned.
	//
	// Just like slices, buffer indices must be non-negative; and no greater
	// than (Len() - 1), or else a panic may occur.
	//
	// Equivilent slice syntax:
	//
	//  b[index]
	//   -> b.At(index)
	//
	At(index int) F64

	// Slice returns a new slice of the buffer, using the low and high
	// parameters.
	//
	// Equivilent slice syntax:
	//
	//  b[low:high]
	//   -> b.Slice(low, high)
	//
	//  b[2:]
	//   -> b.Slice(2, a.Len())
	//
	//  b[:3]
	//   -> b.Slice(0, 3)
	//
	//  b[:]
	//   -> b.Slice(0, a.Len())
	//
	Slice(low, high int) Buffer

	// Make creates and returns a new buffer of this buffers type. This allows
	// allocating a new buffer of exactly the same type for lossless copying of
	// data without knowing about the underlying type.
	//
	// It is exactly the same syntax as the make builtin:
	//
	//  make(MuLawSamples, len, cap)
	//
	// Where cap cannot be less than len.
	Make(length, capacity int) Buffer
}

// BufReader wraps a Buffer to represent a Reader and ReadSeeker.
type BufReader struct {
	original, offset Buffer
}

// Implements Reader interface.
func (r *BufReader) Read(b Buffer) (read int, e error) {
	for read = 0; read < r.offset.Len() && read < b.Len(); read++ {
		b.Set(read, r.offset.At(read))
	}

	// Slice for offset.
	r.offset = r.offset.Slice(read, r.offset.Len())
	if r.offset.Len() == 0 {
		return read, io.EOF
	}
	return read, nil
}

// Implements ReadSeeker interface.
func (r *BufReader) Seek(sample uint64) error {
	if sample > uint64(r.original.Len()) {
		return errors.New("cannot seek past buffer length")
	}
	r.offset = r.original.Slice(int(sample), r.original.Len())
	return nil
}

// NewBufferReader returns a new BufReader wrapping the given buffer.
func NewBufReader(buf Buffer) *BufReader {
	return &BufReader{
		original: buf,
		offset:   buf,
	}
}
