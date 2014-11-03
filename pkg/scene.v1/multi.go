package scene

import (
	"image"

	"azul3d.org/gfx.v1"
)

// Multi implements the Scene interface by performing every scene action on
// every scene in the slice.
type Multi []Scene

// Implements the Scene interface by simply calling Add(d) on every scene in
// the slice.
func (m Multi) Add(d gfx.Drawable) {
	for _, s := range m {
		s.Add(d)
	}
}

// Implements the Scene interface by simply calling Has(d) on every scene in
// the slice, and returning the first (true) result.
func (m Multi) Has(d gfx.Drawable) bool {
	for _, s := range m {
		if s.Has(d) {
			return true
		}
	}
	return false
}

// Implements the Scene interface by simply calling Remove(d) on every scene in
// the slice.
func (m Multi) Remove(d gfx.Drawable) {
	for _, s := range m {
		s.Remove(d)
	}
}

// Implements the Scene interface by simply calling DrawTo(c, bounds, cam) on
// every scene in the slice.
func (m Multi) DrawTo(c gfx.Canvas, bounds image.Rectangle, cam *gfx.Camera) {
	for _, s := range m {
		s.DrawTo(c, bounds, cam)
	}
}
