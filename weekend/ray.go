package main

type Ray struct {
	Origin    Vec3
	Direction Vec3
}

func NewRay(origin Vec3, direction Vec3) Ray {
	return Ray{Origin: origin, Direction: direction}
}

func (r *Ray) PointAt(scalar float64) Vec3 {
	return r.Origin.AddVec3(r.Direction.MultiplyScalar(scalar))
}
