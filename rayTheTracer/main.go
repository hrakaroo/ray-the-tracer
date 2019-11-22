package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"runtime/pprof"
	"sync"
)

var (
	red   = color.RGBA64{R: 0xffff, A: 0xffff}
	green = color.RGBA64{G: 0xffff, A: 0xffff}
	black = color.RGBA64{A: 0xffff}
	white = color.RGBA64{R: 0xffff, G: 0xffff, B: 0xffff, A: 0xffff}
	blue  = color.RGBA64{G: 0xff, B: 0xffff, A: 0xffff}
)

type Environment struct {
	ambientCoefficient float64
	diffuseCoefficient float64
	shapes             []Shape
	width              int
	height             int
	eye                Point
	light              Point
}

type Point2D struct {
	X int
	Y int
}

type ColoredPoint2D struct {
	X     int
	Y     int
	color color.RGBA64
}

func shadeColor(c color.RGBA64, percent float64) color.RGBA64 {

	red := uint16(float64(c.R) * percent)
	green := uint16(float64(c.G) * percent)
	blue := uint16(float64(c.B) * percent)

	return color.RGBA64{R: red, G: green, B: blue, A: 65535}
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

	return color.RGBA64{R: uint16(r / l), G: uint16(g / l), B: uint16(b / l), A: uint16(a / l)}
}

func main() {

	f, err := os.Create("run.prof")
	if err != nil {
		panic("Could not create prof")
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic("failed to start profile")
	}
	defer pprof.StopCPUProfile()

	//ball1 := Sphere{Point{250.0, 100.0, 0.0}, 50.0, red}
	//ball2 := Sphere{Point{150.0, 125.0, 0.0}, 80.0, green}

	// cube := Cube{[8]Ray{top, bottom, left, right, front, back}, blue}
	cube := makeCube(Point{0, 0, 0}, 50, blue)

	// shapes := []Shape{ball1, ball2, cube}
	//shapes := []Shape{ball1, ball2}

	img := image.NewRGBA64(image.Rect(0, 0, 800, 800))

	drawChan := make(chan ColoredPoint2D, 500)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for point := range drawChan {
			img.Set(point.X, point.Y, point.color)
		}
		wg.Done()
	}()

	environment := Environment{
		ambientCoefficient: 0.2,
		diffuseCoefficient: 0.8,
		shapes:             []Shape{cube},
		width:              img.Rect.Dx(),
		height:             img.Rect.Dy(),
		eye:                Point{150, -200, -400.0},
		light:              Point{100, -400, -400},
	}

	renderChan := make(chan Point2D, 5)
	go renderPoint(renderChan, environment, drawChan)

	// Walk the screen left to right, top to bottom
	//  (or bottom to top, I'm not sure)
	for x := 0; x < img.Rect.Dx(); x++ {
		for y := 0; y < img.Rect.Dy(); y++ {
			renderChan <- Point2D{x, y}
		}
	}
	close(renderChan)

	wg.Wait()

	name := "drawing.png"

	fmt.Println("Writing file", name)
	f, err = os.Create(name)
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

func renderPoint(renderChan chan Point2D, environment Environment, drawChan chan ColoredPoint2D) {

	halfx := environment.width / 2
	halfy := environment.height / 2

	var wg sync.WaitGroup

	for point := range renderChan {
		wg.Add(1)
		go func(point Point2D) {
			defer wg.Done()
			// Translate our x and y so that 0, 0 is in the center
			tx := float64(point.X - halfx)
			ty := float64(point.Y - halfy)

			// This will hold all of the dots from our jitter
			dots := make([]color.RGBA64, 9)

			// Our jitter
			index := 0
			for xp := -1; xp < 2; xp++ {
				for yp := -1; yp < 2; yp++ {

					//xx := tx + float64(xp)*0.3
					//yy := ty + float64(yp)*0.3
					xx := tx
					yy := ty

					hits := make(map[float64]color.RGBA64)

					for _, shape := range environment.shapes {

						// Create a unit vector from our eye
						ray := unitRay(environment.eye, Point{xx, yy, 0.0})

						// Calculate the point of intersect
						hit, m := shape.intersect(ray)

						if !hit {
							continue
						}

						p := multiplyRay(ray, m)

						// Normal unit vector out of the sphere
						normal := shape.normal(p)

						shade := -dotProduct(normal, unitRay(environment.light, p).Direction)
						if shade < 0 {
							shade = 0
						}

						pointColor := shadeColor(shape.getColor(),
							environment.ambientCoefficient+environment.diffuseCoefficient*shade)

						hits[m] = pointColor
					}

					pointColor := white

					if len(hits) > 0 {

						set := false
						var last float64

						// Find the closest one
						for k, v := range hits {
							if !set || k < last {
								last = k
								pointColor = v
								set = true
							}
						}
					}

					dots[index] = pointColor
					index++
				}
			}

			drawChan <- ColoredPoint2D{X: point.X, Y: point.Y, color: averageColor(dots)}
		}(point)
	}

	wg.Wait()
	close(drawChan)
}
