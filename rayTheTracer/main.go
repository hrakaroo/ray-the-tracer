package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var img = image.NewRGBA64(image.Rect(0, 0, 800, 800))


func shadeColor(c color.RGBA64, percent float64) color.RGBA64 {

	red := uint16(float64(c.R) * percent)
	green := uint16(float64(c.G) * percent)
	blue := uint16(float64(c.B) * percent)

	return color.RGBA64{red, green, blue, 65535}
}

func averageColor(dots []color.RGBA64) color.RGBA64 {

	r := 0
	g := 0
	b := 0
	a := 0

	for _, c := range dots {
		r += int(c.R)
		g += int(c.G)
		b += int(c.B)
		a += int(c.A)
	}

	l := len(dots)

	return color.RGBA64{uint16(r / l), uint16(g / l), uint16(b / l), uint16(a / l)}
}



func main() {

	 //red := color.RGBA64{0xffff, 0, 0, 0xffff}
	 //green := color.RGBA64{0, 0xffff, 0, 0xffff}
	//black := color.RGBA64{0, 0, 0, 0xffff}
	white := color.RGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	blue := color.RGBA64{0, 0xff, 0xffff, 0xffff}
	ambient_coefficient := .2
	diffuse_coefficient := .8

	 //ball1 := Sphere{Point{250.0, 100.0, 0.0}, 50.0, red}
	 //ball2 := Sphere{Point{150.0, 125.0, 0.0}, 80.0, green}

	// cube := Cube{[8]Ray{top, bottom, left, right, front, back}, blue}
	cube := makeCube(Point{0, 0, 0}, 50, blue)

	// shapes := []Shape{ball1, ball2, cube}
	//shapes := []Shape{ball1, ball2}
	shapes := []Shape{cube}

	eye := Point{150, -200, -400.0}
	light := Point{100, -400, -400}

	halfx := img.Rect.Dx() / 2
	halfy := img.Rect.Dy() / 2

	// Walk the screen left to right, top to bottom (or bottom to top, I'm not sure)
	for x := 0; x < img.Rect.Dx(); x++ {
		for y := 0; y < img.Rect.Dy(); y++ {

			// Translate our x and y so that 0, 0 is in the center
			tx := float64(x - halfx)
			ty := float64(y - halfy)

			// This will hold all of the dots from our jitter
			var dots []color.RGBA64

			// Our jitter
			for xp := -2; xp < 5; xp++ {
				for yp := -2; yp < 5; yp++ {

					hits := make(map[float64]color.RGBA64)

					for _, shape := range shapes {

						xx := tx + float64(xp) * 0.3
						yy := ty + float64(yp) * 0.3

						// Create a unit vector from our eye
						ray := unitRay(eye, Point{xx, yy, 0.0})

						// Calculate the point of intersect
						hit, m := shape.intersect(ray)

						if ! hit {
							continue
						}

						p := multiplyRay(ray, m)

						// Normal unit vector out of the sphere
						normal := shape.normal(p)

						shade := - dotProduct(normal, unitRay(light, p).Direction)
						if ( shade < 0 ) {
							shade = 0
						}

						point_color := shadeColor(shape.getColor(),
							ambient_coefficient + diffuse_coefficient * shade)

						hits[m] = point_color
					}

					color := white

					if len(hits) > 0 {

						set := false
						var last float64

						// Find the closest one
						for k, v := range hits {
							if ! set || k < last {
								last = k
								color = v
								set = true
							}
						}
					}

					dots = append(dots, color)
				}
			}

			img.Set(x, y, averageColor(dots))
		}
	}

	name := "drawing.png"

	fmt.Println("Writing file", name)
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
