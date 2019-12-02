package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
)

/**
Create a random direction in a unit sphere.
 */
func randomInUnitSphere() Vec3 {

	var unit Vec3

	// go doesn't support do/while loops
	done := false
	for ! done {
		// Find a random point in a cube
		// todo - the book calls this a unit cube, but I'm not convinced it is
		x := rand.Float64() * 2.0 - 1.0
		y := rand.Float64() * 2.0 - 1.0
		z := rand.Float64() * 2.0 - 1.0

		// Check if the point is within a unit sphere where x^2 + y^2 + z^2 <= 1
		if x*x + y*y + z*z <= 1.0 {
			unit = NewVec3(x, y, z)
			done = true
		}
	}
	return unit
}

func myColor(ray Ray, world World) Vec3 {

	hit := world.Hit(ray, 0.001, math.MaxFloat64)

	if hit != nil {
		target := hit.Point.AddVec3(hit.Normal).AddVec3(randomInUnitSphere())

		return myColor(NewRay(hit.Point, target.SubtractVec3(hit.Point)), world).MultiplyScalar(0.5)
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

	world := World{[]Object{
		NewSphere(NewVec3(0, 0, -1), 0.5),
		NewSphere(NewVec3(0, -100.5, -1), 100)},
	}

	camera := NewCamera()

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {

			var color Vec3

			for s := 0; s < ns; s++ {
				h := (float64(i) + rand.Float64()) / float64(nx)
				v := (float64(j) + rand.Float64()) / float64(ny)

				ray := camera.GetRay(h, v)
				color = color.AddVec3(myColor(ray, world))
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
