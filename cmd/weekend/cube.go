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
	Center   Vec3
	Planes   [6]Plane
	Material Material
}

func NewCube(center Vec3, size float64, material Material) *Cube {
	bottom := NewPlane(NewVec3(0, -1, 0), size/2.0)
	top := NewPlane(NewVec3(0, 1, 0), size/2.0)
	left := NewPlane(NewVec3(-1, 0, 0), size/2.0)
	right := NewPlane(NewVec3(1, 0, 0), size/2.0)
	front := NewPlane(NewVec3(0, 0, -1), size/2.0)
	back := NewPlane(NewVec3(0, 0, 1), size/2.0)

	return &Cube{
		Center:   center,
		Planes:   [6]Plane{bottom, top, left, right, front, back},
		Material: material,
	}
}

func (c *Cube) Hit(ray Ray, tMin, tMax float64) (*Hit, Material) {

	// Suppose a normal (xn, yn, zn) and a value k such that every point on the plane is given by
	//  x*xn + y*yn + z*zn = k
	//
	// As an example, imagine the normal (0, 0, 1) and k=4.  This is a plane hovering at z=4.
	//  Every point on the plane can be expressed by x*0 + y*0 + z*1 = 4, in other words, it
	//  doesn't mater the x and y as long as z=4, so a plane hovering at z=4
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

	// We need to move the cube to the origin so we can rotate it.  The easiest way
	//  is simply to ignore its center.  But then we also need to move the ray
	//  the same distance so create a copy ray which is translated
	rayCopy := ray.SubtractVec3(c.Center)

	// Compute the intersection with each plane
	var hits []*Hit
	for _, plane := range c.Planes {

		// Determine where our ray intersects the plane
		dot := rayCopy.Direction.Dot(plane.Normal)

		if dot == 0.0 {
			// Okay, so this ray is either above, below, or contained in the plane.
			// Figure out the distance from the center to the plane
			//  m = - k / (xn*xn + yn*yn + zn*zn)

			plane_d := plane.K / (square(plane.Normal.X()) + square(plane.Normal.Y()) + square(plane.Normal.Z()))

			// Now determine the distance from the line to the center.
			m := -(rayCopy.Origin.X()*rayCopy.Direction.X() + rayCopy.Origin.Y()*rayCopy.Direction.Y() + rayCopy.Origin.Z()*rayCopy.Direction.Z())
			closestPoint := ray.PointAt(m)
			line_d := math.Sqrt(square(closestPoint.X()) + square(closestPoint.Y()) + square(closestPoint.Z()))

			if line_d >= plane_d {
				return nil, nil
			}

			// Otherwise, just skip this as we will intersect some other plane.
			continue
		}

		n := plane.K - (rayCopy.Origin.X()*plane.Normal.X() +
			rayCopy.Origin.Y()*plane.Normal.Y() +
			rayCopy.Origin.Z()*plane.Normal.Z())

		m := n / dot

		// For calculating the hit, use the actual ray
		hits = append(hits, &Hit{Scalar: m, Point: ray.PointAt(m), Normal: plane.Normal})
	}

	sort.Slice(hits, func(i, j int) bool {
		return hits[i].Scalar < hits[j].Scalar
	})

	// Assume the first hit is going in
	in := true

	var lastHit *Hit
	for _, hit := range hits {

		dot := rayCopy.Direction.Dot(hit.Normal)

		if dot < 0 {
			// Normal is pointing in opposite direction of ray so this is an IN
			if in {
				// The last one was also going in so this is our new last_m
				lastHit = hit
			} else {
				// The last one was going out and this is an IN so we missed
				return nil, nil
			}
		} else {
			// Normal and ray are pointing in the same direction so this is an OUT
			in = false
		}
	}

	//log.Printf("%v", lastHit)
	//log.Printf("%v", c.Material)

	if lastHit.Scalar < tMax && lastHit.Scalar > tMin {
		return lastHit, c.Material
	}

	return nil, nil
}

func square(x float64) float64 {
	return x * x
}