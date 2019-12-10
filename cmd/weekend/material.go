package main

import (
	"math"
	"math/rand"
)

type Material interface {
	Scatter(ray Ray, hit *Hit) (Vec3, Ray)
}

/**
From Wikipedia:
Lambertian reflectance is the property that defines an ideal "matte" or diffusely reflecting surface.
The apparent brightness of a Lambertian surface to an observer is the same regardless of the observer's
angle of view.
*/
type Lambertian struct {
	Albedo Vec3
}

func NewLambertian(albedo Vec3) *Lambertian {
	return &Lambertian{Albedo: albedo}
}

func (l *Lambertian) Scatter(_ Ray, hit *Hit) (Vec3, Ray) {
	// Since apparent brightness is the same regardless of the observers angle, the
	//  incoming ray is not needed.
	// Our bounce is just the normal plus some random direction
	scattered := NewRay(hit.Point, hit.Normal.AddVec3(randomInUnitSphere()))
	return l.Albedo, scattered
}

type Metal struct {
	Albedo Vec3
	Fuzz   float64
}

func NewMetal(albedo Vec3, fuzz float64) *Metal {
	return &Metal{Albedo: albedo, Fuzz: fuzz}
}

func (m *Metal) Scatter(ray Ray, hit *Hit) (Vec3, Ray) {
	reflected := ray.Direction.UnitVector().Reflect(hit.Normal)
	scattered := NewRay(hit.Point, reflected.AddVec3(randomInUnitSphere().MultiplyScalar(m.Fuzz)))

	if scattered.Direction.Dot(hit.Normal) <= 0 {
		// There is no scatter
		return NewVec3(0.0, 0.0, 0.0), scattered
	}

	return m.Albedo, scattered
}

type Dieletric struct {
	RefractionIndex float64
}

func NewDieletric(refractionIndex float64) *Dieletric {
	return &Dieletric{RefractionIndex: refractionIndex}
}

func (d *Dieletric) Scatter(ray Ray, hit *Hit) (Vec3, Ray) {
	reflected := ray.Direction.Reflect(hit.Normal)
	color := NewVec3(1.0, 1.0, 1.0)

	var outwardNormal Vec3
	var niOverNt float64
	var cosine float64

	if ray.Direction.Dot(hit.Normal) > 0 {
		outwardNormal = hit.Normal.Negate()
		niOverNt = d.RefractionIndex
		//cosine = ray.Direction.Dot(hit.Normal) * d.RefractionIndex / ray.Direction.Length()
		cosine = ray.Direction.Dot(hit.Normal) / ray.Direction.Length()
		cosine = math.Sqrt(1 - d.RefractionIndex*d.RefractionIndex*(1-cosine*cosine))
	} else {
		outwardNormal = hit.Normal
		niOverNt = 1.0 / d.RefractionIndex
		cosine = -ray.Direction.Dot(hit.Normal) / ray.Direction.Length()
	}

	if refracts, refracted := ray.Direction.Refract(outwardNormal, niOverNt); refracts {
		if rand.Float64() > schlick(cosine, d.RefractionIndex) {
			return color, NewRay(hit.Point, refracted)
		}
	}

	return color, NewRay(hit.Point, reflected)
}

func schlick(cosine, refractionIndex float64) float64 {
	r0 := (1.0 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}

/**
Create a random direction in a unit sphere.
*/
func randomInUnitSphere() Vec3 {

	var unit Vec3

	// go doesn't support do/while loops
	done := false
	for !done {
		// Find a random point in a cube
		x := rand.Float64()*2.0 - 1.0
		y := rand.Float64()*2.0 - 1.0
		z := rand.Float64()*2.0 - 1.0

		// Check if the point is within a unit sphere where x^2 + y^2 + z^2 <= 1
		if x*x+y*y+z*z < 1.0 {
			unit = NewVec3(x, y, z)
			done = true
		}
	}
	return unit
}
