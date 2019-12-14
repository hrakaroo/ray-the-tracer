package main

import (
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"
)

type ColoredPoint2D struct {
	X int
	Y int
	C color.RGBA64
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

func bookCover() []Object {

	var objects []Object

	// Table top
	objects = append(objects, NewBlock(NewVec3(0, -0.25, 0), 100, 100, 0.5, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.03)))

	//objects = append(objects, NewSphere(NewVec3(0, -1000, 0), 1000, NewLambertian(NewVec3(0.5, 0.5, 0.5))))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())
			if center.SubtractVec3(NewVec3(4, 0.2, 0)).Length() <= 0.9 {
				continue
			}

			chooseMaterial := rand.Float64()
			var material Material
			if chooseMaterial < 0.6 {
				// diffuse
				material = NewLambertian(NewVec3(rand.Float64()*rand.Float64(),
					rand.Float64()*rand.Float64(),
					rand.Float64()*rand.Float64()))
			} else if chooseMaterial < 0.85 {
				// metal
				material = NewMetal(NewVec3(0.5*(1+rand.Float64()),
					0.5*(1+rand.Float64()),
					0.5*(1+rand.Float64())), 0.5*rand.Float64())
			} else {
				// glass
				refractionIndex := rand.Float64()/2.0 + 1.0
				material = NewDieletric(refractionIndex)
			}

			chooseShape := rand.Float64()
			var shape Object
			if chooseShape < 0.25 {
				shape = NewBlock(center, 0.4, 0.4, 0.4, material)
			} else {
				shape = NewSphere(center, 0.2, material)
			}
			objects = append(objects, shape)
		}
	}

	objects = append(objects, NewSphere(NewVec3(0, 1, 0), 1.0, NewDieletric(1.5)))
	objects = append(objects, NewSphere(NewVec3(-1, 1, -4), 1.0, NewLambertian(NewVec3(0.4, 0.2, 0.1))))
	objects = append(objects, NewSphere(NewVec3(1, 1, 4), 1.0, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)))
	//
	//objects = append(objects, NewSphere(NewVec3(0, 1, 0), 1.0, NewDieletric(1.5)))
	//objects = append(objects, NewSphere(NewVec3(-4, 1, 0), 1.0, NewLambertian(NewVec3(0.4, 0.2, 0.1))))
	//objects = append(objects, NewSphere(NewVec3(4, 1, 0), 1.0, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)))

	return objects
}

func basicWorld() []Object {

	var objects []Object

	// Table top
	objects = append(objects, NewBlock(NewVec3(0, -0.25, 0), 100, 100, 0.5, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.03)))

	//objects = append(objects, NewSphere(NewVec3(0, 1, 0), 1.0, NewDieletric(1.5)))
	//objects = append(objects, NewSphere(NewVec3(-1, 1, -4), 1.0, NewLambertian(NewVec3(0.4, 0.2, 0.1))))
	objects = append(objects, NewBlock(NewVec3(0, 1, 0), 2.0, 2.0,2.0, NewLambertian(NewVec3(0.4, 0.2, 0.1))))
	//objects = append(objects, NewSphere(NewVec3(1, 1, 4), 1.0, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)))

	return objects
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
	rand.Seed(time.Now().UnixNano())

	width  := 2400
	height := 1600
	lookFrom := NewVec3(-5, 3, 20)
	lookAt := NewVec3(0, 0, 0)
	distToFocus := 20.0
	aperture := 0.05

	canvas := NewCanvas(width, height)

	environment := Environment{
		canvas:   canvas,
		camera:   NewCamera(lookFrom, lookAt, NewVec3(0, 1, 0), 20.0, canvas.AspectRatio(), aperture, distToFocus),
		sampling: 100,
		bounce:   30,
		objects:  bookCover(),
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
