package benott

import "math"

const epsilon = 1e-9

// Point represents a point in 2D space.
type Point struct {
	X, Y float64
}

// Segment represents a line segment in 2D space.
type Segment struct {
	P1, P2 Point
}

// intersection finds the intersection point of two line segments.
// It returns the point and true if they intersect within their bounds,
// otherwise an empty point and false. This implementation focuses on
// simple, non-collinear intersections.
func (s1 Segment) intersection(s2 Segment) (Point, bool) {
	p1, p2 := s1.P1, s1.P2
	p3, p4 := s2.P1, s2.P2

	// Vector representation of the segments
	r := Point{X: p2.X - p1.X, Y: p2.Y - p1.Y}
	s := Point{X: p4.X - p3.X, Y: p4.Y - p3.Y}

	// Cross product of the direction vectors
	rxs := r.X*s.Y - r.Y*s.X

	// Vector from p1 to p3
	qp := Point{X: p3.X - p1.X, Y: p3.Y - p1.Y}

	// If rxs is zero, the lines are parallel or collinear.
	// A full implementation might handle collinear overlaps, but for counting
	// discrete intersection points, we can ignore this case.
	if math.Abs(rxs) < epsilon {
		return Point{}, false
	}

	// Solve for t and u such that p1 + t*r = p3 + u*s
	t := (qp.X*s.Y - qp.Y*s.X) / rxs
	u := (qp.X*r.Y - qp.Y*r.X) / rxs

	// An intersection exists if t and u are between 0 and 1 (inclusive).
	if (t >= -epsilon && t <= 1+epsilon) && (u >= -epsilon && u <= 1+epsilon) {
		intersectionPoint := Point{X: p1.X + t*r.X, Y: p1.Y + t*r.Y}
		return intersectionPoint, true
	}

	return Point{}, false
}
