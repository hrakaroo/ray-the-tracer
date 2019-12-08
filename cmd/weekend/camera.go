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
	U Vec3
	V Vec3
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
		Horizontal:      u.MultiplyScalar(halfWidth * focusDistance * 2.0),
		Vertical:        v.MultiplyScalar(halfHeight * focusDistance * 2.0),
		Origin:          lookFrom,
		LensRadius:      aperture / 2.0,
		U: u,
		V: v,
	}
}

func (c *Camera) GetRay(s, t float64) Ray {
	rd :=  randomInUnitDisk().MultiplyScalar(c.LensRadius)
	offset := c.U.MultiplyScalar(rd.X()).AddVec3(c.V.MultiplyScalar(rd.Y()))
	return NewRay(c.Origin.AddVec3(offset),
		c.LowerLeftCorner.AddVec3(c.Horizontal.MultiplyScalar(s)).AddVec3(c.Vertical.MultiplyScalar(t)).SubtractVec3(c.Origin).SubtractVec3(offset))
}

func randomInUnitDisk() Vec3 {
	var p Vec3
	done := false
	for !done {
		p := NewVec3(rand.Float64(), rand.Float64(), 0).MultiplyScalar(2.0).SubtractVec3(NewVec3(1, 1, 0))
		if p.Dot(p) < 1.0 {
			done = true
		}
	}
	return p
}