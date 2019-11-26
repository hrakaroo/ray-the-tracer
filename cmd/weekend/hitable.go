package main

type HitRecord struct {
	Scalar float64
	Point Vec3
	Normal Vec3
}

type Hitable interface {
	Hit(ray Ray, tMin, tMax float64) (bool, *HitRecord)
}

type HitableList struct {
	Hitables  []Hitable
}

func (h *HitableList) Hit(ray Ray, tMin, tMax float64) (bool, *HitRecord) {

	var hitRecord *HitRecord

	for _, hitable := range h.Hitables {
		// Calculate the hit
		hit, tempHit := hitable.Hit(ray, tMin, tMax)
		if ! hit {
			continue
		}

		// Check for closest hit
		if hitRecord == nil || tempHit.Scalar < hitRecord.Scalar {
			hitRecord = tempHit
		}
	}

	return hitRecord != nil, hitRecord
}
