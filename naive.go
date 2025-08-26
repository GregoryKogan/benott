package benott

// CountIntersectionsNaive calculates the number of intersections using a simple
// O(n^2) brute-force algorithm. It checks every pair of segments for intersections.
//
// This function is provided primarily as a testing and validation tool to verify
// the results of the more complex Bentley-Ottmann implementation. While robust,
// it is significantly slower for large numbers of segments.
func CountIntersectionsNaive(segments []Segment) int {
	count := 0
	for i := range segments {
		for j := i + 1; j < len(segments); j++ {
			s1 := segments[i]
			s2 := segments[j]

			// Don't count intersections at shared endpoints as a true "crossing".
			if segmentsShareEndpoint(&s1, &s2) {
				continue
			}

			// Check for a proper intersection using the CCW test.
			if segmentsIntersectCCW(s1.P1, s1.P2, s2.P1, s2.P2) {
				count++
			}
		}
	}
	return count
}

// segmentsShareEndpoint checks if two segments share a common vertex.
// Direct struct comparison works here because Point structs are comparable.
func segmentsShareEndpoint(s1, s2 *Segment) bool {
	return s1.P1 == s2.P1 || s1.P1 == s2.P2 || s1.P2 == s2.P1 || s1.P2 == s2.P2
}

// segmentsIntersectCCW tests for segment intersection using the counter-clockwise
// (CCW) orientation test. This is a standard geometric primitive that avoids
// floating-point division and is highly robust.
//
// It returns true if the segment (p1, p2) and the segment (p3, p4) have a
// "proper" intersection, meaning they cross each other at a single point.
func segmentsIntersectCCW(p1, p2, p3, p4 Point) bool {
	// This internal helper function determines the orientation of the ordered
	// triplet (u, v, w). The logic is adapted directly from your provided code.
	ccw := func(u, v, w Point) bool {
		// The formula calculates a value proportional to the signed area of the
		// triangle (u,v,w). The sign indicates whether the turn from vector
		// uv to vw is clockwise or counter-clockwise.
		return (w.Y-u.Y)*(v.X-u.X) > (v.Y-u.Y)*(w.X-u.X)
	}

	// The segments intersect if and only if the endpoints of each segment are
	// on opposite sides of the line containing the other segment.
	//
	// Check if p3 and p4 are on opposite sides of the line defined by p1 and p2.
	test1 := ccw(p1, p3, p4) != ccw(p2, p3, p4)
	// Check if p1 and p2 are on opposite sides of the line defined by p3 and p4.
	test2 := ccw(p1, p2, p3) != ccw(p1, p2, p4)

	return test1 && test2
}
