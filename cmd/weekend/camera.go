package main

type Camera struct {
	 LowerLeftCorner Vec3
	 Horizontal Vec3
	 Vertical Vec3
	 Origin Vec3
}

// todo - I don't understand why the range for x is 4 but 2 for y
// todo - I think it's to match the dimensions of the drawing image so the image itself doesn't look warped
func NewCamera() *Camera {
	return &Camera{
		LowerLeftCorner: NewVec3(-2.0, -1.0, -1.0),
		Horizontal:      NewVec3(4.0, 0.0, 0.0),
		Vertical:        NewVec3(0.0, 2.0, 0.0),
		Origin:          NewVec3(0.0, 0.0, 0.0),
	}
}

func (c *Camera) GetRay(h, v float64) Ray {
	return NewRay(c.Origin, c.LowerLeftCorner.AddVec3(c.Horizontal.MultiplyScalar(h)).AddVec3(c.Vertical.MultiplyScalar(v)).SubtractVec3(c.Origin))
}
