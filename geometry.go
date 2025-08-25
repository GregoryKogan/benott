package benott

import "math"

// epsilon is a small constant used for floating-point comparisons to handle
// precision errors. It defines the tolerance for considering two numbers equal.
const epsilon = 1e-9

// Point represents a point in a 2D Cartesian coordinate system.
type Point struct {
	X, Y float64
}

// Segment represents a line segment in 2D space, defined by two endpoints, P1 and P2.
type Segment struct {
	P1, P2 Point
}

// intersection calculates the intersection point of two line segments, s1 and s2.
// It uses a standard vector cross-product method to solve for the intersection.
//
// The method returns the intersection Point and a boolean `true` if the segments
// intersect within their finite bounds. If they are parallel, collinear, or
// intersect outside their bounds, it returns a zero Point and `false`.
func (s1 Segment) intersection(s2 Segment) (Point, bool) {
	p1, p2 := s1.P1, s1.P2
	p3, p4 := s2.P1, s2.P2

	// r and s are the direction vectors of the segments.
	r := Point{X: p2.X - p1.X, Y: p2.Y - p1.Y}
	s := Point{X: p4.X - p3.X, Y: p4.Y - p3.Y}

	// rxs is the cross product of the direction vectors. If it's zero, the
	// lines are parallel or collinear and do not intersect in a single point.
	rxs := r.X*s.Y - r.Y*s.X
	if math.Abs(rxs) < epsilon {
		return Point{}, false
	}

	// qp is the vector from the start of s1 to the start of s2.
	qp := Point{X: p3.X - p1.X, Y: p3.Y - p1.Y}

	// t and u are the scalar parameters along the lines. The intersection point
	// is at p1 + t*r and p3 + u*s.
	t := (qp.X*s.Y - qp.Y*s.X) / rxs
	u := (qp.X*r.Y - qp.Y*r.X) / rxs

	// An intersection exists only if both t and u are between 0 and 1,
	// meaning the intersection point lies on both finite segments.
	if (t >= -epsilon && t <= 1+epsilon) && (u >= -epsilon && u <= 1+epsilon) {
		intersectionPoint := Point{X: p1.X + t*r.X, Y: p1.Y + t*r.Y}
		return intersectionPoint, true
	}

	return Point{}, false
}
