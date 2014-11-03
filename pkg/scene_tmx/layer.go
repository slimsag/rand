// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

import (
	"fmt"
)

type xmlLayer struct {
	Name    string  `xml:"name,attr"`
	Opacity float64 `xml:"opacity,attr"`
	Visible int     `xml:"visible,attr"`
	Data    xmlData `xml:"data"`
}

func (x xmlLayer) toLayer(width, height int) (*Layer, error) {
	tiles, err := x.Data.tiles(width, height)
	if err != nil {
		return nil, err
	}
	return &Layer{
		Name:    x.Name,
		Opacity: x.Opacity,
		Visible: x.Visible != 0,
		Tiles:   tiles,
	}, nil
}

// Coord represents a single 2D coordinate pair (x, y)
type Coord struct {
	X, Y int
}

// Layer represents a single map layer and all of it's tiles
type Layer struct {
	// The name of the layer.
	Name string

	// Value between 0 and 1 representing the opacity of the layer.
	Opacity float64

	// Boolean value representing whether or not the layer is visible.
	Visible bool

	// A map of 2D coordinates in this layer to so called "global tile IDs"
	// (gids).
	//
	// 2D coordinates whose gid's are zero (I.e. 'no tile') are not stored in
	// the map for efficiency reasons (as a good majority are zero).
	//
	// gids are global, since they may refere to a tile from any of the
	// tilesets used by the map. In order to find out from which tileset the
	// tile is you need to find the tileset with the highest Firstgid that is
	// still lower or equal than the gid. The tilesets are always stored with
	// increasing firstgids.
	Tiles map[Coord]uint32
}

// String returns a string representation of this layer.
func (l *Layer) String() string {
	return fmt.Sprintf("Layer(Name=%q, Opacity=%1.f, Visible=%v)", l.Name, l.Opacity, l.Visible)
}
