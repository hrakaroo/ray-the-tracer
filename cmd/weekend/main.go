package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
)

func myColor(ray Ray, world World, depth int) Vec3 {

	hit, material := world.Hit(ray, 0.001, math.MaxFloat64)

	if hit != nil {

		if depth > 50 {
			return NewVec3(0, 0, 0)
		}

		color, scattered := material.Scatter(ray, hit)
		if color.IsZero() {
			// It's essentially black so don't keep bouncing
			return color
		} else {
			return color.MultiplyVec3(myColor(scattered, world, depth+1))
		}
	}

	// Compute the background color
	unitDirection := ray.Direction.UnitVector()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return NewVec3(1.0, 1.0, 1.0).
		MultiplyScalar(1.0 - t).
		AddVec3(NewVec3(0.5, 0.7, 1.0).MultiplyScalar(t))
}

func main() {

	nx := 200
	ny := 100
	ns := 100
	img := image.NewRGBA64(image.Rect(0, 0, nx, ny))

	world := World{
		[]Object{
			NewSphere(NewVec3(0, 0, -1), 0.5, NewLambertian(NewVec3(0.1, 0.2, 0.5))),
			NewSphere(NewVec3(0, -100.5, -1), 100, NewLambertian(NewVec3(0.8, 0.8, 0.0))),
			NewSphere(NewVec3(1, 0, -1), 0.5, NewMetal(NewVec3(0.8, 0.6, 0.2), 0.0)),
			NewSphere(NewVec3(-1, 0, -1), 0.5, NewDieletric(1.5)),
			NewSphere(NewVec3(-1, 0, -1), -0.45, NewDieletric(1.5)),
		},
	}

	lookFrom := NewVec3(3, 3, 2)
	lookAt := NewVec3( 0, 0, -1)
	distToFocus := lookFrom.SubtractVec3(lookAt).Length()
	aperture := 2.0

	camera := NewCamera(lookFrom, lookAt, NewVec3(0, 1, 0),20.0, float64(nx)/float64(ny), aperture, distToFocus)

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {

			var color Vec3

			for s := 0; s < ns; s++ {
				h := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)

				ray := camera.GetRay(h, v)
				color = color.AddVec3(myColor(ray, world, 0))
			}

			color = color.DivideScalar(float64(ns)).Gamma2()

			img.Set(i, ny-(j+1), color.RGBA())
		}
	}

	name := "drawing.png"

	fmt.Println("Writing file", name)
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, img)
	if err != nil {
		panic("Failed to encode image")
	}

	err = f.Close()
	if err != nil {
		panic("Failed to close file")
	}
}
