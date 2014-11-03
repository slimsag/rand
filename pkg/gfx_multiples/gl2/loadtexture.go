// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl2

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/native/gl"
	"image"
	"image/draw"
	"runtime"
	"unsafe"
)

type nativeTexture struct {
	r             *Renderer
	id            uint32
	width, height int
}

func (n *nativeTexture) Destroy() {
	finalizeTexture(n)
}

func finalizeTexture(n *nativeTexture) {
	n.r.texturesToFree.Lock()
	n.r.texturesToFree.slice = append(n.r.texturesToFree.slice, n.id)
	n.r.texturesToFree.Unlock()
}

func (n *nativeTexture) Download(rect image.Rectangle, complete chan image.Image) {
	if !n.r.glArbFramebufferObject {
		// We don't have GL_ARB_framebuffer_object extensions, we can't do
		// this at all.
		complete <- nil
		return
	}

	n.r.LoaderExec <- func() {
		// Create a FBO, bind it now.
		var fbo uint32
		n.r.loader.GenFramebuffers(1, &fbo)
		n.r.loader.Execute()
		n.r.loader.BindFramebuffer(gl.FRAMEBUFFER, fbo)

		// Attach the texture to the FBO.
		n.r.loader.FramebufferTexture2D(
			gl.FRAMEBUFFER,
			gl.COLOR_ATTACHMENT0,
			gl.TEXTURE_2D,
			n.id, // texture ID
			0,    // level
		)

		// If the rectangle is empty use the entire area.
		bounds := image.Rect(0, 0, n.width, n.height)
		if rect.Empty() {
			rect = bounds
		} else {
			// Intersect the rectangle with the texture's bounds.
			rect = bounds.Intersect(rect)
		}

		// Read texture pixels.
		img := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
		x, y, w, h := convertRect(rect, bounds)
		n.r.loader.ReadPixels(
			x, y, w, h,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(&img.Pix[0]),
		)

		// Delete the FBO.
		n.r.loader.BindFramebuffer(gl.FRAMEBUFFER, 0)
		n.r.loader.DeleteFramebuffers(1, &fbo)

		// Flush and execute.
		n.r.loader.Flush()
		n.r.loader.Execute()

		complete <- img
	}
}

func (r *Renderer) RenderToTexture(t *gfx.Texture) gfx.Canvas {
	return r
}

func verticalFlip(img *image.RGBA) {
	b := img.Bounds()
	rowCpy := make([]uint8, b.Dx()*4)
	for r := 0; r < (b.Dy() / 2); r++ {
		topRow := img.Pix[img.PixOffset(0, r):img.PixOffset(b.Dx(), r)]

		bottomR := b.Dy() - r - 1
		bottomRow := img.Pix[img.PixOffset(0, bottomR):img.PixOffset(b.Dx(), bottomR)]

		// Save bottom row.
		copy(rowCpy, bottomRow)

		// Copy top row to bottom row.
		copy(bottomRow, topRow)

		// Copy saved bottom row to top row.
		copy(topRow, rowCpy)
	}
}

func (r *Renderer) Download(rect image.Rectangle, complete chan image.Image) {
	r.RenderExec <- func() bool {
		bounds := r.Bounds()

		// If the rectangle is empty use the entire area.
		if rect.Empty() {
			rect = bounds
		} else {
			// Intersect the rectangle with the renderer's bounds.
			rect = bounds.Intersect(rect)
		}

		img := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
		x, y, w, h := convertRect(rect, bounds)
		r.render.ReadPixels(
			x, y, w, h,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(&img.Pix[0]),
		)

		// Flush and execute.
		r.render.Flush()
		r.render.Execute()

		// We must vertically flip the image.
		verticalFlip(img)

		// Yield for occlusion query results, if any are available.
		r.queryYield()

		complete <- img
		return false
	}
}

func convertWrap(w gfx.TexWrap) int32 {
	switch w {
	case gfx.Repeat:
		return gl.REPEAT
	case gfx.Clamp:
		return gl.CLAMP_TO_EDGE
	case gfx.BorderColor:
		return gl.CLAMP_TO_BORDER
	case gfx.Mirror:
		return gl.MIRRORED_REPEAT
	}
	panic("Invalid wrap mode")
}

func convertFilter(f gfx.TexFilter) int32 {
	switch f {
	case gfx.Nearest:
		return gl.NEAREST
	case gfx.Linear:
		return gl.LINEAR
	case gfx.NearestMipmapNearest:
		return gl.NEAREST_MIPMAP_NEAREST
	case gfx.LinearMipmapNearest:
		return gl.LINEAR_MIPMAP_NEAREST
	case gfx.NearestMipmapLinear:
		return gl.NEAREST_MIPMAP_LINEAR
	case gfx.LinearMipmapLinear:
		return gl.LINEAR_MIPMAP_LINEAR
	}
	panic("invalid filter.")
}

func (r *Renderer) freeTextures() {
	// Lock the list.
	r.texturesToFree.Lock()

	if len(r.texturesToFree.slice) > 0 {
		// Free the textures.
		r.loader.DeleteTextures(uint32(len(r.texturesToFree.slice)), &r.texturesToFree.slice[0])

		// Flush and execute OpenGL commands.
		r.loader.Flush()
		r.loader.Execute()
	}

	// Slice to zero, and unlock.
	r.texturesToFree.slice = r.texturesToFree.slice[:0]
	r.texturesToFree.Unlock()
}

const (
	// We really should try to get our GL bindings to wrap extensions..
	// See: http://www.opengl.org/registry/specs/EXT/texture_compression_s3tc.txt
	glCOMPRESSED_RGB_S3TC_DXT1_EXT  = 0x83F0
	glCOMPRESSED_RGBA_S3TC_DXT1_EXT = 0x83F1
	glCOMPRESSED_RGBA_S3TC_DXT3_EXT = 0x83F2
	glCOMPRESSED_RGBA_S3TC_DXT5_EXT = 0x83F3
)

func convertTexFormat(f gfx.TexFormat) int32 {
	switch f {
	case gfx.RGBA:
		return gl.RGBA
	case gfx.RGB:
		return gl.RGB
	case gfx.DXT1:
		return glCOMPRESSED_RGB_S3TC_DXT1_EXT
	case gfx.DXT1RGBA:
		return glCOMPRESSED_RGBA_S3TC_DXT1_EXT
	case gfx.DXT3:
		return glCOMPRESSED_RGBA_S3TC_DXT3_EXT
	case gfx.DXT5:
		return glCOMPRESSED_RGBA_S3TC_DXT5_EXT
	}
	panic("unknown format")
}

func (r *Renderer) LoadTexture(t *gfx.Texture, done chan *gfx.Texture) {
	// Lock the texture until we are done loading it.
	t.Lock()
	if t.Loaded {
		// Texture is already loaded, signal completion if needed and return
		// after unlocking.
		t.Unlock()
		select {
		case done <- t:
		default:
		}
		return
	}

	f := func() {
		// Create texture ID.
		native := &nativeTexture{
			r: r,
		}
		r.loader.GenTextures(1, &native.id)
		r.loader.Execute()

		// Bind the texture.
		r.loader.BindTexture(gl.TEXTURE_2D, native.id)

		// Set wrap mode.
		r.loader.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_BASE_LEVEL, 0)
		r.loader.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)

		// Determine appropriate internal image format.
		targetFormat := convertTexFormat(t.Format)
		internalFormat := gl.RGBA
		for _, format := range r.compressedTextureFormats {
			if format == targetFormat {
				internalFormat = format
				break
			}
		}

		// Upload the image.
		bounds := t.Source.Bounds()
		rgba, ok := t.Source.(*image.RGBA)
		if !ok {
			// Convert the image to RGBA.
			rgba = image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
			draw.Draw(rgba, rgba.Bounds(), t.Source, bounds.Min, draw.Src)
		}
		r.loader.TexImage2D(
			gl.TEXTURE_2D,
			0,
			internalFormat,
			uint32(bounds.Dx()),
			uint32(bounds.Dy()),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(&rgba.Pix[0]),
		)
		native.width = bounds.Dx()
		native.height = bounds.Dy()

		// Unbind texture to avoid carrying OpenGL state.
		r.loader.BindTexture(gl.TEXTURE_2D, 0)

		// Flush and execute OpenGL commands.
		r.loader.Flush()
		r.loader.Execute()

		// Mark the texture as loaded.
		t.Loaded = true
		t.NativeTexture = native
		t.ClearData()

		// Attach a finalizer to the texture that will later free it.
		runtime.SetFinalizer(native, finalizeTexture)

		// Unlock, signal completion, and return.
		t.Unlock()
		select {
		case done <- t:
		default:
		}
	}

	select {
	case r.LoaderExec <- f:
	default:
		go func() {
			r.LoaderExec <- f
		}()
	}
}
