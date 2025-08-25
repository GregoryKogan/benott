package benott

import (
	"math"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type sweepLineComparator struct {
	currentX float64
}

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
	return seg.P1.Y + (c.currentX-seg.P1.X)*(seg.P2.Y-seg.P1.Y)/(seg.P2.X-seg.P1.X)
}

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

	// Tie-breaking using slope for co-linear segments at the event point.
	// This ensures a consistent, stable ordering and handles vertical lines.
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

type Status struct {
	tree       *rbt.Tree
	comparator *sweepLineComparator
}

func NewStatus() *Status {
	comp := &sweepLineComparator{}
	return &Status{
		tree:       rbt.NewWith(comp.Compare),
		comparator: comp,
	}
}
func (s *Status) SetX(x float64)      { s.comparator.currentX = x }
func (s *Status) Add(seg *Segment)    { s.tree.Put(seg, true) }
func (s *Status) Remove(seg *Segment) { s.tree.Remove(seg) }

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
