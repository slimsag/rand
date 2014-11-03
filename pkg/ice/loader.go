// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ice

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Load loads and returns a scene from a binary Ice file and returns it. If any
// error occurs it and a nil scene is returned immedietly.
//
// The resolver parameter allows for specifying a file path resolver used when
// loading textures, shaders, etc. This allows for loading of all resources
// over, for example, a network. The default file system resolver is exposed
// by this package as FileSystem.
func Load(r io.Reader, resolver Resolver) (*Scene, error) {
	var err error

	// Use gzip to decompress the file.
	r, err = gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	// Use gob to decode the scene data.
	dec := gob.NewDecoder(r)
	s := new(Scene)
	err = dec.Decode(&s)
	if err != nil {
		return nil, err
	}
	// FIXME: invoke scene()
	return s, nil
}

// LoadJSON loads and returns a scene from a JSON Ice file and returns it. If
// any error occurs it and a nil scene is returned immedietly.
//
// The resolver parameter allows for specifying a file path resolver used when
// loading textures, shaders, etc. This allows for loading of all resources
// over, for example, a network. The default file system resolver is exposed
// by this package as FileSystem.
func LoadJSON(r io.Reader, resolver Resolver) (*Scene, error) {
	var err error

	// Use JSON to decode the scene data.
	dec := json.NewDecoder(r)
	s := new(jsonScene)
	err = dec.Decode(&s)
	if err != nil {
		return nil, err
	}
	return s.scene(resolver), nil
}

// LoadFile examines the extension of the filepath, if the extension is .ice
// then this method invokes Load, if it is .json then this method invokes
// LoadJSON.
//
// This function uses the default FileSystem resolver to load textures,
// shaders, etc.
//
// Any errors that occur are returned along with a nil scene.
func LoadFile(path string) (*Scene, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var scene *Scene
	ext := filepath.Ext(path)
	if ext == ".ice" {
		scene, err = Load(f, FileSystem)
	} else if ext == ".json" {
		scene, err = LoadJSON(f, FileSystem)
	} else {
		return nil, InvalidExt
	}
	if err != nil {
		return nil, err
	}
	return scene, nil
}
