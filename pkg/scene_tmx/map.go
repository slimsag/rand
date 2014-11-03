// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

import (
	"fmt"
	"image"
	"image/color"
)

// Map represents a single TMX map file.
//
// Although TileWidth and TileHeight describe the general size of tiles in
// pixels, individual tiles may have different sizes. Larger tiles will extend
// at the top and right (E.g. they are anchored to the bottom left).
type Map struct {
	// Version of the map.
	//
	// E.g. VersionMajor=1, VersionMinor=0 for "1.0"
	VersionMajor, VersionMinor int

	// Orientation of the map.
	//
	// Like "orthogonal", "isometric" or "staggered".
	Orientation Orientation

	// Width and height of the map in tiles.
	Width, Height int

	// Width and height of a tile in pixels.
	TileWidth, TileHeight int

	// Background color of the map.
	//
	// Like "#FF0000".
	BackgroundColor color.RGBA

	// Map of property names and values for all properties set on the map.
	Properties map[string]string

	// A list of all loaded tilesets of this map.
	Tilesets []*Tileset

	// A list of all the layers of this map.
	Layers []*Layer
}

// String returns a string representation of this map.
func (m *Map) String() string {
	return fmt.Sprintf("Map(Version=%d.%d, Size=%dx%d, TileSize=%dx%dpx)", m.VersionMajor, m.VersionMinor, m.Width, m.Height, m.TileWidth, m.TileHeight)
}

// FindTileset returns the proper tileset for the given global tile id.
//
// If the global tile id is invalid this function will return nil.
func (m *Map) FindTileset(gid uint32) *Tileset {
	gid &^= (FLIPPED_HORIZONTALLY_FLAG | FLIPPED_VERTICALLY_FLAG | FLIPPED_DIAGONALLY_FLAG)

	for i := len(m.Tilesets) - 1; i >= 0; i-- {
		ts := m.Tilesets[i]
		if ts.Firstgid <= gid {
			return ts
		}
	}
	return nil
}

// TilesetTile returns the proper tile definition for the given global tile id.
//
// If there is no tile definition for the given gid (can be common), or if the
// global tile id is invalid this function will return nil.
func (m *Map) TilesetTile(ts *Tileset, gid uint32) *Tile {
	gid &^= (FLIPPED_HORIZONTALLY_FLAG | FLIPPED_VERTICALLY_FLAG | FLIPPED_DIAGONALLY_FLAG)
	id := int(gid - ts.Firstgid)
	return ts.Tiles[id]
}

// TilesetRect returns a image rectangle describing what part of the tileset
// image represents the tile for the given gid.
//
// The image width and height must be passed as parameters because
// ts.Image.Width and ts.Image.Height are not always available.
//
// If spacingAndMargins is true, then spacing and margins are applied to the
// rectangle.
func (m *Map) TilesetRect(ts *Tileset, width, height int, spacingAndMargins bool, gid uint32) image.Rectangle {
	gid &^= (FLIPPED_HORIZONTALLY_FLAG | FLIPPED_VERTICALLY_FLAG | FLIPPED_DIAGONALLY_FLAG)
	id := int(gid - ts.Firstgid)
	coord := toCoord(id, width/ts.Width, height/ts.Height)
	var cx, cy int
	if spacingAndMargins {
		cx = coord.X * (ts.Width + ts.Spacing)
		cx += ts.Spacing
		cy = coord.Y * (ts.Height + ts.Spacing)
		cy += ts.Margin
	} else {
		cx = coord.X * ts.Width
		cy = coord.Y * ts.Height
	}
	return image.Rect(cx, cy, cx+ts.Width, cy+ts.Height)
}
