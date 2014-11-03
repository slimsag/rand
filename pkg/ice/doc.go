// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ice implements a reader and writer for the Azul3D scene format.
//
// Ice is a binary 3D scene storage format designed specifically for Azul3D. It
// offers gzip compression and data is gob-encoded. Ice files have a large
// amount of versatility within the engine because it was specifically designed
// for it.
//
// It is common place for users to directly ship ice files with applications
// made in Azul3D instead of other model formats because it offers great
// compression and utility.
//
// Modeling Software
//
// Ice is strictly for use within Go applications since the binary file format
// in use is gob (see encoding/gob). Therefore Ice binary files are not a very
// good target for modeling software exporters.
//
// For this very reason Ice also has a JSON file format, which should be used
// for compatability with other non-Go software. JSON allows for easy
// interchange of model formats from modeling software to Azul3D applications.
//
// Conversion tools from Ice's binary and JSON file formats are available such
// that converting from an Ice JSON file to an Ice binary file (or vice-versa)
// is made easy.
//
// Converters for existing file formats to the Ice JSON file format will surely
// arise within the community as the project grows.
//
// Loading
//
// Loading an Ice binary or JSON file is easy:
//  scene, err := ice.LoadFile("path/to/file.ice")
//  handle(err)
//
// Or
//
//  scene, err := ice.LoadFile("path/to/file.json")
//  handle(err)
//
// Saving
//
// Saving an Ice binary file is easy:
//  err := scene.SaveFile("path/to/file.ice")
//  handle(err)
//
// And saving an Ice JSON file is also easy:
//  err := scene.SaveFile("path/to/file.json")
//  handle(err)
package ice
