package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

func render(ray Ray, world World, depth int) Vec3 {

	hit, material := world.Hit(ray, 0.001, math.MaxFloat64)

	if hit != nil {

		if depth <= 0 {
			return NewVec3(0, 0, 0)
		}

		color, scattered := material.Scatter(ray, hit)
		if color.IsZero() {
			// It's essentially black so don't keep bouncing
			return color
		} else {
			return color.MultiplyVec3(render(scattered, world, depth-1))
		}
	}

	// Compute the background color
	unitDirection := ray.Direction.UnitVector()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return NewVec3(1.0, 1.0, 1.0).
		MultiplyScalar(1.0 - t).
		AddVec3(NewVec3(0.5, 0.7, 1.0).MultiplyScalar(t))
}

func bookCover() World {

	var objects []Object

	// Add in the massive "world"
	objects = append(objects, NewSphere(NewVec3(0, -1000, 0), 1000, NewLambertian(NewVec3(0.5, 0.5, 0.5))))

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
				material = NewDieletric(1.5)
			}

			objects = append(objects, NewSphere(center, 0.2, material))
		}
	}

	objects = append(objects, NewSphere(NewVec3(0, 1, 0), 1.0, NewDieletric(1.5)))
	objects = append(objects, NewSphere(NewVec3(-4, 1, 0), 1.0, NewLambertian(NewVec3(0.4, 0.2, 0.1))))
	objects = append(objects, NewSphere(NewVec3(4, 1, 0), 1.0, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)))

	return World{Objects: objects}
}

func main() {

	nx := 1200
	ny := 800
	ns := 10
	depth := 10
	img := image.NewRGBA64(image.Rect(0, 0, nx, ny))

	rand.Seed(time.Now().UnixNano())

	//world := basicWorld()
	world := bookCover()

	lookFrom := NewVec3(13, 2, 3)
	lookAt := NewVec3(0, 0, 0)
	distToFocus := 10.0
	aperture := 0.05

	camera := NewCamera(lookFrom, lookAt, NewVec3(0, 1, 0), 20.0, float64(nx)/float64(ny), aperture, distToFocus)

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {

			var color Vec3

			for s := 0; s < ns; s++ {
				u := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)

				ray := camera.GetRay(u, v)
				color = color.AddVec3(render(ray, world, depth))
			}

			color = color.DivideScalar(float64(ns)).Gamma2()

			img.SetRGBA64(i, ny-(j+1), color.RGBA())
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
