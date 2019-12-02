package main

import (
	"math/rand"
)

type Material interface {
	Scatter(ray Ray, hit *Hit) (Vec3, Ray)
}

/**
aka Matte
*/
type Lambertian struct {
	Color Vec3
}

func NewLambertian(color Vec3) *Lambertian {
	return &Lambertian{Color: color}
}

func (l *Lambertian) Scatter(ray Ray, hit *Hit) (Vec3, Ray) {
	target := hit.Point.AddVec3(hit.Normal).AddVec3(randomInUnitSphere())
	scattered := NewRay(hit.Point, target.SubtractVec3(hit.Point))
	return l.Color, scattered
}

type Metal struct {
	Color Vec3
	Fuzz  float64
}

func NewMetal(albedo Vec3, fuzz float64) *Metal {
	return &Metal{Color: albedo, Fuzz: fuzz}
}

func (m *Metal) Scatter(ray Ray, hit *Hit) (Vec3, Ray) {
	reflected := ray.Direction.UnitVector().Reflect(hit.Normal)
	scattered := NewRay(hit.Point, reflected.AddVec3(randomInUnitSphere().MultiplyScalar(m.Fuzz)))

	if scattered.Direction.Dot(hit.Normal) <= 0 {
		// There is no scatter
		return NewVec3(0.0, 0.0, 0.0), scattered
	}

	return m.Color, scattered
}

/**
Create a random direction in a unit sphere.
*/
func randomInUnitSphere() Vec3 {

	var unit Vec3

	// go doesn't support do/while loops
	done := false
	for !done {
		// Find a random point in a cube
		// todo - the book calls this a unit cube, but I'm not convinced it is
		x := rand.Float64()*2.0 - 1.0
		y := rand.Float64()*2.0 - 1.0
		z := rand.Float64()*2.0 - 1.0

		// Check if the point is within a unit sphere where x^2 + y^2 + z^2 <= 1
		if x*x+y*y+z*z <= 1.0 {
			unit = NewVec3(x, y, z)
			done = true
		}
	}
	return unit
}
