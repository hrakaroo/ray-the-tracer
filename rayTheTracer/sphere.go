package main

import (
	"image/color"
	"math"
)

// All points on the sphere satisfy
//     x^2 + y^2 + z^2 = r^2
//  basically the distance equation.
//  for a sphere at point p the equation becomes
//     (x-p.x)^2 + (y-p.y)^2 + (z-p.z)^2 = r^2
//
type Sphere struct {
	Center Point
	Radius float64
	Color  color.RGBA64
}

func (sphere Sphere) getColor() color.RGBA64 {
	return sphere.Color
}

func (sphere Sphere) normal(point Point) Ray {
	return unitRay(sphere.Center, point)
}

func (sphere Sphere) intersect(ray Ray) (bool, float64) {

	// The point at which the ray intersects the sphere is
	//  where the two equations are equal

	// Sphere (x - cx)^2 + (y - cy)^2 + (z - cz)^2 = r^2
	//  where cx,cy,cz is the center of the sphere
	//
	// Ray (x,y,z) = (sx,sy,sz) + m (dx,dy,dz)
	// where (sx,sy,sz) = source
	//       (dx,dy,dz) = direction
	//                m = magnitude
	// which breaks down to
	//         x = sx + m * dx
	//         y = sy + m * dy
	//         z = sz + m * dz
	//
	// (m * dx + sx - cx)^2 + (m * dy + sy - cy)^2 + (m * dz + sz - cz)^2 = r^2
	//
	// to simplify let
	//     nx = sx - cx
	//     ny = sy - cy
	//     nz = sz - cz
	//
	// (m * dx + nx)^2 + (m * dy + ny)^2 + (m * dz + nz)^2
	//
	// m^2dx^2 + 2m*dx*nx + nx^2 +
	// m^2dy^2 + 2m*dy*ny + ny^2 +
	// m^2dz^2 + 2m*dz*nz + nz^2 = r^2
	//
	// solving for m we get
	//
	// (dx^2 + dy^2 + dz^2)m^2 + (2dx*nx + 2dy*ny + 2dz*nz)m +
	//     nx^2 + ny^2 + nz^2 - r2 = 0
	//
	// Solve with quadratic equation

	nx := ray.Source.X - sphere.Center.X
	ny := ray.Source.Y - sphere.Center.Y
	nz := ray.Source.Z - sphere.Center.Z

	//a := ray.direction.x * ray.direction.x +
	//     ray.direction.y * ray.direction.y +
	//     ray.direction.z * ray.direction.z

	// Since we are dealing with unit vector a == 1

	b := 2.0*ray.Direction.X*nx +
		2.0*ray.Direction.Y*ny +
		2.0*ray.Direction.Z*nz

	c := nx*nx +
		ny*ny +
		nz*nz -
		sphere.Radius*sphere.Radius

	// Quadratic equation.  -b +/- sqrt(b^2 - 4ac) / 2a
	//   again, a == 1 so this reduces to
	//                      -b +/- sqrt(b^2 - 4c) / 2
	// if the value in sqrt() is negative then it is imaginary and there is no
	//  intersection
	//

	g := b*b - 4.0*c

	// if this value is less than zero there is no solution so no hit
	// if it is equal zero then there is one intersection point
	if g < 0 {
		return false, 0
	}

	// Take the sqrt root
	g = math.Sqrt(g)

	if g == 0 {
		m := -b / 2.0
		return true, m
	} else {
		m1 := (-b - g) / 2.0
		m2 := (-b + g) / 2.0

		// The smaller magnitude is closer so order the points
		if m1 > m2 {
			t := m1
			m1 = m2
			m2 = t
		}

		return true, m1
	}
}
