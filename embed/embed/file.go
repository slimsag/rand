// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

// Package embed allows for creating an embedded file within a Go binary.
package embed

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var (
	InvalidFileErr = errors.New("Invalid or corrupt embedded file.")
	tag            = []byte("-embedded-")
)

func Offset(bin *os.File) (int64, error) {
	var n int

	fi, err := bin.Stat()
	if err != nil {
		return 0, err
	}
	buf := make([]byte, len(tag))
	n, err = bin.ReadAt(buf, fi.Size()-int64(len(buf))-8)
	if n < len(buf) {
		return 0, InvalidFileErr
	}

	for i, tagByte := range tag {
		if buf[i] != tagByte {
			return 0, InvalidFileErr
		}
	}

	buf = make([]byte, 8)
	n, err = bin.ReadAt(buf, fi.Size()-int64(len(buf)))
	if n < len(buf) {
		return 0, InvalidFileErr
	}

	dataSize := int64(binary.LittleEndian.Uint64(buf))
	offset := fi.Size() - dataSize - int64(len(tag)) - 8
	return offset, nil
}

// Preamble returns an limited reader around the preamble section of an embedded file
func Preamble(bin *os.File) (io.Reader, error) {
	offset, err := Offset(bin)
	if err != nil {
		return nil, err
	}
	return io.LimitReader(bin, offset), nil
}

func WriteFooter(bin *os.File, dataSize int64) error {
	buf := make([]byte, len(tag))
	for i, c := range tag {
		buf[i] = c
	}
	n, err := bin.Write(buf)
	if n < len(buf) && err != nil {
		return err
	}

	buf = make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(dataSize))
	n, err = bin.Write(buf)
	if n < len(buf) && err != nil {
		return err
	}
	return nil
}
