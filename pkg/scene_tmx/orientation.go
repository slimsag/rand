// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

// Orientation represents the map's orientation. It will always be one of the
// predefined Orientation constants and will never be Invalid.
type Orientation int

const (
	// Invalid orientation for catching zero-value related issues
	Invalid Orientation = iota

	// Orthogonal map orientation
	Orthogonal

	// Isometric map orientation
	Isometric

	// Staggered map orientation
	Staggered
)
