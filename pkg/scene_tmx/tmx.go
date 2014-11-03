// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

// Package tmx implements a Tiled Map XML file loader.
//
// The Tiled Map XML file specification can be found at:
//
//    https://github.com/bjorn/tiled/wiki/TMX-Map-Format
//
// This package supports all of the current file specification with the
// exception of embedded image data (I.e. non-external tileset images).
//
package tmx

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// hexColorToRGBA converts hex color strings to color.RGBA
//
// Alpha value in returned color will always be 255
func hexToRGBA(c string) color.RGBA {
	// There isin't really a color specification I can find on TMX file format,
	// but Tiled exports #RRGGBB hex values, but this also supports #RGB ones
	// just in case some abstract tool uses them by coincidence.

	// Strip leading # if there is one
	if len(c) > 0 && c[0] == '#' {
		c = c[1:]
	}

	// If an invalid length value then simply return
	if len(c) != 6 && len(c) != 3 {
		return color.RGBA{0, 0, 0, 255}
	}

	var r, g, b uint8
	if len(c) == 6 {
		// Parse RRGGBB color
		rgb, err := strconv.ParseUint(c, 16, 48)
		if err != nil {
			return color.RGBA{0, 0, 0, 255}
		}
		r = uint8(rgb >> 16)
		g = uint8(rgb >> 8)
		b = uint8(rgb)
	} else {
		// Parse #RGB values
		rgb, err := strconv.ParseUint(c, 16, 24)
		if err != nil {
			return color.RGBA{0, 0, 0, 255}
		}
		r = uint8(rgb>>8) & 0xf
		g = uint8(rgb>>4) & 0xf
		b = uint8(rgb) & 0xf
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}
	return color.RGBA{r, g, b, 255}
}

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type xmlProperties struct {
	Property []xmlProperty `xml:"property"`
}

func (p xmlProperties) toMap() map[string]string {
	m := make(map[string]string, len(p.Property))
	for _, p := range p.Property {
		m[p.Name] = p.Value
	}
	return m
}

type xmlTileoffset struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}

type xmlTerrain struct {
	Name string `xml:"name,attr"`
	Tile int    `xml:"id,attr"`
}

type xmlTerraintypes struct {
	Terrain []xmlTerrain `xml:"terrain"`
}

type xmlMap struct {
	Version         string        `xml:"version,attr"`
	Orientation     string        `xml:"orientation,attr"`
	Width           int           `xml:"width,attr"`
	Height          int           `xml:"height,attr"`
	TileWidth       int           `xml:"tilewidth,attr"`
	TileHeight      int           `xml:"tileheight,attr"`
	BackgroundColor string        `xml:"backgroundcolor,attr"`
	Properties      xmlProperties `xml:"properties"`
	Tileset         []xmlTileset  `xml:"tileset"`
	Layer           []xmlLayer    `xml:"layer"`
}

// Load loads the TMX map file data and returns a loaded *Map.
//
// nil and a error will be returned if there are any problems loading the data.
func Load(data []byte) (*Map, error) {
	// Unmarshal map data
	x := new(xmlMap)
	err := xml.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}

	// Parse version string
	split := strings.Split(x.Version, ".")
	var major, minor int
	if len(split) == 2 {
		major, err = strconv.Atoi(split[0])
		if err != nil {
			return nil, err
		}

		minor, err = strconv.Atoi(split[1])
		if err != nil {
			return nil, err
		}
	}

	// Find map orientation
	var orient Orientation
	switch x.Orientation {
	case "orthogonal":
		orient = Orthogonal
	case "isometric":
		orient = Isometric
	case "staggered":
		orient = Staggered
	default:
		return nil, fmt.Errorf("unknown map orientation.")
	}

	// Find map properties
	props := make(map[string]string, len(x.Properties.Property))
	for _, prop := range x.Properties.Property {
		props[prop.Name] = prop.Value
	}

	// Convert the tilesets
	tilesets := make([]*Tileset, len(x.Tileset))
	for i, tsx := range x.Tileset {
		ts := &Tileset{
			Name:     tsx.Name,
			Firstgid: tsx.Firstgid,
			Source:   tsx.Source,
			Width:    tsx.TileWidth,
			Height:   tsx.TileHeight,
			Spacing:  tsx.Spacing,
			Margin:   tsx.Margin,
		}

		// Find tileset offset
		ts.OffsetX, ts.OffsetY = tsx.Tileoffset.X, tsx.Tileoffset.Y

		// Find tileset properties
		ts.Properties = tsx.Properties.toMap()

		// Find image properties
		ts.Image = &Image{
			Source: tsx.Image.Source,
			Width:  tsx.Image.Width,
			Height: tsx.Image.Height,
		}

		// Find tile definitions
		ts.Tiles = tsx.tilesMap()

		// Find terrain definitions
		ts.Terrain = tsx.terrainTypes()

		tilesets[i] = ts
	}

	// Manage loading layers
	layers := make([]*Layer, len(x.Layer))
	for i, xl := range x.Layer {
		var err error
		layers[i], err = xl.toLayer(x.Width, x.Height)
		if err != nil {
			return nil, err
		}
	}

	// Create actual map
	m := &Map{
		VersionMajor:    major,
		VersionMinor:    minor,
		Orientation:     orient,
		Width:           x.Width,
		Height:          x.Height,
		TileWidth:       x.TileWidth,
		TileHeight:      x.TileHeight,
		BackgroundColor: hexToRGBA(x.BackgroundColor),
		Properties:      props,
		Tilesets:        tilesets,
		Layers:          layers,
	}

	return m, nil
}
