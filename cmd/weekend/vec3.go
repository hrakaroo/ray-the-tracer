package main

import (
	"image/color"
	"math"
)

/**
Our generic 3 value struct.  We will use this for points, colors and directions.
*/
type Vec3 struct {
	v1 float64
	v2 float64
	v3 float64
}

func NewVec3(v1, v2, v3 float64) Vec3 {
	return Vec3{v1, v2, v3}
}

func (v *Vec3) X() float64 {
	return v.v1
}

func (v *Vec3) Y() float64 {
	return v.v2
}

func (v *Vec3) Z() float64 {
	return v.v3
}

func (v *Vec3) R() float64 {
	return v.v1
}

func (v *Vec3) G() float64 {
	return v.v2
}

func (v *Vec3) B() float64 {
	return v.v3
}

func (v *Vec3) Negate() Vec3 {
	return NewVec3(-v.v1, -v.v2, -v.v3)
}

func (v Vec3) AddVec3(o Vec3) Vec3 {
	return NewVec3(v.v1+o.v1, v.v2+o.v2, v.v3+o.v3)
}

func (v Vec3) SubtractVec3(o Vec3) Vec3 {
	return NewVec3(v.v1-o.v1, v.v2-o.v2, v.v3-o.v3)
}

func (v Vec3) MultiplyVec3(o Vec3) Vec3 {
	return NewVec3(v.v1*o.v1, v.v2*o.v2, v.v3*o.v3)
}

func (v Vec3) DivideVec3(o Vec3) Vec3 {
	return NewVec3(v.v1/o.v1, v.v2/o.v2, v.v3/o.v3)
}

func (v Vec3) MultiplyScalar(scalar float64) Vec3 {
	return NewVec3(v.v1*scalar, v.v2*scalar, v.v3*scalar)
}

func (v Vec3) DivideScalar(scalar float64) Vec3 {
	return NewVec3(v.v1/scalar, v.v2/scalar, v.v3/scalar)
}

func (v Vec3) Dot(o Vec3) float64 {
	return v.v1*o.v1 + v.v2*o.v2 + v.v3*o.v3
}

func (v Vec3) Cross(o Vec3) Vec3 {
	return NewVec3(v.v2*o.v3-v.v3*o.v2,
		-(v.v1*o.v3 - v.v3*o.v1),
		v.v1*o.v2-v.v2*o.v1)
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.v1*v.v1 + v.v2*v.v2 + v.v3*v.v3)
}

func (v Vec3) UnitVector() Vec3 {
	l := v.Length()
	return NewVec3(v.v1/l, v.v2/l, v.v3/l)
}

func (v Vec3) Gamma2() Vec3 {
	return NewVec3(math.Sqrt(v.v1), math.Sqrt(v.v2), math.Sqrt(v.v3))
}

/**
Is essentially close enough to zero
*/
func (v Vec3) IsZero() bool {
	return math.Abs(v.v1) < 0.00001 && math.Abs(v.v2) < 0.00001 && math.Abs(v.v3) < 0.00001
}

/**
From the book:
The reflected ray direction is just v+2B where B = dot(v, N).  The
subtract is to point v out.
*/
func (v Vec3) Reflect(normal Vec3) Vec3 {
	return v.SubtractVec3(normal.MultiplyScalar(v.Dot(normal) * 2.0))
}

func (v Vec3) Refract(normal Vec3, niOverNt float64) (bool, Vec3) {
	uv := v.UnitVector()
	dt := uv.Dot(normal)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		return true, uv.SubtractVec3(normal.MultiplyScalar(dt)).MultiplyScalar(niOverNt).SubtractVec3(normal.MultiplyScalar(math.Sqrt(discriminant)))
	}
	return false, Vec3{}
}

func (v Vec3) RGBA() color.RGBA64 {
	return color.RGBA64{
		R: uint16(float64(0xfffe) * v.v1),
		G: uint16(float64(0xfffe) * v.v2),
		B: uint16(float64(0xfffe) * v.v3),
		A: 0xffff,
	}
}
