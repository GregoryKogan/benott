package benott

// EventType defines the type of an event.
type EventType int

const (
	SegmentStart EventType = iota
	SegmentEnd
	Intersection
)

// Event represents an event in the sweep-line algorithm.
type Event struct {
	Point    Point
	Type     EventType
	Segments []*Segment // Holds 1 segment for start/end, 2 for intersection
}

// EventQueue is a priority queue of events, implemented using container/heap.
type EventQueue []*Event

func (eq EventQueue) Len() int { return len(eq) }

// Less compares two events. It sorts primarily by X-coordinate,
// and secondarily by Y-coordinate as a tie-breaker.
func (eq EventQueue) Less(i, j int) bool {
	if eq[i].Point.X != eq[j].Point.X {
		return eq[i].Point.X < eq[j].Point.X
	}
	return eq[i].Point.Y < eq[j].Point.Y
}

func (eq EventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

func (eq *EventQueue) Push(x any) {
	*eq = append(*eq, x.(*Event))
}

func (eq *EventQueue) Pop() any {
	old := *eq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*eq = old[0 : n-1]
	return item
}
