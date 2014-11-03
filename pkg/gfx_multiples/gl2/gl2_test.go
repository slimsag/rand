// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl2

import (
	"azul3d.org/v1/gfx"
	"testing"
)

func TestRendererInterface(t *testing.T) {
	var r *Renderer
	_ = gfx.Renderer(r)
}
