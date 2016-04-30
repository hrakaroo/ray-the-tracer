// Contains basic geometry concepts
package main

import (
	"math"
	"image/color"
)

type Row struct {
	X float64
	Y float64
	Z float64
	T float64
}

type Matrix struct {
	Rows [4]Row
}

type Point struct {
	X float64
	Y float64
	Z float64
}

func (row Row) multiplyPoint(point Point) float64 {
	//fmt.Println("Row: ", row)
	return row.X * point.X + row.Y * point.Y + row.Z * point.Z + row.T
}

func (matrix Matrix) multiplyPoint(point Point) Point {

	//fmt.Println("Pre point: ", point)

	x1 := matrix.Rows[0].multiplyPoint(point)
	y1 := matrix.Rows[1].multiplyPoint(point)
	z1 := matrix.Rows[2].multiplyPoint(point)

	point1 := Point{x1, y1, z1}
	//fmt.Println("Post Point: ", point1)

	return point1

	//return Point{x1, y1, z1}
}

// A vector is basically the same thing as a point with a source at the origin
type Vector Point

func (row Row) multiplyVector(vector Vector) float64 {
	return row.X * vector.X + row.Y * vector.Y + row.Z * vector.Z + row.T
}

func (matrix Matrix) multiplyVector(vector Vector) Vector {

	x1 := matrix.Rows[0].multiplyVector(vector)
	y1 := matrix.Rows[1].multiplyVector(vector)
	z1 := matrix.Rows[2].multiplyVector(vector)

	return Vector{x1, y1, z1}
}


// A ray is best represented by a source and a direction
//  The distance between the source and direction point
//  should be equal to 1 (unit vector)
type Ray struct {
	Source Point
	Direction Vector
}

func (matrix Matrix) multiplyRay(ray Ray) Ray {

	return Ray{matrix.multiplyPoint(ray.Source), matrix.multiplyVector(ray.Direction)}

}

// A plane is represented by a normal and a distance k from the origin
type Plane struct {
	Normal Vector
	K      float64
}


func square(x float64) float64 {
	return x * x
}

func unitVector(point Point) Vector {
	d := math.Sqrt(point.X * point.X + point.Y * point.Y + point.Z * point.Z)

	return Vector{point.X / d, point.Y / d, point.Z / d}
}

// Create a unit ray from two points
func unitRay(point1, point2 Point) Ray {

	// Distance for x, y, z
	dx := point2.X - point1.X
	dy := point2.Y - point1.Y
	dz := point2.Z - point1.Z

	// Total distance
	d := math.Sqrt(dx * dx + dy * dy + dz * dz)

	// Unit x, y, z
	ux := dx / d
	uy := dy / d
	uz := dz / d

	return Ray{point1, Vector{ux, uy, uz}}
}


// Multiply a ray by a magnitude to get a new point
func multiplyRay(ray Ray, magnitude float64) Point {

	x := ray.Source.X + ray.Direction.X * magnitude
	y := ray.Source.Y + ray.Direction.Y * magnitude
	z := ray.Source.Z + ray.Direction.Z * magnitude

	return Point{x, y, z}
}

func dotProduct(vector1, vector2 Vector) float64 {

	return vector1.X * vector2.X + vector1.Y * vector2.Y + vector1.Z * vector2.Z
}

type Shape interface {
	intersect(ray Ray) (bool, float64)
	normal(point Point) Vector
	getColor() color.RGBA64
}
