// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
)

const (
	// Flags representing horizontal, vertical, and diagonal tile flipping.
	// These can be used in combination with gid's, like so:
	//
	// if (gid & FLIPPED_HORIZONTALLY_FLAG) > 0 {
	//     ...draw the tiled flipped horizontally...
	// }
	FLIPPED_HORIZONTALLY_FLAG uint32 = 0x80000000
	FLIPPED_VERTICALLY_FLAG   uint32 = 0x40000000
	FLIPPED_DIAGONALLY_FLAG   uint32 = 0x20000000
)

type xmlTile struct {
	ID          int           `xml:"id,attr"`
	Terrain     []byte        `xml:"terrain,attr"`
	Probability float64       `xml:"probability,attr"`
	Properties  xmlProperties `xml:"properties"`
	Image       xmlImage      `xml:"image"`
}

func (x xmlTile) terrainArray() (indices [4]int) {
	var (
		buf = bytes.NewBuffer(x.Terrain)
		k   = 0
	)
	indices = [4]int{-1, -1, -1, -1}
	terrain, err := csv.NewReader(buf).ReadAll()
	if err != nil {
		return
	}
	for _, i := range terrain {
		for _, j := range i {
			index, _ := strconv.Atoi(j)
			indices[k] = index
			k++
			if k >= 4 {
				return
			}
		}
	}
	return
}

func (x xmlTile) toTile() *Tile {
	return &Tile{
		ID:          x.ID,
		Terrain:     x.terrainArray(),
		Probability: x.Probability,
		Properties:  x.Properties.toMap(),
		Image:       x.Image.toImage(),
	}
}

// Tile represents a single tile definition and it's properties
type Tile struct {
	// The ID of the tile
	ID int

	// An array defining the terrain type of each corner of the tile, as indices
	// into the terrain types slice of the tileset this tile came from, in the
	// order of: top left, top right, bottom left, bottom right.
	//
	// -1 values have a meaning of 'no terrain'.
	Terrain [4]int

	// Percentage chance indicating the probability that this tile is chosen
	// when editing with the terrain tool.
	Probability float64

	// Map of properties for the tile
	Properties map[string]string

	// Image for the tile
	Image *Image
}

// String returns a string representation of this tileset.
func (t *Tile) String() string {
	return fmt.Sprintf("Tile(ID=%v)", t.ID)
}
