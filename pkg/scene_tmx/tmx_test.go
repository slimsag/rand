// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func verify(t *testing.T, name string) {
	// Open file
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}

	// Read file data
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	// Load the map
	m, err := Load(data)
	if err != nil {
		t.Fatal(err)
	}

	// External tilesets in the map must be loaded seperately
	for _, ts := range m.Tilesets {
		if len(ts.Source) > 0 {
			// Open tsx file
			f, err := os.Open(filepath.Join("testdata", ts.Source))
			if err != nil {
				t.Fatal(err)
			}

			// Read file data
			data, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			// Load the tileset
			err = ts.Load(data)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	if m.VersionMajor != 1 || m.VersionMinor != 0 {
		t.Log(m.VersionMajor, m.VersionMinor)
		t.Fatal("incorrect version number")
	}
	if m.Orientation != Orthogonal {
		t.Log(m.Orientation)
		t.Fatal("incorrect orientation")
	}
	if m.Width != 60 || m.Height != 10 {
		t.Log(m.Width)
		t.Log(m.Height)
		t.Fatal("incorrect width/height")
	}
	if m.TileWidth != 32 || m.TileHeight != 32 {
		t.Log(m.TileWidth)
		t.Log(m.TileHeight)
		t.Fatal("incorrect tile width/height")
	}
	c := m.BackgroundColor
	if c.R != 255 || c.G != 0 || c.B != 0 || c.A != 255 {
		t.Log(m.BackgroundColor)
		t.Fatal("incorrect background color")
	}

	v := m.Properties["mymap_prop"]
	if v != "mymap_prop_value" {
		t.Log(m.Properties)
		t.Fatal("incorrect map property value")
	}

	for _, layer := range m.Layers {
		for x := 0; x < m.Width; x++ {
			for y := 0; y < m.Height; y++ {
				gid, hasTile := layer.Tiles[Coord{x, y}]
				if hasTile {
					tileset := m.FindTileset(gid)
					if tileset == nil {
						t.Fatal("FindTileset failed to find correct tileset.")
					}

					_ = m.TilesetRect(tileset, tileset.Image.Width, tileset.Image.Height, true, gid)
				}
			}
		}
	}
}

func TestXMLMap(t *testing.T) {
	verify(t, "test_xml.tmx")
}

func TestXMLDTDMap(t *testing.T) {
	verify(t, "test_xml_dtd.tmx")
}

func TestCSVMap(t *testing.T) {
	verify(t, "test_csv.tmx")
}

func TestCSVTSXMap(t *testing.T) {
	verify(t, "test_csv_tsx.tmx")
}

func TestBase64Map(t *testing.T) {
	verify(t, "test_base64.tmx")
}

func TestBase64GZipMap(t *testing.T) {
	verify(t, "test_base64_gzip.tmx")
}

func TestBase64ZLibMap(t *testing.T) {
	verify(t, "test_base64_zlib.tmx")
}
