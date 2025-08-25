package benott

import (
	"container/heap"
	"math"
	"sort"
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

	for eq.Len() > 0 {
		event := heap.Pop(&eq).(*Event)

		// Set the sweep line position for all comparisons related to this event.
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
			// The simple pairwise swap is insufficient for multi-line intersections.
			// The robust approach is to re-evaluate the local order.

			// 1. Identify all segments involved in intersections at this specific point.
			intersectingSegs := make(map[*Segment]bool)
			intersectingSegs[event.Segments[0]] = true
			intersectingSegs[event.Segments[1]] = true

			// Peek at next events in the queue to find all intersections at this exact point.
			for eq.Len() > 0 && (*eq[0]).Point == event.Point && (*eq[0]).Type == Intersection {
				nextEvent := heap.Pop(&eq).(*Event)
				intersectingSegs[nextEvent.Segments[0]] = true
				intersectingSegs[nextEvent.Segments[1]] = true
			}

			// 2. Count all intersections at this point.
			// n segments intersecting at one point produce n*(n-1)/2 intersections.
			k := len(intersectingSegs)
			intersections += k * (k - 1) / 2

			// 3. Find the neighbors of the entire block of intersecting segments.
			segsAsSlice := make([]*Segment, 0, k)
			for seg := range intersectingSegs {
				segsAsSlice = append(segsAsSlice, seg)
			}
			sort.Slice(segsAsSlice, func(i, j int) bool {
				return status.comparator.Compare(segsAsSlice[i], segsAsSlice[j]) < 0
			})

			bottomSeg, topSeg := segsAsSlice[0], segsAsSlice[k-1]
			aboveTop, _ := status.FindNeighbors(topSeg)
			_, belowBottom := status.FindNeighbors(bottomSeg)

			// 4. Remove and re-insert all intersecting segments to update their order.
			for _, seg := range segsAsSlice {
				status.Remove(seg)
			}
			status.SetX(event.Point.X + epsilon) // Use a point slightly after to get the new order
			for i := k - 1; i >= 0; i-- {        // Re-add in reverse order
				status.Add(segsAsSlice[i])
			}

			// 5. Check for new intersections between the block's new boundaries and old neighbors.
			if aboveTop != nil {
				checkIntersection(segsAsSlice[0], aboveTop, event.Point, &eq) // new top is old bottom
			}
			if belowBottom != nil {
				checkIntersection(segsAsSlice[k-1], belowBottom, event.Point, &eq) // new bottom is old top
			}
		}
	}

	return intersections
}

// checkIntersection determines if two segments intersect at a valid future point
// and, if so, adds an intersection event to the queue.
func checkIntersection(s1, s2 *Segment, currentPoint Point, eq *EventQueue) {
	if s1 == nil || s2 == nil {
		return
	}
	if p, ok := s1.intersection(*s2); ok {
		// Only add events that are strictly to the right of the current sweep line,
		// or on the sweep line but "after" the current event point.
		if p.X > currentPoint.X || (math.Abs(p.X-currentPoint.X) < epsilon && p.Y > currentPoint.Y) {
			heap.Push(eq, &Event{Point: p, Type: Intersection, Segments: []*Segment{s1, s2}})
		}
	}
}
