package benott

import (
	"container/heap"
	"math"
)

// CountIntersections counts the number of intersection points in a set of segments.
func CountIntersections(segments []Segment) int {
	eq := make(EventQueue, 0)
	segmentCopies := make([]Segment, len(segments))
	copy(segmentCopies, segments)

	for i := range segmentCopies {
		p1, p2 := segmentCopies[i].P1, segmentCopies[i].P2
		if p1.X > p2.X || (p1.X == p2.X && p1.Y > p2.Y) {
			segmentCopies[i].P1, segmentCopies[i].P2 = p2, p1
		}
		heap.Push(&eq, &Event{Point: segmentCopies[i].P1, Type: SegmentStart, Segments: []*Segment{&segmentCopies[i]}})
		heap.Push(&eq, &Event{Point: segmentCopies[i].P2, Type: SegmentEnd, Segments: []*Segment{&segmentCopies[i]}})
	}

	status := NewStatus()
	intersections := 0
	processedIntersections := make(map[Point]bool)

	for eq.Len() > 0 {
		event := heap.Pop(&eq).(*Event)

		if event.Type == Intersection {
			if _, found := processedIntersections[event.Point]; found {
				continue
			}
			processedIntersections[event.Point] = true
		}

		status.SetX(event.Point.X)

		switch event.Type {
		case SegmentStart:
			seg := event.Segments[0]
			status.Add(seg)
			above, below := status.FindNeighbors(seg)
			if above != nil {
				checkIntersection(seg, above, event.Point, &eq)
			}
			if below != nil {
				checkIntersection(seg, below, event.Point, &eq)
			}
		case SegmentEnd:
			seg := event.Segments[0]
			above, below := status.FindNeighbors(seg)
			status.Remove(seg)
			if above != nil && below != nil {
				checkIntersection(above, below, event.Point, &eq)
			}
		case Intersection:
			intersections++
			s1, s2 := event.Segments[0], event.Segments[1]

			// Determine which segment is upper and which is lower *before* the intersection point.
			var lowerSeg, upperSeg *Segment
			if status.comparator.Compare(s1, s2) < 0 {
				lowerSeg, upperSeg = s1, s2
			} else {
				lowerSeg, upperSeg = s2, s1
			}

			// CORRECTED: Find outer neighbors with separate, valid function calls.
			aboveUpper, _ := status.FindNeighbors(upperSeg)
			_, belowLower := status.FindNeighbors(lowerSeg)

			// At the intersection point, the two segments swap their vertical order.
			// We achieve this by removing and re-adding them. The comparator
			// will place them in their new, correct positions.
			status.Remove(lowerSeg)
			status.Remove(upperSeg)
			status.Add(lowerSeg)
			status.Add(upperSeg)

			// CORRECTED: Check for new intersections between the swapped segments
			// and their former outer neighbors.
			if aboveUpper != nil {
				checkIntersection(lowerSeg, aboveUpper, event.Point, &eq)
			}
			if belowLower != nil {
				checkIntersection(upperSeg, belowLower, event.Point, &eq)
			}
		}
	}

	return intersections
}

// checkIntersection (No changes)
func checkIntersection(s1, s2 *Segment, currentPoint Point, eq *EventQueue) {
	if p, ok := s1.intersection(*s2); ok {
		if p.X > currentPoint.X || (math.Abs(p.X-currentPoint.X) < epsilon && p.Y > currentPoint.Y) {
			heap.Push(eq, &Event{Point: p, Type: Intersection, Segments: []*Segment{s1, s2}})
		}
	}
}
