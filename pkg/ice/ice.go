// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ice

import (
	"azul3d.org/v1/gfx"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	// Error returned by LoadFile and SaveFile for invalid file extensions.
	InvalidExt = errors.New("invalid file extension")
)

// Scene represents a single graphics scene.
type Scene struct {
	// A map of properties for the scene by name.
	Props map[string]interface{}

	// A map of cameras by name.
	Cameras map[string]*gfx.Camera

	// A map of objects by name.
	Objects map[string]*gfx.Object
}

// Save saves the scene to an Ice binary file to the given writer and returns
// an error, if any.
//
// See also the SaveFile method which operates on a fixed file-path easilly.
func (s *Scene) Save(w io.Writer) error {
	w = gzip.NewWriter(w)
	enc := gob.NewEncoder(w)
	err := enc.Encode(s)
	return err
}

// SaveJSON saves the scene to an Ice JSON file to the given writer and returns
// an error, if any.
//
// See also the SaveJSONFile method which operates on a fixed file-path
// easilly.
func (s *Scene) SaveJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(s)
}

// SaveFile saves the scene to an Ice binary or JSON file (determined by the
// file extension) and returns an error, if any.
//
// See also the Save and SaveJSON method which operate on a io.Writer instead
// of a fixed file-path.
func (s *Scene) SaveFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	ext := filepath.Ext(path)
	if ext == ".ice" {
		return s.Save(f)
	} else if ext == ".json" {
		return s.SaveJSON(f)
	}
	return InvalidExt
}
