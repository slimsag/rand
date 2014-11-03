// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

import (
	"encoding/xml"
	"fmt"
)

type xmlTileset struct {
	Raw          []byte `xml:",innerxml"`
	Firstgid     uint32 `xml:"firstgid,attr"`
	Source       string `xml:"source,attr"`
	Name         string `xml:"name,attr"`
	TileWidth    int    `xml:"tilewidth,attr"`
	TileHeight   int    `xml:"tileheight,attr"`
	Spacing      int    `xml:"spacing,attr"`
	Margin       int    `xml:"margin,attr"`
	Tileoffset   xmlTileoffset
	Properties   xmlProperties   `xml:"properties"`
	Image        xmlImage        `xml:"image"`
	Tile         []xmlTile       `xml:"tile"`
	Terraintypes xmlTerraintypes `xml:"terraintypes"`
}

func (x *xmlTileset) tilesMap() map[int]*Tile {
	tiles := make(map[int]*Tile, len(x.Tile))
	for _, xt := range x.Tile {
		tiles[xt.ID] = xt.toTile()
	}
	return tiles
}

func (x *xmlTileset) terrainTypes() []TerrainType {
	terrainTypes := make([]TerrainType, len(x.Terraintypes.Terrain))
	for i, xt := range x.Terraintypes.Terrain {
		terrainTypes[i] = TerrainType{
			Name: xt.Name,
			Tile: xt.Tile,
		}
	}
	return terrainTypes
}

// TerrainType defines a single terrain with a name and associated tile ID
type TerrainType struct {
	// Name of the terrain type
	Name string

	// Tile ID
	Tile int
}

// Tileset represents a tileset of a map, as loaded from the TMX file or from
// a external TSX file.
type Tileset struct {
	// The name of this tileset.
	Name string

	// The first global tile ID of this tileset (this global ID maps to the
	// first tile in this tileset).
	Firstgid uint32

	// The tilset source (tsx) file, if this tileset was loaded externally from
	// the TMX map.
	Source string

	// The maximum width/height of tiles in this tileset in pixels.
	Width, Height int

	// The horizontal and vertical offset of tiles in this tileset in pixels,
	// where +Y is down.
	OffsetX, OffsetY int

	// The spacing in pixels between the tiles in this tileset.
	Spacing int

	// The margin in pixels around the tiles in this tileset.
	Margin int

	// Map of property names and values for all properties set on the map.
	Properties map[string]string

	// The image of the tileset
	Image *Image

	// Tiles represents a map of tile ID's and their associated definitions.
	Tiles map[int]*Tile

	// The slice of terrain types
	Terrain []TerrainType
}

// String returns a string representation of this tileset.
func (t *Tileset) String() string {
	return fmt.Sprintf("Tileset(Name=%q, Firstgid=%v, Source=%q, Size=%dx%dpx, Offset=%dx%dpx, Spacing=%dpx, Margin=%dpx)", t.Name, t.Firstgid, t.Source, t.Width, t.Height, t.OffsetX, t.OffsetY, t.Spacing, t.Margin)
}

// Load loads the specified data as this tileset or returns a error if the data
// is invalid.
//
// If len(m.Source) == nil (I.e. if this tileset is not an external tsx file)
// then a panic will occur.
//
// Clients should ensure properly synchronized read/write access to the tileset
// structure as this function write's to it's memory and does not attempt any
// synchronization with other goroutines who are reading from it (data race).
func (t *Tileset) Load(data []byte) error {
	x := new(xmlTileset)
	err := xml.Unmarshal(data, &x)
	if err != nil {
		return err
	}
	t.Name = x.Name
	t.Width = x.TileWidth
	t.Height = x.TileHeight
	t.Spacing = x.Spacing
	t.Margin = x.Margin

	// Find tileset offset
	t.OffsetX, t.OffsetY = x.Tileoffset.X, x.Tileoffset.Y

	// Find tileset properties
	t.Properties = x.Properties.toMap()

	// Find image properties
	t.Image = &Image{
		Source: x.Image.Source,
		Width:  x.Image.Width,
		Height: x.Image.Height,
	}

	// Find tile definitions
	t.Tiles = x.tilesMap()

	// Find terrain definitions
	t.Terrain = x.terrainTypes()

	return nil
}
