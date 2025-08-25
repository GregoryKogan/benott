package benott

import (
	"container/heap"
	"math"
	"sort"
	"sync"
)

// This pool allows us to reuse Event structs instead of allocating new ones on the
// heap for every endpoint and intersection. This dramatically reduces GC pressure.
var eventPool = sync.Pool{
	New: func() any {
		return &Event{}
	},
}

// CountIntersections implements the Bentley-Ottmann algorithm to find the total
// number of intersection points in a given set of line segments.
//
// The algorithm uses a sweep-line approach, processing events (segment endpoints
// and intersections) from left to right. It maintains a status structure (a
// Red-Black Tree) of segments that cross the sweep line, ordered vertically.
// This allows it to efficiently find new intersections only between adjacent
// segments.
//
// The time complexity is O((n+k) log n), where n is the number of segments and
// k is the number of intersections. The space complexity is O(n+k).
//
// It correctly handles complex cases, including vertical segments and multiple
// segments intersecting at a single point.
func CountIntersections(segments []Segment) int {
	// The event queue stores all segment endpoints to initialize the sweep.
	// Pre-allocate the event queue with a known initial size.
	// Each segment generates two initial events (start and end).
	eq := make(EventQueue, 0, len(segments)*2)
	segmentCopies := make([]Segment, len(segments))
	copy(segmentCopies, segments)

	// This single loop correctly normalizes, pre-computes, and creates events
	// for each segment in a logical, efficient order.
	for i := range segmentCopies {
		s := &segmentCopies[i] // Use a pointer to modify the copy

		// 1. NORMALIZE FIRST: Ensure P1 is always the leftmost endpoint.
		if s.P1.X > s.P2.X || (s.P1.X == s.P2.X && s.P1.Y > s.P2.Y) {
			s.P1, s.P2 = s.P2, s.P1
		}

		// 2. PRE-COMPUTE SECOND: Calculate properties based on the final, normalized points.
		p1, p2 := s.P1, s.P2 // Use the now-normalized points
		if math.Abs(p1.X-p2.X) < epsilon {
			s.isVertical = true
			s.slope = math.Inf(1)
			s.yIntercept = math.NaN() // Not applicable
		} else {
			s.isVertical = false // Ensure this is set correctly
			s.slope = (p2.Y - p1.Y) / (p2.X - p1.X)
			s.yIntercept = p1.Y - s.slope*p1.X
		}

		// 3. PUSH EVENTS THIRD: Create and push the start and end events.
		startEvent := eventPool.Get().(*Event)
		startEvent.Point = s.P1
		startEvent.Type = SegmentStart
		startEvent.Seg1 = s
		heap.Push(&eq, startEvent)

		endEvent := eventPool.Get().(*Event)
		endEvent.Point = s.P2
		endEvent.Type = SegmentEnd
		endEvent.Seg1 = s
		heap.Push(&eq, endEvent)
	}

	status := NewStatus()
	intersections := 0

	// Pre-allocate the slice used for sorting intersecting segments.
	// We declare it once outside the loop and reset its length to 0 on each use.
	// This avoids re-allocating this slice in every intersection event.
	segsAsSlice := make([]*Segment, 0, 16) // Start with a reasonable capacity

	for eq.Len() > 0 {
		event := heap.Pop(&eq).(*Event)

		// Set the sweep line's current position for all status comparisons.
		status.SetX(event.Point.X)

		switch event.Type {
		case SegmentStart:
			// A new segment is added to the status.
			seg := event.Seg1
			status.Add(seg)
			// Check for intersections with its new neighbors.
			above, below := status.FindNeighbors(seg)
			if above != nil {
				checkIntersection(seg, above, event.Point, &eq)
			}
			if below != nil {
				checkIntersection(seg, below, event.Point, &eq)
			}

		case SegmentEnd:
			// A segment is removed from the status.
			seg := event.Seg1
			above, below := status.FindNeighbors(seg)
			status.Remove(seg)
			// Its former neighbors are now adjacent; check if they intersect.
			if above != nil && below != nil {
				checkIntersection(above, below, event.Point, &eq)
			}

		case Intersection:
			// This block handles the complex case of one or more intersections
			// occurring at the exact same point.

			// 1. Identify all segments involved in intersections at this point.
			intersectingSegs := make(map[*Segment]bool)
			intersectingSegs[event.Seg1] = true
			intersectingSegs[event.Seg2] = true

			// Peek at subsequent events in the queue to find all other
			// intersections happening at this exact coordinate.
			for eq.Len() > 0 && (*eq[0]).Point == event.Point && (*eq[0]).Type == Intersection {
				nextEvent := heap.Pop(&eq).(*Event)
				intersectingSegs[nextEvent.Seg1] = true
				intersectingSegs[nextEvent.Seg2] = true

				// When returning to pool, nil out pointers to prevent memory leaks.
				nextEvent.Seg1 = nil
				nextEvent.Seg2 = nil
				eventPool.Put(nextEvent)
			}

			// 2. Count all intersections. If k segments meet at one point,
			// this creates C(k, 2) = k*(k-1)/2 intersections.
			k := len(intersectingSegs)
			intersections += k * (k - 1) / 2

			// 3. Find the neighbors of the entire block of intersecting segments
			// before their order is changed.
			// Reset the pre-allocated slice instead of making a new one.
			segsAsSlice = segsAsSlice[:0]
			for seg := range intersectingSegs {
				segsAsSlice = append(segsAsSlice, seg)
			}
			sort.Slice(segsAsSlice, func(i, j int) bool {
				return status.comparator.Compare(segsAsSlice[i], segsAsSlice[j]) < 0
			})

			bottomSeg, topSeg := segsAsSlice[0], segsAsSlice[k-1]
			aboveTop, _ := status.FindNeighbors(topSeg)
			_, belowBottom := status.FindNeighbors(bottomSeg)

			// 4. Update the status by removing and re-inserting all intersecting
			// segments. This correctly reorders them after the intersection point.
			for _, seg := range segsAsSlice {
				status.Remove(seg)
			}
			status.SetX(event.Point.X + epsilon) // Use a point slightly after to get the new order
			for i := k - 1; i >= 0; i-- {        // Re-add in their new (reversed) order
				status.Add(segsAsSlice[i])
			}

			// 5. Check for new intersections between the block's new boundaries
			// and their old outer neighbors.
			if aboveTop != nil {
				checkIntersection(segsAsSlice[0], aboveTop, event.Point, &eq) // New top vs. old neighbor
			}
			if belowBottom != nil {
				checkIntersection(segsAsSlice[k-1], belowBottom, event.Point, &eq) // New bottom vs. old neighbor
			}
		}

		// Return the processed event to the pool.
		event.Seg1 = nil
		event.Seg2 = nil
		eventPool.Put(event)
	}

	return intersections
}

// checkIntersection checks if two segments s1 and s2 intersect at a point that
// is to the right of the current sweep line. If they do, a new Intersection
// event is created and pushed onto the event queue.
func checkIntersection(s1, s2 *Segment, currentPoint Point, eq *EventQueue) {
	if s1 == nil || s2 == nil {
		return
	}
	if p, ok := s1.intersection(*s2); ok {
		// Only add events that are in the future (to the right of the sweep line,
		// or on the line but with a greater Y value). This prevents adding
		// duplicate events or events that have already been processed.
		// This condition is critical to prevent infinite loops from floating-point errors.
		// A new intersection event is only added if its X-coordinate is *meaningfully*
		// to the right of the current event's X-coordinate, or if it's on the same
		// vertical line but *meaningfully* above the current event point.
		isFutureEvent := (p.X-currentPoint.X > epsilon) ||
			(math.Abs(p.X-currentPoint.X) < epsilon && p.Y-currentPoint.Y > epsilon)

		if isFutureEvent {
			// Get event from the pool.
			newEvent := eventPool.Get().(*Event)
			newEvent.Point = p
			newEvent.Type = Intersection
			newEvent.Seg1 = s1
			newEvent.Seg2 = s2
			heap.Push(eq, newEvent)
		}
	}
}
