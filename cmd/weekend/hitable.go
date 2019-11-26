package main

type Hit struct {
	Scalar float64
	Point Vec3
	Normal Vec3
}

type Object interface {
	ComputeHit(ray *Ray, tMin, tMax float64) *Hit
}


type World struct {
	Objects  []Object
}

func (w *World) Hit(ray *Ray, tMin, tMax float64) *Hit {

	var hit *Hit

	for _, object := range w.Objects {
		// Calculate the hit
		if tempHit := object.ComputeHit(ray, tMin, tMax); tempHit != nil {
			// Check if its the closest
			if hit == nil || tempHit.Scalar < hit.Scalar {
				hit = tempHit
			}
		}
	}

	return hit
}
