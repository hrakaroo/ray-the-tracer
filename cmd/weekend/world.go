package main

type Hit struct {
	Scalar float64
	Point  Vec3
	Normal Vec3
}

type Object interface {
	Hit(ray Ray, tMin, tMax float64) (*Hit, Material)
}

type World struct {
	Objects []Object
}

func (w *World) Hit(ray Ray, tMin, tMax float64) (*Hit, Material) {

	var hit *Hit
	var material Material

	closestSoFar := tMax

	for _, object := range w.Objects {
		// Calculate the hit
		if tempHit, tempMaterial := object.Hit(ray, tMin, closestSoFar); tempHit != nil {
			hit = tempHit
			material = tempMaterial
			closestSoFar = tempHit.Scalar
		}
	}

	return hit, material
}
