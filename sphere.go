package math

// Sphere describes a 3D sphere composed of a center point and radius.
type Sphere struct {
	Center Vec3
	Radius float64
}

// Contains tells if the point p is within this sphere.
func (s Sphere) Contains(p Vec3) bool {
	return s.Center.Sub(p).LengthSq() < s.Radius
}

// In tells if the sphere s inside the sphere b.
func (s Sphere) In(b Sphere) bool {
	return false
}

// Overlaps reports whether s and b have a non-empty intersection.
func (s Sphere) Overlaps(b Sphere) bool {
	// Real-Time Collision Detection, 4.3.1:
	//  Sphere-sphere Intersection

	// Calculate squared distance between centers.
	dist := s.Center.Sub(b.Center)
	dist2 := dist.Dot(dist)

	// Spheres intersect if squared distance is less than squared sum of radii.
	radiusSum := s.Radius + b.Radius
	return dist2 <= radiusSum*radiusSum
}

// Rect3 returns a 3D rectangle encapsulating this sphere.
func (s Sphere) Rect3() Rect3 {
	return Rect3{
		Min: s.Center.SubScalar(s.Radius),
		Max: s.Center.AddScalar(s.Radius),
	}
}

// InRect3 reports whether the sphere s is completely inside the rectangle r.
// It is short-hand for:
//  s.Rect3().In(r)
func (s Sphere) InRect3(r Rect3) bool {
	return s.Rect3().In(r)
}

// OverlapsRect3 reports whether the sphere s has a non-empty intersection with
// the rectangle r.
func (s Sphere) OverlapsRect3(r Rect3) bool {
	// Real-Time Collision Detection, 5.2.5:
	//  Testing Sphere Against AABB

	// Compute squared distance to the center of s.
	sqDist := r.SqDistToPoint(s.Center)

	// Sphere and AABB have an intersection if the squared distance between
	// them is less than the squared sphere radius.
	return sqDist <= s.Radius*s.Radius
}
