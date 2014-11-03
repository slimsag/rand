// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl2

import (
	"azul3d.org/v1/clock"
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/native/gl"
	"errors"
	"image"
	"runtime"
	"sync"
)

// Used when attempting to create an OpenGL 2.0 renderer in a lesser OpenGL context.
var ErrInvalidVersion = errors.New("invalid OpenGL version; must be at least OpenGL 2.0")

type pendingQuery struct {
	// The ID of the pending occlusion query.
	id uint32

	// The object of the pending occlusion query.
	o *gfx.Object
}

// Renderer is an OpenGL 2 based graphics renderer, it runs independant of the
// window management library being used (GLFW, SDL, Chippy, QML, etc).
//
// The renderer primarily uses two independant OpenGL contexts, one is used for
// rendering and one is used for managing resources like meshes, textures, and
// shaders which allows for asynchronous loading (although it is also possible
// to use only a single OpenGL context for windowing libraries that do not
// support multiple, but this will inheritly disable asynchronous loading).
type Renderer struct {
	// Render and loader execution channels.
	RenderExec chan func() bool
	LoaderExec chan func()

	// render/loader context API.
	render, loader *gl.Context

	// The graphics clock.
	clock *clock.Clock

	// The bounds of the gfx.Drawable, must be updated whenever the window size
	// changes.
	bounds struct {
		sync.RWMutex
		rect image.Rectangle
	}

	// GPU limitations.
	gpuInfo gfx.GPUInfo

	// Whether or not certain extensions we use are present or not.
	glArbMultisample, glArbFramebufferObject, glArbOcclusionQuery bool

	// Number of multisampling samples, buffers.
	samples, sampleBuffers int32

	// List of OpenGL texture compression format identifiers.
	compressedTextureFormats []int32

	// A channel which will have one empty struct inside it in the event that
	// a finalizer for a mesh, texture, etc has ran and something needs to be
	// free'd.
	wantFree chan struct{}

	// List of native meshes to free at next frame.
	meshesToFree struct {
		sync.RWMutex
		slice []*nativeMesh
	}

	// List of native shaders to free at next frame.
	shadersToFree struct {
		sync.RWMutex
		slice []*nativeShader
	}

	// List of native texture id's to free at next frame.
	texturesToFree struct {
		sync.RWMutex
		slice []uint32
	}

	// Structure for various previously set render states, mostly to avoid
	// useless OpenGL calls.
	last struct {
		scissor                                                                          image.Rectangle
		clearColor, blendColor                                                           gfx.Color
		clearDepth                                                                       float64
		clearStencil                                                                     int
		colorWrite                                                                       [4]bool
		depthFunc                                                                        gfx.Cmp
		blendFuncSeparate, blendEquationSeparate                                         gfx.BlendState
		stencilOpFront, stencilOpBack, stencilFuncFront, stencilFuncBack                 gfx.StencilState
		stencilMaskFront, stencilMaskBack                                                uint
		texturing, dithering, depthTest, depthWrite, stencilTest, blend, alphaToCoverage bool
		faceCulling                                                                      gfx.FaceCullMode
		program                                                                          uint32
		shader                                                                           *gfx.Shader
	}

	// Structure used to manage pending occlusion queries.
	pending struct {
		sync.Mutex
		queries []pendingQuery
	}

	// Whether or not our global OpenGL state has been set for this frame.
	stateSetForFrame bool

	// The precision of this renderer.
	precision gfx.Precision

	// The MSAA state.
	msaa struct {
		sync.RWMutex
		enabled bool
	}

	// Channel to wait for a Render() call to finish.
	renderComplete chan struct{}
}

// Implements gfx.Renderer interface.
func (r *Renderer) Clock() *clock.Clock {
	return r.clock
}

// Implements gfx.Renderer interface.
func (r *Renderer) Bounds() image.Rectangle {
	r.bounds.RLock()
	b := r.bounds.rect
	r.bounds.RUnlock()
	return b
}

// Implements gfx.Renderer interface.
func (r *Renderer) Clear(rect image.Rectangle, bg gfx.Color) {
	r.RenderExec <- func() bool {
		r.performClear(rect, bg)
		r.queryYield()
		return false
	}
}

// Implements gfx.Renderer interface.
func (r *Renderer) ClearDepth(rect image.Rectangle, depth float64) {
	r.RenderExec <- func() bool {
		r.performClearDepth(rect, depth)
		r.queryYield()
		return false
	}
}

// Implements gfx.Renderer interface.
func (r *Renderer) ClearStencil(rect image.Rectangle, stencil int) {
	r.RenderExec <- func() bool {
		r.performClearStencil(rect, stencil)
		r.queryYield()
		return false
	}
}

// Tries to receive pending occlusion query results, returns immediately if
// none are available yet. Returns the number of queries still pending.
func (r *Renderer) queryYield() int {
	if !r.glArbOcclusionQuery {
		return 0
	}
	r.pending.Lock()
	var (
		available, result int32
	)
	for queryIndex, query := range r.pending.queries {
		r.render.GetQueryObjectiv(query.id, gl.QUERY_RESULT_AVAILABLE, &available)
		r.render.Execute()
		if available == gl.TRUE {
			// Get the result then.
			r.render.GetQueryObjectiv(query.id, gl.QUERY_RESULT, &result)

			// Delete the query.
			r.render.DeleteQueries(1, &query.id)
			r.render.Execute()

			// Update object's sample count.
			nativeObj := query.o.NativeObject.(nativeObject)
			nativeObj.sampleCount = int(result)
			query.o.NativeObject = nativeObj

			// Remove from pending slice.
			r.pending.queries = append(r.pending.queries[:queryIndex], r.pending.queries[queryIndex+1:]...)
		}
	}
	length := len(r.pending.queries)
	r.pending.Unlock()
	return length
}

// Blocks until all pending occlusion query results are received.
func (r *Renderer) queryWait() {
	if !r.glArbOcclusionQuery {
		return
	}

	// We have no choice except to busy-wait until the results come: OpenGL
	// doesn't provide a blocking mechanism for waiting for query results but
	// at least we can runtime.Gosched() other goroutines.
	for i := 0; r.queryYield() > 0; i++ {
		// Only runtime.Gosched() every 16th iteration to avoid bogging down
		// rendering.
		if i != 0 && (i%16) == 0 {
			runtime.Gosched()
		}
	}
}

// Implements gfx.Renderer interface.
func (r *Renderer) QueryWait() {
	// Ask the render channel to wait for query results now.
	r.RenderExec <- func() bool {
		// Flush and execute any pending OpenGL commands.
		r.render.Flush()
		r.render.Execute()

		// Wait for occlusion query results to come in.
		r.queryWait()

		// signal render completion.
		r.renderComplete <- struct{}{}
		return false
	}
	<-r.renderComplete
}

// Implements gfx.Renderer interface.
func (r *Renderer) Render() {
	// If any finalizers have ran and actually want us to free something, then
	// we will ask the loader to do so now.
	r.LoaderExec <- func() {
		r.freeMeshes()
		r.freeShaders()
		r.freeTextures()
	}

	// Ask the render channel to render things now.
	r.RenderExec <- func() bool {
		// Execute all pending operations.
		for i := 0; i < len(r.RenderExec); i++ {
			f := <-r.RenderExec
			f()
		}

		// Flush and execute any pending OpenGL commands.
		r.render.Flush()
		r.render.Execute()

		// Clear our OpenGL state now.
		r.clearGlobalState()

		// Wait for occlusion query results to come in.
		r.queryWait()

		// Tick the clock.
		r.clock.Tick()

		// signal render completion.
		r.renderComplete <- struct{}{}
		return true
	}
	<-r.renderComplete
}

// Implements gfx.Renderer interface.
func (r *Renderer) Precision() gfx.Precision {
	return r.precision
}

// Implements gfx.Renderer interface.
func (r *Renderer) GPUInfo() gfx.GPUInfo {
	return r.gpuInfo
}

// Implements gfx.Canvas interface.
func (r *Renderer) SetMSAA(msaa bool) {
	r.msaa.Lock()
	r.msaa.enabled = msaa
	r.msaa.Unlock()
}

// Implements gfx.Canvas interface.
func (r *Renderer) MSAA() (msaa bool) {
	r.msaa.RLock()
	msaa = r.msaa.enabled
	r.msaa.RUnlock()
	return
}

func (r *Renderer) performClear(rect image.Rectangle, bg gfx.Color) {
	r.setGlobalState()

	// Only if the clear color has changed do we need to make the OpenGL call.
	r.stateScissor(rect)
	r.stateClearColor(bg)
	r.render.Clear(uint32(gl.COLOR_BUFFER_BIT))
}

func (r *Renderer) performClearDepth(rect image.Rectangle, depth float64) {
	r.setGlobalState()

	// Perform clearing.
	r.stateScissor(rect)
	r.stateClearDepth(depth)
	r.render.Clear(uint32(gl.DEPTH_BUFFER_BIT))
}

func (r *Renderer) performClearStencil(rect image.Rectangle, stencil int) {
	r.setGlobalState()

	// Perform clearing.
	r.stateScissor(rect)
	r.stateClearStencil(stencil)
	r.render.Clear(uint32(gl.STENCIL_BUFFER_BIT))
}

func (r *Renderer) UpdateBounds(bounds image.Rectangle) {
	r.bounds.Lock()
	r.bounds.rect = bounds
	r.bounds.Unlock()
}

func (r *Renderer) setGlobalState() {
	if !r.stateSetForFrame {
		r.stateSetForFrame = true

		// Update viewport bounds.
		r.bounds.Lock()
		r.render.Viewport(0, 0, uint32(r.bounds.rect.Dx()), uint32(r.bounds.rect.Dy()))
		r.bounds.Unlock()

		// Enable scissor testing.
		r.render.Enable(gl.SCISSOR_TEST)

		// Enable multisampling, if available and wanted.
		if r.glArbMultisample {
			r.msaa.RLock()
			msaa := r.msaa.enabled
			r.msaa.RUnlock()
			if msaa {
				r.render.Enable(gl.MULTISAMPLE)
			}
		}
	}
}

func (r *Renderer) clearGlobalState() {
	if r.stateSetForFrame {
		r.stateSetForFrame = false

		// Clear last used state.
		r.clearLastState()

		// Reset last shader so that uniforms are loaded again next frame.
		r.last.shader = nil

		// Disable scissor testing.
		r.render.Disable(gl.SCISSOR_TEST)

		// Disable multisampling, if available.
		if r.glArbMultisample {
			r.render.Disable(gl.MULTISAMPLE)
		}
	}
}

func New() (*Renderer, error) {
	r := &Renderer{
		RenderExec:     make(chan func() bool, 1024),
		LoaderExec:     make(chan func(), 1024),
		renderComplete: make(chan struct{}, 8),
		wantFree:       make(chan struct{}, 1),
		clock:          clock.New(),
	}

	// MSAA is enabled by default.
	r.msaa.enabled = true

	// Initialize r.render now.
	r.render = gl.New()
	r.render.SetBatching(false)
	if !r.render.AtLeastVersion(2, 0) {
		return nil, ErrInvalidVersion
	}

	// Initialize r.loader now.
	r.loader = gl.New()
	r.loader.SetBatching(false)

	// Note: we don't need r.ctx.Lock() here because no other goroutines
	// can be using r.ctx yet since we haven't returned from New().

	// Find the renderer's precision.
	var redBits, greenBits, blueBits, alphaBits, depthBits, stencilBits int32
	r.render.GetIntegerv(gl.RED_BITS, &redBits)
	r.render.GetIntegerv(gl.GREEN_BITS, &greenBits)
	r.render.GetIntegerv(gl.BLUE_BITS, &blueBits)
	r.render.GetIntegerv(gl.ALPHA_BITS, &alphaBits)
	r.render.GetIntegerv(gl.DEPTH_BITS, &depthBits)
	r.render.GetIntegerv(gl.STENCIL_BITS, &stencilBits)
	r.render.Execute()

	r.precision.RedBits = uint8(redBits)
	r.precision.GreenBits = uint8(greenBits)
	r.precision.BlueBits = uint8(blueBits)
	r.precision.AlphaBits = uint8(alphaBits)
	r.precision.DepthBits = uint8(depthBits)
	r.precision.StencilBits = uint8(stencilBits)

	// Query whether we have the GL_ARB_framebuffer_object.
	r.glArbFramebufferObject = r.render.Extension("GL_ARB_framebuffer_object")

	// Query whether we have the GL_ARB_occlusion_query.
	r.glArbOcclusionQuery = r.render.Extension("GL_ARB_occlusion_query")

	// Query whether we have the GL_ARB_multisample extension.
	r.glArbMultisample = r.render.Extension("GL_ARB_multisample")
	if r.glArbMultisample {
		// Query the number of samples and sample buffers we have, if any.
		r.render.GetIntegerv(gl.SAMPLES, &r.samples)
		r.render.GetIntegerv(gl.SAMPLE_BUFFERS, &r.sampleBuffers)
		r.render.Execute() // Needed because glGetIntegerv must execute now.
		r.precision.Samples = int(r.samples)
	}

	// Store GPU info.
	var maxTextureSize, maxVaryingFloats, maxVertexInputs, maxFragmentInputs, occlusionQueryBits int32
	r.render.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxTextureSize)
	r.render.GetIntegerv(gl.MAX_VARYING_FLOATS, &maxVaryingFloats)
	r.render.GetIntegerv(gl.MAX_VERTEX_UNIFORM_COMPONENTS, &maxVertexInputs)
	r.render.GetIntegerv(gl.MAX_FRAGMENT_UNIFORM_COMPONENTS, &maxFragmentInputs)
	if r.glArbOcclusionQuery {
		r.render.GetQueryiv(gl.SAMPLES_PASSED, gl.QUERY_COUNTER_BITS, &occlusionQueryBits)
	}
	r.render.Execute()

	r.gpuInfo.MaxTextureSize = int(maxTextureSize)
	r.gpuInfo.GLSLMaxVaryingFloats = int(maxVaryingFloats)
	r.gpuInfo.GLSLMaxVertexInputs = int(maxVertexInputs)
	r.gpuInfo.GLSLMaxFragmentInputs = int(maxFragmentInputs)
	r.gpuInfo.GLExtensions = r.render.Extensions()
	r.gpuInfo.AlphaToCoverage = r.glArbMultisample && r.samples > 0 && r.sampleBuffers > 0
	r.gpuInfo.Name = gl.String(r.render.GetString(gl.RENDERER))
	r.gpuInfo.Vendor = gl.String(r.render.GetString(gl.VENDOR))
	r.gpuInfo.GLMajor, r.gpuInfo.GLMinor, _ = r.render.Version()
	r.gpuInfo.GLSLMajor, r.gpuInfo.GLSLMinor, _ = r.render.ShaderVersion()
	r.gpuInfo.OcclusionQuery = r.glArbOcclusionQuery && occlusionQueryBits > 0
	r.gpuInfo.OcclusionQueryBits = int(occlusionQueryBits)

	// Grab the current renderer bounds (opengl viewport).
	var viewport [4]int32
	r.render.GetIntegerv(gl.VIEWPORT, &viewport[0])
	r.render.Execute()
	r.bounds.rect = image.Rect(0, 0, int(viewport[2]), int(viewport[3]))

	// Update scissor rectangle.
	r.stateScissor(r.bounds.rect)

	// Grab the number of texture compression formats.
	var numFormats int32
	r.render.GetIntegerv(gl.NUM_COMPRESSED_TEXTURE_FORMATS, &numFormats)
	r.render.Execute() // Needed because glGetIntegerv must execute now.

	// Store the slice of texture compression formats.
	if numFormats > 0 {
		r.compressedTextureFormats = make([]int32, numFormats)
		r.render.GetIntegerv(gl.COMPRESSED_TEXTURE_FORMATS, &r.compressedTextureFormats[0])
		r.render.Execute() // Needed because glGetIntegerv must execute now.
	}
	return r, nil
}
