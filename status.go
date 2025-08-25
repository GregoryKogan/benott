package benott

import (
	"math"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

// sweepLineComparator provides the dynamic comparison logic for the Red-Black Tree.
// The ordering of segments in the sweep-line status depends on their y-coordinate
// at the current x-position of the sweep line. This struct holds that `currentX`
// state, allowing the comparator to function correctly at each event point.
type sweepLineComparator struct {
	currentX float64
}

// getY calculates the y-coordinate of a segment at the comparator's currentX.
func (c *sweepLineComparator) getY(seg *Segment) float64 {
	if seg.P1.X == seg.P2.X {
		return seg.P1.Y
	}
	if c.currentX <= seg.P1.X {
		return seg.P1.Y
	}
	if c.currentX >= seg.P2.X {
		return seg.P2.Y
	}
	// Linear interpolation: y = y1 + (x - x1) * (y2 - y1) / (x2 - x1)
	return seg.P1.Y + (c.currentX-seg.P1.X)*(seg.P2.Y-seg.P1.Y)/(seg.P2.X-seg.P1.X)
}

// Compare implements the github.com/emirpasic/gods/utils.Comparator interface.
// It compares two segments based on their y-coordinates at the current sweep-line
// position. If y-coordinates are equal, it uses the segment's slope as a tie-breaker
// to ensure a consistent and stable ordering.
func (c *sweepLineComparator) Compare(a, b interface{}) int {
	segA := a.(*Segment)
	segB := b.(*Segment)
	yA := c.getY(segA)
	yB := c.getY(segB)

	if math.Abs(yA-yB) > epsilon {
		if yA < yB {
			return -1
		}
		return 1
	}

	// Tie-breaking with slope handles collinear segments and ensures stable ordering.
	var slopeA, slopeB float64
	if math.Abs(segA.P2.X-segA.P1.X) < epsilon {
		slopeA = math.Inf(1) // Vertical slope is infinite
	} else {
		slopeA = (segA.P2.Y - segA.P1.Y) / (segA.P2.X - segA.P1.X)
	}
	if math.Abs(segB.P2.X-segB.P1.X) < epsilon {
		slopeB = math.Inf(1)
	} else {
		slopeB = (segB.P2.Y - segB.P1.Y) / (segB.P2.X - segB.P1.X)
	}

	if slopeA < slopeB {
		return -1
	}
	if slopeA > slopeB {
		return 1
	}
	return 0
}

// Status represents the sweep-line status structure. It maintains the set of
// segments that are currently intersecting the sweep line, ordered vertically.
//
// It is implemented using a Red-Black Tree to achieve efficient O(log n)
// for add, remove, and neighbor-finding operations.
type Status struct {
	tree       *rbt.Tree
	comparator *sweepLineComparator
}

// NewStatus creates and initializes a new Status structure.
func NewStatus() *Status {
	comp := &sweepLineComparator{}
	return &Status{
		tree:       rbt.NewWith(comp.Compare),
		comparator: comp,
	}
}

// SetX updates the current x-coordinate of the sweep line for the status comparator.
// This is a critical step and MUST be called before any tree operations at a new
// event point to ensure segments are compared correctly.
func (s *Status) SetX(x float64) { s.comparator.currentX = x }

// Add inserts a segment into the status tree.
func (s *Status) Add(seg *Segment) { s.tree.Put(seg, true) }

// Remove deletes a segment from the status tree.
func (s *Status) Remove(seg *Segment) { s.tree.Remove(seg) }

// findSuccessor finds the in-order successor of a node in the tree (the next largest element).
func findSuccessor(node *rbt.Node) *rbt.Node {
	if node.Right != nil {
		curr := node.Right
		for curr.Left != nil {
			curr = curr.Left
		}
		return curr
	}
	p := node.Parent
	curr := node
	for p != nil && curr == p.Right {
		curr = p
		p = p.Parent
	}
	return p
}

// findPredecessor finds the in-order predecessor of a node in the tree (the next smallest element).
func findPredecessor(node *rbt.Node) *rbt.Node {
	if node.Left != nil {
		curr := node.Left
		for curr.Right != nil {
			curr = curr.Right
		}
		return curr
	}
	p := node.Parent
	curr := node
	for p != nil && curr == p.Left {
		curr = p
		p = p.Parent
	}
	return p
}

// FindNeighbors finds the segments immediately above and below a given segment
// in the status tree. It returns `nil` for a neighbor if one does not exist.
func (s *Status) FindNeighbors(seg *Segment) (above, below *Segment) {
	node := s.tree.GetNode(seg)
	if node == nil {
		return nil, nil
	}
	if predNode := findPredecessor(node); predNode != nil {
		below = predNode.Key.(*Segment)
	}
	if succNode := findSuccessor(node); succNode != nil {
		above = succNode.Key.(*Segment)
	}
	return above, below
}
