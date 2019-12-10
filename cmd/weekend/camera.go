package main

import (
	"math"
	"math/rand"
)

type Camera struct {
	LowerLeftCorner Vec3
	Horizontal      Vec3
	Vertical        Vec3
	Origin          Vec3
	LensRadius      float64
	U               Vec3
	V               Vec3
}

// vfov = vertical field of view from top to bottom in degrees
// aspect = aspect ratio
func NewCamera(lookFrom, lookAt, vup Vec3, vfov, aspect, aperture, focusDistance float64) *Camera {

	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight
	w := lookFrom.SubtractVec3(lookAt).UnitVector()
	u := vup.Cross(w).UnitVector()
	v := w.Cross(u)

	return &Camera{
		LowerLeftCorner: lookFrom.
			SubtractVec3(u.MultiplyScalar(halfWidth * focusDistance)).
			SubtractVec3(v.MultiplyScalar(halfHeight * focusDistance)).
			SubtractVec3(w.MultiplyScalar(focusDistance)),
		Horizontal: u.MultiplyScalar(halfWidth * focusDistance * 2.0),
		Vertical:   v.MultiplyScalar(halfHeight * focusDistance * 2.0),
		Origin:     lookFrom,
		LensRadius: aperture / 2.0,
		U:          u,
		V:          v,
	}
}

func (c *Camera) GetRay(s, t float64) Ray {
	rd := randomInUnitDisk().MultiplyScalar(c.LensRadius)
	offset := c.U.MultiplyScalar(rd.X()).AddVec3(c.V.MultiplyScalar(rd.Y()))
	return NewRay(c.Origin.AddVec3(offset),
		c.LowerLeftCorner.AddVec3(c.Horizontal.MultiplyScalar(s)).AddVec3(c.Vertical.MultiplyScalar(t)).SubtractVec3(c.Origin).SubtractVec3(offset))
}

func randomInUnitDisk() Vec3 {
	var unit Vec3
	done := false
	for !done {
		x := rand.Float64()*2.0 - 1.0
		y := rand.Float64()*2.0 - 1.0
		if x*x+y*y < 1.0 {
			unit = NewVec3(x, y, 0)
			done = true
		}
	}
	return unit
}
