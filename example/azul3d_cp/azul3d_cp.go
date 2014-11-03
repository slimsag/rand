package main

import cp "azul3d.org/native/cp.v1"
import "time"
import "fmt"

var fs []func(body *cp.Body, gravity cp.Vect, damping, dt float64)

func spacey(space *cp.Space) {
	for i := 0; i < 100; i++ {
		body := cp.BodyNew(10, cp.MomentForCircle(1, 0, 32, cp.V(0, 0)))
		shape := body.CircleShapeNew(32, cp.Vect{})
		space.AddBody(body)
		space.AddShape(shape)

		f := func(body *cp.Body, gravity cp.Vect, damping, dt float64) {}
		fs = append(fs, f)
		body.SetVelocityUpdateFunc(f)
	}
}

func main() {
	space := cp.SpaceNew()
	spacey(space)

	for {
		time.Sleep(160 * time.Millisecond)
		spacey(space)
		space.Step(1)
		a := make([]byte, 1024*1024)
		b := a
		_ = b
	}

	for _, fs := range fs {
		fmt.Println(fs)
	}
}
