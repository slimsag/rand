// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package window is the easiest way to open a window and render graphics.
package window

import (
	"azul3d.org/v1/chippy"
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/gfx/gl2"
	"fmt"
	"image"
	"log"
	"os"
	"runtime"
	"time"
)

func program(gfxLoop func(w *chippy.Window, r gfx.Renderer)) {
	defer chippy.Exit()

	window := chippy.NewWindow()
	window.SetTitle("Azul3D")
	screen := chippy.DefaultScreen()
	window.SetPositionCenter(screen)

	events := window.Events()

	// Actually open the windows
	err := window.Open(screen)
	if err != nil {
		log.Fatal(err)
	}

	// All OpenGL related calls must occur in the same OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Choose and set a frame buffer configuration.
	configs := window.GLConfigs()
	bestConfig := chippy.GLChooseConfig(configs, chippy.GLWorstConfig, chippy.GLBestConfig)
	log.Println("Chosen configuration:", bestConfig)
	window.GLSetConfig(bestConfig)

	// Create the OpenGL rendering context.
	ctx, err := window.GLCreateContext(2, 1, chippy.GLCoreProfile, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create the OpenGL loader context.
	loaderCtx, err := window.GLCreateContext(2, 1, chippy.GLCoreProfile, ctx)
	if err != nil {
		log.Fatal(err)
	}

	// OpenGL rendering context must be active to create the renderer.
	window.GLMakeCurrent(ctx)
	defer window.GLMakeCurrent(nil)

	// Disable vertical sync.
	//window.GLSetVerticalSync(chippy.NoVerticalSync)

	// Create the renderer.
	r, err := gl2.New()
	if err != nil {
		log.Fatal(err)
	}

	// Start the graphics rendering loop.
	go gfxLoop(window, r)

	// Channel to signal shutdown to renderer and loader.
	shutdown := make(chan bool, 2)

	// Start event loop.
	go func() {
		cl := r.Clock()
		printFPS := time.Tick(1 * time.Second)

		for {
			select {
			case <-printFPS:
				window.SetTitle(fmt.Sprintf("Azul3D %vFPS (%f Avg.)", cl.FrameRate(), cl.AverageFrameRate()))

			case e := <-events:
				switch ev := e.(type) {
				case chippy.ResizedEvent:
					r.UpdateBounds(image.Rect(0, 0, ev.Width, ev.Height))

				case chippy.CloseEvent, chippy.DestroyedEvent:
					shutdown <- true
					shutdown <- true
					return
				}
			}
		}
	}()

	// Start loading goroutine.
	go func() {
		// All OpenGL related calls must occur in the same OS thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		// OpenGL loading context must be active.
		window.GLMakeCurrent(loaderCtx)
		defer window.GLMakeCurrent(nil)

		for {
			select {
			case <-shutdown:
				return
			case fn := <-r.LoaderExec:
				fn()
			}
		}
	}()

	// Enter rendering loop.
	for {
		select {
		case <-shutdown:
			return

		case fn := <-r.RenderExec:
			if renderedFrame := fn(); renderedFrame {
				// Swap OpenGL buffers.
				window.GLSwapBuffers()
			}
		}
	}
}

func Run(gfxLoop func(w *chippy.Window, r gfx.Renderer)) {
	// Enable debug messages.
	chippy.SetDebugOutput(os.Stdout)

	// Initialize Chippy
	err := chippy.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer chippy.Exit()

	// Start the program.
	go program(gfxLoop)

	// Enter the main loop.
	chippy.MainLoop()
}
