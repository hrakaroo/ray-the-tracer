package main

import (
	"fmt"
	goimage "image"
	gocolor "image/color"
	"image/png"
	"os"
)

type Canvas struct {
	Width int
	Height int
	Image *goimage.RGBA64
}

/**
Some things of note
x - right (+) and left (-)
y - up (+) and down (-)
z - towards (+) and away (-)
 */
type Environment struct {
	canvas *Canvas
	camera *Camera
	objects World
	sampling int
	bounce int
}

func NewCanvas(width, height int) *Canvas {
	return &Canvas{Width: width, Height: height, Image: goimage.NewRGBA64(goimage.Rect(0, 0, width, height))}
}

func (c *Canvas) AspectRatio() float64 {
	return float64(c.Width)/float64(c.Height)
}

func (c * Canvas) Draw(x, y int, color gocolor.RGBA64) {
	c.Image.SetRGBA64(x, c.Height - ( y + 1 ), color)
}


func (c *Canvas) write(name string) {

	fmt.Println("Writing file", name)
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, c.Image)
	if err != nil {
		panic("Failed to encode image")
	}

	err = f.Close()
	if err != nil {
		panic("Failed to close file")
	}
}