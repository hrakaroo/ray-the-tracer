package main

import (
	"image/color"
	"math"
)

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

func (v Vec3) Length() float64 {
	return math.Sqrt(v.v1*v.v1 + v.v2*v.v2 + v.v3*v.v3)
}

func (v Vec3) SquaredLength() float64 {
	return v.v1*v.v1 + v.v2*v.v2 + v.v3*v.v3
}

func (v Vec3) UnitVector() Vec3 {
	l := v.Length()
	return NewVec3(v.v1/l, v.v2/l, v.v3/l)
}

func (v Vec3) RGBA() color.RGBA64 {
	return color.RGBA64{
		R: uint16(float64(0xfffe) * v.v1),
		G: uint16(float64(0xfffe) * v.v2),
		B: uint16(float64(0xfffe) * v.v3),
		A: 0xffff,
	}
}
