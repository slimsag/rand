// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ice

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Resolver is a file resolver. Because it is an interface you may define your
// own resolver that, for example, loads resources over a network.
//
// This package exposes a file system resolver by default, as FileSystem, which
// works for most cases.
type Resolver interface {
	// Resolve should resolve the given filepath and return an io.ReadCloser
	// usable for reading the named file.
	//
	// If any error occurs while trying to resolve the filepath, it may be
	// returned only if the returned io.ReadCloser is nil.
	Resolve(filepath string) (io.ReadCloser, error)
}

// FileSystem is the default file system resolver. Because it is often harmful
// to use absolute paths as Ice resources, this resolver will explicitly
// disallow any resource whose filepath is not relative (by simply returning
// the AbsPathError).
var FileSystem Resolver

type fsResolver struct{}

// AbsPathError is the error returned by the default FileSystem resolver if any
// resource's filepath is absolute and not relative. This is because absolute
// filepaths are generally harmful as Ice resources.
var AbsPathError = errors.New("resource is using an absolute path")

// Implements the Resolver interface.
func (r *fsResolver) Resolve(path string) (io.ReadCloser, error) {
	if filepath.IsAbs(path) {
		return nil, AbsPathError
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func init() {
	FileSystem = &fsResolver{}
}
