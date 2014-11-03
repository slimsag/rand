package scene

import (
	"image"
	"sync"

	"azul3d.org/gfx.v1"
)

// Scene represents an arbitrary list of drawable objects which effectively
// compose a graphical scene.
//
// All methods on a Scene are safe for access from multiple goroutines
// concurrently.
type Scene interface {
	// Add adds the drawable object to this scene. If the object already
	// existed in the scene, false is returned.
	Add(d gfx.Drawable) bool

	// Has tells whether or not the drawable object is within this scene or
	// not.
	Has(d gfx.Drawable) bool

	// Remove removes the drawable object from this scene. If the object does
	// not exist in the scene then false is returned.
	Remove(d gfx.Drawable) bool

	// Iter iterates over all of the drawable objects in this scene. The order
	// in which the objects are iterated is dependent on the implementation of
	// the scene.
	//
	// If the callback returns false, the iteration is stopped.
	Iter(callback func(d gfx.Drawable) (stop bool))

	// DrawTo draws all of the objects in this scene to the given canvas, using
	// the given camera and bounding box. It is effectively the same as:
	//  for drawable := range scene {
	//  	drawable.DrawTo(c, bounds, cam)
	//  }
	//
	// Drawables may choose to ignore the bounds or camera parameters and
	// instead substitute their own for the draw (this allows individual
	// drawable objects to explicitly use their own camera or draw area).
	DrawTo(c gfx.Canvas, bounds image.Rectangle, cam *gfx.Camera)
}

// Basic implements the Scene interface using a simple slice.
type Basic struct {
	sync.RWMutex
	S []gfx.Drawable
}

// Implements the Scene interface.
func (s Basic) Add(d gfx.Drawable) bool {
	s.Lock()

	// Is it already in the scene? If so don't add it.
	for _, drawable := range s.S {
		if d == drawable {
			s.Unlock()
			return false
		}
	}

	// Append it to the slice.
	s.S = append(s.S, d)
	s.Unlock()
	return true
}

// Implements the Scene interface.
func (s Basic) Has(d gfx.Drawable) bool {
	s.RLock()
	for _, drawable := range s.S {
		if d == drawable {
			s.RUnlock()
			return true
		}
	}
	s.RUnlock()
	return false
}

// Implements the Scene interface.
func (s Basic) Remove(d gfx.Drawable) bool {
	s.Lock()
	for i, drawable := range s.S {
		if d == drawable {
			// Delete it from the slice.
			s.S = append(s.S[:i], s.S[i+1:]...)
			s.Unlock()
			return true
		}
	}
	s.Unlock()
	return false
}

// Implements the Scene interface.
func (s Basic) Iter(callback func(d gfx.Drawable) (stop bool)) {
	s.RLock()
	for _, drawable := range s.S {
		s.RUnlock() // Unlock during callback.
		ret := callback(drawable)
		s.RLock() // Lock again.

		if !ret {
			break
		}
	}
	s.RUnlock()
}

// Implements the Scene interface.
func (s Basic) DrawTo(c gfx.Canvas, bounds image.Rectangle, cam *gfx.Camera) {
	s.RLock()
	for _, drawable := range s.S {
		drawable.DrawTo(c, bounds, cam)
	}
	s.RUnlock()
}

// New returns a new basic scene. It is short-handed for:
//  s := Scene(&Basic{
//      S: make([]gfx.Drawable, 0, 128)
//  })
func New() Scene {
	return Scene(&Basic{
		S: make([]gfx.Drawable, 0, 128),
	})
}
