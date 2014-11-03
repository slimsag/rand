// Copyright 2012 Lightpoke. All rights reserved.
// This source code is subject to the terms and
// conditions defined in the "License.txt" file.

package tmx

import (
	"fmt"
	"image/color"
)

type xmlImage struct {
	Format string `xml:"format,attr"`
	// TODO: expose image <data> element for embedded image data ?
	Source string `xml:"source,attr"`
	Trans  string `xml:"trans,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

func (x xmlImage) toImage() *Image {
	return &Image{
		Format: x.Format,
		Source: x.Source,
		Trans:  hexToRGBA(x.Trans),
		Width:  x.Width,
		Height: x.Height,
	}
}

// Image represents the source and properties of a image
type Image struct {
	// Format of the embedded image data (if any).
	Format string

	// The file path at which the image may be found
	Source string

	// The color in the image representing transparency (if any).
	//
	// The alpha (A) component of the color will always be 255.
	Trans color.RGBA

	// The width and height of the image (useful mostly only for correction
	// when the image's size changes from that known to the TMX file).
	Width, Height int
}

// String returns a string representation of this image.
func (i *Image) String() string {
	return fmt.Sprintf("Image(Source=%q, Size=%dx%dpx)", i.Source, i.Width, i.Height)
}
