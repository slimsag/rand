// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
)

type xmlDataTile struct {
	Gid uint32 `xml:"gid,attr"`
}

type xmlData struct {
	Data []byte `xml:",innerxml"`

	// base64, csv
	Encoding string `xml:"encoding,attr"`

	// gzip, zlib
	Compression string `xml:"compression,attr"`

	Tile []xmlDataTile `xml:"tile"`
}

var (
	// Error representing an unknown encoding method inside of a tmx file
	ErrBadEncoding = errors.New("tile data encoding type is not supported")

	// Error representing an unknown compression method inside of a tmx file
	ErrBadCompression = errors.New("tile data compression type is not supported")
)

func toCoord(index, width, height int) Coord {
	if width == 0 {
		panic("width == 0")
	}
	return Coord{
		X: index % width,
		Y: index / width,
	}
}

func (x xmlData) tiles(width, height int) (map[Coord]uint32, error) {
	tiles := make(map[Coord]uint32)
	switch x.Encoding {
	case "":
		// No encoding, plain XML elements
		for coordIndex, xt := range x.Tile {
			if xt.Gid != 0 {
				tiles[toCoord(coordIndex, width, height)] = xt.Gid
			}
		}

	case "csv":
		buf := bytes.NewBuffer(x.Data)
		r := csv.NewReader(buf)
		r.FieldsPerRecord = -1
		csvLines, err := r.ReadAll()
		if err != nil {
			return nil, err
		}
		coordIndex := 0
		for _, line := range csvLines {
			for _, csvString := range line {
				if len(csvString) > 0 {
					gid, err := strconv.ParseUint(csvString, 10, 0)
					if err != nil {
						return nil, err
					}
					if gid != 0 {
						tiles[toCoord(coordIndex, width, height)] = uint32(gid)
					}
					coordIndex++
				}
			}
		}

	case "base64":
		data := bytes.Replace(x.Data, []byte(" "), []byte(""), -1)
		data = bytes.Replace(data, []byte("\r"), []byte(""), -1)
		data = bytes.Replace(data, []byte("\n"), []byte(""), -1)
		buf := bytes.NewBuffer(data)
		decoded := base64.NewDecoder(base64.StdEncoding, buf)

		var decompressed io.Reader
		switch x.Compression {
		case "":
			// No compression
			decompressed = decoded
		case "zlib":
			r, err := zlib.NewReader(decoded)
			if err != nil {
				return nil, err
			}
			defer r.Close()
			decompressed = r

		case "gzip":
			r, err := gzip.NewReader(decoded)
			if err != nil {
				return nil, err
			}
			defer r.Close()
			decompressed = r

		default:
			return nil, ErrBadCompression
		}
		coordIndex := 0
		for {
			var gid uint32
			err := binary.Read(decompressed, binary.LittleEndian, &gid)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			if gid != 0 {
				tiles[toCoord(coordIndex, width, height)] = gid
			}
			coordIndex++
		}

	default:
		return nil, ErrBadEncoding
	}
	return tiles, nil
}
