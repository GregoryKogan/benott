package benott

// EventType defines the nature of an event in the sweep-line algorithm.
type EventType int

const (
	// SegmentStart signifies that the sweep line has reached the left endpoint of a segment.
	SegmentStart EventType = iota
	// SegmentEnd signifies that the sweep line has reached the right endpoint of a segment.
	SegmentEnd
	// Intersection signifies that the sweep line has reached a point where two or more segments cross.
	Intersection
)

// Event represents a significant point in the 2D plane to be processed by the
// Bentley-Ottmann algorithm. Each event has a location (Point), a Type, and one
// or more associated Segments.
type Event struct {
	// Point is the coordinate in the 2D plane where the event occurs.
	Point Point
	// Type is the nature of the event (e.g., SegmentStart, Intersection).
	Type EventType
	// Seg1 is the primary segment associated with this event.
	Seg1 *Segment
	// Seg2 is the second segment, used only for Intersection events.
	Seg2 *Segment
}

// EventQueue is a min-priority queue of events, implemented using Go's container/heap.
// Events are ordered primarily by their X-coordinate, then by their Y-coordinate
// as a tie-breaker. This ensures the sweep-line processes points from left-to-right,
// bottom-to-top.
type EventQueue []*Event

// Len returns the number of events in the queue.
func (eq EventQueue) Len() int { return len(eq) }

// Less reports whether the event at index i should be sorted before the event at index j.
func (eq EventQueue) Less(i, j int) bool {
	if eq[i].Point.X != eq[j].Point.X {
		return eq[i].Point.X < eq[j].Point.X
	}
	return eq[i].Point.Y < eq[j].Point.Y
}

// Swap swaps the events at indices i and j.
func (eq EventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

// Push adds an event to the priority queue.
func (eq *EventQueue) Push(x any) {
	*eq = append(*eq, x.(*Event))
}

// Pop removes and returns the event with the lowest value (highest priority) from the queue.
func (eq *EventQueue) Pop() any {
	old := *eq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*eq = old[0 : n-1]
	return item
}
