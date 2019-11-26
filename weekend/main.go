package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func hitSphere(center Vec3, radius float64, ray Ray) bool {
	oc := ray.Origin.SubtractVec3(center)
	a := ray.Direction.Dot(ray.Direction)
	b := 2.0 * oc.Dot(ray.Direction)
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant > 0
}

func myColor(ray Ray) Vec3 {
	if hitSphere(NewVec3(0, 0, -1), -0.5, ray) {
		return NewVec3(1, 0, 0)
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

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {

			h := float64(i) / float64(nx)
			v := float64(j) / float64(ny)

			ray := NewRay(origin, lowerLeftCorner.AddVec3(horizontal.MultiplyScalar(h)).AddVec3(vertical.MultiplyScalar(v)))

			col := myColor(ray)

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
