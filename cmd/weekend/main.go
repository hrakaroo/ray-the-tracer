package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)


func myColor(ray *Ray, world World) Vec3 {

	hit := world.Hit(ray, 0.0, math.MaxFloat64)

	if hit != nil {
		return NewVec3(hit.Normal.X() + 1, hit.Normal.Y()+1, hit.Normal.Z()+1).MultiplyScalar(0.5)
	}

	unitDirection := ray.Direction.UnitVector()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return NewVec3(1.0, 1.0, 1.0).
		MultiplyScalar(1.0 - t).
		AddVec3(NewVec3(0.5, 0.7, 1.0).MultiplyScalar(t))
}

func main() {

	nx := 200
	ny := 100
	img := image.NewRGBA64(image.Rect(0, 0, nx, ny))

	lowerLeftCorner := NewVec3(-2.0, -1.0, -1.0)
	horizontal := NewVec3(4.0, 0.0, 0.0)
	vertical := NewVec3(0.0, 2.0, 0.0)
	origin := NewVec3(0.0, 0.0, 0.0)

	world := World{[]Object{
		NewSphere(NewVec3(0, 0, -1), 0.5),
		NewSphere(NewVec3(0, -100.5, -1), 100)},
	}

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {

			h := float64(i) / float64(nx)
			v := float64(j) / float64(ny)

			ray := NewRay(origin, lowerLeftCorner.AddVec3(horizontal.MultiplyScalar(h)).AddVec3(vertical.MultiplyScalar(v)))

			col := myColor(ray, world)

			//unitDirection := ray.Direction.UnitVector()
			//t := 0.5 * (unitDirection.Y() + 1.0)
			//col := NewVec3(1.0, 1.0, 1.0).MultiplyScalar(1.0 - t).AddVec3(NewVec3(0.5, 0.7, 1.0).MultiplyScalar(t))

			img.Set(i, ny-(j+1), col.Color())
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
