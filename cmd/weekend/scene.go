package main

import "math/rand"


func blocksScene() []Object {
	var objects []Object

	// Table top
	objects = append(objects, NewBlock(NewVec3(0, -0.25, 0), 100, 100, 0.5, NewLambertian(NewVec3(0.7, 0.6, 0.5)), false))

	for x := -2; x < 2; x++ {
		for y := -2; y < 2; y++ {

			center := NewVec3(float64(x)+5.0*rand.Float64(), 0.2, float64(y)+5.0*rand.Float64())

			chooseMaterial := rand.Float64()
			var material Material
			if chooseMaterial < 0.7 {
				// diffuse
				material = NewLambertian(NewVec3(rand.Float64()*rand.Float64(),
					rand.Float64()*rand.Float64(),
					rand.Float64()*rand.Float64()))
			} else if chooseMaterial < 0.9 {
				// metal
				material = NewMetal(NewVec3(0.5*(1+rand.Float64()),
					0.5*(1+rand.Float64()),
					0.5*(1+rand.Float64())), 0.5*rand.Float64())
			} else {
				// glass
				refractionIndex := rand.Float64()/2.0 + 1.0
				material = NewDieletric(NewVec3(0.5, 0.5, 1.0), refractionIndex)
			}

			chooseShape := rand.Float64()
			var shape Object
			if chooseShape < 0.5 {
				shape = NewBlock(center, 0.4, 0.4, 0.4, material, false)
			} else {
				shape = NewSphere(center, 0.2, material)
			}
			objects = append(objects, shape)
		}
	}

	return objects
}

func smallScene() []Object {

	var objects []Object

	// Table top
	objects = append(objects, NewBlock(NewVec3(0, -0.25, 0), 100, 100, 0.5, NewLambertian(NewVec3(0.7, 0.6, 0.5)), false))
	//objects = append(objects, NewBlock(NewVec3(0, -0.25, 0), 100, 100, 0.5, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.03), false))

	objects = append(objects, NewSphere(NewVec3(-1, 1, -4), 1.0, NewDieletric(NewVec3(0.7, 0.5, 0.5),1.1)))

	block := NewBlock(NewVec3(0, 1, 0), 2.0, 2.0,2.0, NewLambertian(NewVec3(0.4, 0.2, 0.1)), false).
		RotateX(10).
		RotateY(10)

	objects = append(objects, block)
	objects = append(objects, NewSphere(NewVec3(1, 1, 4), 1.0, NewDieletric(NewVec3(0.5, 0.5, 1.0),1.5)))

	return objects
}


func bigScene() []Object {

	var objects []Object

	// Table top
	objects = append(objects, NewBlock(NewVec3(0, -0.25, 0), 100, 100, 0.5, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.03), false))

	//objects = append(objects, NewSphere(NewVec3(0, -1000, 0), 1000, NewLambertian(NewVec3(0.5, 0.5, 0.5))))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := NewVec3(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())
			if center.SubtractVec3(NewVec3(4, 0.2, 0)).Length() <= 0.9 {
				continue
			}

			chooseMaterial := rand.Float64()
			var material Material
			if chooseMaterial < 0.6 {
				// diffuse
				material = NewLambertian(NewVec3(rand.Float64()*rand.Float64(),
					rand.Float64()*rand.Float64(),
					rand.Float64()*rand.Float64()))
			} else if chooseMaterial < 0.85 {
				// metal
				material = NewMetal(NewVec3(0.5*(1+rand.Float64()),
					0.5*(1+rand.Float64()),
					0.5*(1+rand.Float64())), 0.5*rand.Float64())
			} else {
				// glass
				refractionIndex := rand.Float64()/2.0 + 1.0
				material = NewDieletric(NewVec3(1.0, 1.0, 1.0), refractionIndex)
			}

			chooseShape := rand.Float64()
			var shape Object
			if chooseShape < 0.25 {
				shape = NewBlock(center, 0.4, 0.4, 0.4, material, false)
			} else {
				shape = NewSphere(center, 0.2, material)
			}
			objects = append(objects, shape)
		}
	}

	objects = append(objects, NewSphere(NewVec3(0, 1, 0), 1.0, NewDieletric(NewVec3(1.0, 1.0, 1.0), 1.5)))
	objects = append(objects, NewSphere(NewVec3(-1, 1, -4), 1.0, NewLambertian(NewVec3(0.4, 0.2, 0.1))))
	objects = append(objects, NewSphere(NewVec3(1, 1, 4), 1.0, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)))
	//
	//objects = append(objects, NewSphere(NewVec3(0, 1, 0), 1.0, NewDieletric(1.5)))
	//objects = append(objects, NewSphere(NewVec3(-4, 1, 0), 1.0, NewLambertian(NewVec3(0.4, 0.2, 0.1))))
	//objects = append(objects, NewSphere(NewVec3(4, 1, 0), 1.0, NewMetal(NewVec3(0.7, 0.6, 0.5), 0.0)))

	return objects
}

