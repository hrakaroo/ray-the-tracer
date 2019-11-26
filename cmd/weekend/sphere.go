package main

import "math"

type Sphere struct {
	Center Vec3
	Radius float64
}

func NewSphere(center Vec3, radius float64) *Sphere {
	return &Sphere{
		Center: center,
		Radius: radius,
	}
}

func (s *Sphere) ComputeHit(ray *Ray, tMin, tMax float64) *Hit {

	oc := ray.Origin.SubtractVec3(s.Center)

	a := ray.Direction.Dot(ray.Direction)
	b := oc.Dot(ray.Direction)
	c := oc.Dot(oc) - s.Radius*s.Radius

	discriminant := b*b - a*c

	if discriminant < 0 {
		// miss
		return nil
	}

	// Compute both points
	// Note from the book: I eliminated a bunch of redundant 2's that cancel each other out
	sqrt := math.Sqrt(b*b - a*c)
	scalar1 := (-b - sqrt)/a
	scalar2 := (-b + sqrt)/a

	// scalar1 is closer than scalar2 so test it first
	for _, scalar := range []float64{scalar1, scalar2} {
		if scalar < tMax && scalar > tMin {
			point := ray.PointAt(scalar)

			return &Hit{
				Scalar: scalar,
				Point:  point,
				Normal: point.SubtractVec3(s.Center).DivideScalar(s.Radius),
			}
		}
	}

	// There was a hit, but it was outside our tMin - tMax range
	return nil
}
