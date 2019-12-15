package main

import (
	gocolor "image/color"
	"math"
	"math/rand"
	"sync"
)

type ColoredPoint2D struct {
	X int
	Y int
	C gocolor.RGBA64
}

type Point2D struct {
	X int
	Y int
}

func render(ray Ray, objects World, depth int) Vec3 {

	hit, material := objects.Hit(ray, 0.001, math.MaxFloat64)

	if hit != nil {

		if depth <= 0 {
			return NewVec3(0, 0, 0)
		}

		color, scattered := material.Scatter(ray, hit)
		if color.IsZero() {
			// It's essentially black so don't keep bouncing
			return color
		} else {
			return color.MultiplyVec3(render(scattered, objects, depth-1))
		}
	}

	// Compute the background color
	unitDirection := ray.Direction.UnitVector()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return NewVec3(1.0, 1.0, 1.0).
		MultiplyScalar(1.0 - t).
		AddVec3(NewVec3(0.5, 0.7, 1.0).MultiplyScalar(t))
}


func renderPoint(renderChan chan Point2D, environment Environment, drawChan chan ColoredPoint2D) {

	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for point := range renderChan {
				var c Vec3

				for s := 0; s < environment.sampling; s++ {
					u := (float64(point.X) + rand.Float64()) / float64(environment.canvas.Width)
					v := (float64(point.Y) + rand.Float64()) / float64(environment.canvas.Height)

					ray := environment.camera.GetRay(u, v)
					c = c.AddVec3(render(ray, environment.objects, environment.bounce))
				}

				c = c.DivideScalar(float64(environment.sampling)).Gamma2()

				drawChan <- ColoredPoint2D{X: point.X, Y: point.Y, C: c.RGBA()}
			}
		}()
	}

	wg.Wait()
}

func main() {

	// I'm not 100% sure I really need to do this.  Without it the random numbers will always be the same
	//  which isn't going to really negatively impact the rendered output, except that multiple runs
	//  will create the same image.
	//rand.Seed(time.Now().UnixNano())

	width  := 1200
	height := 800
	lookFrom := NewVec3(-5, 3, 20)
	lookAt := NewVec3(0, 0, 0)
	distToFocus := 20.0
	aperture := 0.05

	canvas := NewCanvas(width, height)

	environment := Environment{
		canvas:   canvas,
		camera:   NewCamera(lookFrom, lookAt, NewVec3(0, 1, 0), 20.0, canvas.AspectRatio(), aperture, distToFocus),
		sampling: 40,
		bounce:   10,
		objects:  smallScene(),
	}

	// Serialize drawing to the actual image as I don't know/doubt it can handle multiple write requests
	drawChan := make(chan ColoredPoint2D, 500)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for point := range drawChan {
			canvas.Draw(point.X, point.Y, point.C)
		}
		wg.Done()
	}()

	// Create a channel to pull the points to render from
	renderChan := make(chan Point2D, 100)

	// Start slamming points into the channel
	go func() {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				renderChan <- Point2D{X: x, Y: y}
			}
		}
		// All done sending points to render
		close(renderChan)
	}()

	// Our main rendering, this blocks until all points have been rendered
	renderPoint(renderChan, environment, drawChan)

	// Everything has been sent to the draw channel so close it off
	close(drawChan)

	// Wait for everything to finish drawing to the canvas
	wg.Wait()

	// Write our image
	canvas.write("drawing.png")
}
