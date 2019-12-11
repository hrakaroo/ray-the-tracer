package main

import (
	"math"
	"sort"
)

// A plane is a normal vector and a multiplier along that vector
type Plane struct {
	Normal Vec3
	K      float64
}

func NewPlane(normal Vec3, k float64) Plane {
	return Plane{
		Normal: normal,
		K:      k,
	}
}

type Cube struct {
	Planes   [6]Plane
	Material Material
}

func NewCube(center Vec3, size float64, material Material) *Cube {
	bottom := NewPlane(NewVec3(0, -1, 0), size)
	top := NewPlane(NewVec3(0, 1, 0), size)
	left := NewPlane(NewVec3(-1, 0, 0), size)
	right := NewPlane(NewVec3(1, 0, 0), size)
	front := NewPlane(NewVec3(0, 0, -1), size)
	back := NewPlane(NewVec3(0, 0, 1), size)

	return &Cube{
		Planes:   [6]Plane{bottom, top, left, right, front, back},
		Material: material,
	}
}

func (c *Cube) Hit(ray Ray, tMin, tMax float64) (*Hit, Material) {

	// Suppose a normal (xn, yn, zn) and a value k such that every point on the plane is given by
	//  x*xn + y*yn + z*zn = k
	//
	// If the normal to the plane is (0, 1, 0) with k = 4.
	//  Every point on the plane is given by x*0, y*1, z*0 = 4
	//
	// Ray (x,y,z) = (sx,sy,sz) + m (dx,dy,dz)
	// where (sx,sy,sz) = source
	//       (dx,dy,dz) = direction
	//                m = magnitude
	// which breaks down to
	//         x = sx + m * dx
	//         y = sy + m * dy
	//         z = sz + m * dz
	//
	// substituting that into the normal we get
	// (sx + m * dx) * xn + (sy + m * dy) * yn + (sz + m * dz) * zn = k
	//
	// sx * xn + m * dx * xn + sy * yn + m * dy * yn + sz * zn + m * dz * zn = k
	//
	// m (dx * xn + dy * yn + dz * zn) + (sx * xn + sy * yn + sz * zn) + k = 0
	//
	// m = (k - (sx * xn + sy * yn + sz * zn)) / (dx * xn + dy * yn + dz * zn)
	// m = - (sx * xn + sy * yn + sz * zn + k) / (dx * xn + dy * yn + dz * zn)

	// The one issue here is that if the dot product of the ray direction and
	//  the plane normal is zero it means that the ray is parallel or contained
	//  within the plane.  For this case we need to determine the distance from
	//  the center to the line and from the center to the plane and determine
	//  which is closer.

	// For the distance between the point and the plane, lets use the normal
	//  x*xn + y*yn + z*zn = k
	// There exists some point on the plane such that some magnitude of
	//  the normal plus this point intersects the origin.
	//    (x, y, z) + m * N = (0, 0, 0)
	//    x + m * xn = 0  => x = - m * xn
	//    y + m * yn = 0  => y = - m * yn
	//    z + m * zn = 0  => z = - m * zn
	// Substituting
	//  (-m * xn * xn) + (-m * yn * yn) + (-z * zn * zn) = k
	// Solve for m
	//  - m (xn*xn + yn*yn + zn*zn) = k
	//  m = - k / (xn*xn + yn*yn + zn*zn)
	// Since we really just care about distance we can take the absolute value and drop the negative
	//  Interestingly, as long as K is guaranteed to be positive we can drop the absolute all together
	//  m = k / (xn*xn + yn*yn + zn*zn)
	//
	//
	// For the distance between a point (the origin) and a line we know that x, y, z is on the line so
	//   rx + m * dx = x
	//   ry + m * dy = y
	//   rz + m * dz = z
	//
	//  We also know that a the dot product between a vector from that point and the origin and
	//  the original line must be zero as they are perpendicular
	//  (x, y, z) dot ( dx, dy, dz) = 0
	//   x * dx + y * dy + z * dz = 0
	//
	//  Substituting in the line info
	//   (rx + m * dx) * dx + (ry + m * dy) * dy + (rz + m * dz) * dz = 0
	//   rx * dx + m * dx^2 + ry * dy + m * dy^2 + rz * dz + m * dz^2 = 0
	//   m ( dx^2 + dy^2 + dz^2 ) = - (rx * dx + ry * dy + rz * dz)
	//   m = - (rx * dx + ry * dy + rz * dz) / ( dx^2 + dy^2 + dz^2 )
	//  Since ray.direction is a unit vector ( dx^2 + dy^2 + dz^2 ) == 1
	//   m = - (rx * dx + ry * dy + rz * dz)
	//  From here we can calculate the point and then the distance to the origin
	//

	hits := []Hit{}
	//hits := Hits{hits: make([]Hit, 6)}

	for _, plane := range c.Planes {

		// Determine where our ray intersects the plane
		dot := ray.Direction.Dot(plane.Normal)

		//fmt.Println("d: ", dot)

		if dot == 0.0 {
			// Okay, so this ray is either above, below, or contained in the plane.
			// Figure out the distance from the center to the plane
			//  m = - k / (xn*xn + yn*yn + zn*zn)

			plane_d := plane.K / (square(plane.Normal.X()) + square(plane.Normal.Y()) + square(plane.Normal.Z()))

			// Now determine the distance from the line to the center.
			m := -(ray.Origin.X()*ray.Direction.X() + ray.Origin.Y()*ray.Direction.Y() + ray.Origin.Z()*ray.Direction.Z())
			closestPoint := ray.PointAt(m)
			line_d := math.Sqrt(square(closestPoint.X()) + square(closestPoint.Y()) + square(closestPoint.Z()))

			if line_d >= plane_d {
				return nil, nil
			}
			continue
		}

		n := plane.K - (ray.Origin.X()*plane.Normal.X() +
			ray.Origin.Y()*plane.Normal.Y() +
			ray.Origin.Z()*plane.Normal.Z())

		m := n / dot

		hits = append(hits, Hit{Scalar: m, Point: ray.PointAt(m), Normal: plane.Normal})
	}

	sort.Sort(hits)

	// Assume the first hit is going in
	in := true

	last_m := 0.0
	for i := 0; i < hits.Len(); i++ {
		hit := hits.Get(i)

		if hit.dot < 0 {
			// Normal is pointing in opposite direction of ray so this is an IN
			if in {
				// The last one was also going in so this is our new last_m
				last_m = hit.m
			} else {
				// The last one was going out and this is an IN so we missed
				return false, 0
			}
		} else {
			// Normal and ray are pointing in the same direction so this is an OUT
			in = false
		}
	}

	return true, last_m
}

func square(x float64) float64 {
	return x * x
}