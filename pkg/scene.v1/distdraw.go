package scene

import (
	"image"
	"sync"
)

// Dist implements the Scene interface using a simple slice.
type Dist struct {
	sync.RWMutex
	S []gfx.Drawable
}

// Implements the Scene interface.
func (s Dist) Add(d gfx.Drawable) bool {
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
func (s Dist) Has(d gfx.Drawable) bool {
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
func (s Dist) Remove(d gfx.Drawable) bool {
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
func (s Dist) Iter(callback func(d gfx.Drawable) (stop bool)) {
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
func (s Dist) DrawTo(c gfx.Canvas, bounds image.Rectangle, cam *gfx.Camera) {
	s.RLock()
	for _, drawable := range s.S {
		drawable.DrawTo(c, bounds, cam)
	}
	s.RUnlock()
}

// New returns a new Dist scene. It is short-handed for:
//  s := Scene(&Dist{
//      S: make([]gfx.Drawable, 0, 128)
//  })
func New() Scene {
	return Scene(&Dist{
		S: make([]gfx.Drawable, 0, 128),
	})
}
