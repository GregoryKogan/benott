package benott_test

import (
	"testing"

	"github.com/GregoryKogan/benott"
)

// helper function to check results
func check(t *testing.T, segments []benott.Segment, expected int) {
	// t.Helper() marks this as a test helper.
	// When t.Errorf is called, the line number reported
	// will be from the calling function, not from inside check().
	t.Helper()

	actual := benott.CountIntersections(segments)
	if actual != expected {
		t.Errorf("Expected %d intersections, but got %d", expected, actual)
	}
}

// --- Basic Cases ---

func TestSimpleSingleIntersection(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{X: 0, Y: 0}, P2: benott.Point{X: 10, Y: 10}},
		{P1: benott.Point{X: 0, Y: 10}, P2: benott.Point{X: 10, Y: 0}},
	}
	check(t, segments, 1)
}

func TestNoIntersection(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{0, 0}, P2: benott.Point{10, 10}},
		{P1: benott.Point{0, 1}, P2: benott.Point{10, 11}},
	}
	check(t, segments, 0)
}

func TestEmptySet(t *testing.T) {
	segments := []benott.Segment{}
	check(t, segments, 0)
}

func TestSingleSegment(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{0, 0}, P2: benott.Point{10, 10}},
	}
	check(t, segments, 0)
}

// --- Vertical and Horizontal Lines ---

func TestVerticalLineIntersection(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{5, 0}, P2: benott.Point{5, 10}},
		{P1: benott.Point{0, 5}, P2: benott.Point{10, 5}},
	}
	check(t, segments, 1)
}

func TestHorizontalLinesNoIntersection(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{0, 5}, P2: benott.Point{10, 5}},
		{P1: benott.Point{0, 6}, P2: benott.Point{10, 6}},
	}
	check(t, segments, 0)
}

func TestVerticalLinesNoIntersection(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{5, 0}, P2: benott.Point{5, 10}},
		{P1: benott.Point{6, 0}, P2: benott.Point{6, 10}},
	}
	check(t, segments, 0)
}

// --- Endpoint and Collinear Cases ---

func TestTJunctionIntersectionAtEndpoint(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{5, 0}, P2: benott.Point{5, 10}},
		{P1: benott.Point{0, 5}, P2: benott.Point{5, 5}},
	}
	check(t, segments, 1)
}

func TestVShapeNoIntersection(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{0, 0}, P2: benott.Point{5, 5}},
		{P1: benott.Point{10, 0}, P2: benott.Point{5, 5}},
	}
	check(t, segments, 0)
}

func TestCollinearNonOverlapping(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{0, 0}, P2: benott.Point{5, 5}},
		{P1: benott.Point{6, 6}, P2: benott.Point{10, 10}},
	}
	check(t, segments, 0)
}

func TestCollinearOverlapping(t *testing.T) {
	// Overlaps are not considered "intersections" in the crossing sense.
	segments := []benott.Segment{
		{P1: benott.Point{0, 0}, P2: benott.Point{10, 10}},
		{P1: benott.Point{2, 2}, P2: benott.Point{8, 8}},
	}
	check(t, segments, 0)
}

// --- Multiple Intersection Cases ---

func TestThreeLinesIntersectingAtOnePoint(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{5, 0}, P2: benott.Point{5, 10}},  // Vertical
		{P1: benott.Point{0, 5}, P2: benott.Point{10, 5}},  // Horizontal
		{P1: benott.Point{0, 0}, P2: benott.Point{10, 10}}, // Diagonal
	}
	// (v,h), (v,d), (h,d)
	check(t, segments, 3)
}

func TestFourLinesIntersectingAtOnePoint(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{5, 0}, P2: benott.Point{5, 10}},  // Vertical
		{P1: benott.Point{0, 5}, P2: benott.Point{10, 5}},  // Horizontal
		{P1: benott.Point{0, 0}, P2: benott.Point{10, 10}}, // Diagonal 1
		{P1: benott.Point{0, 10}, P2: benott.Point{10, 0}}, // Diagonal 2
	}
	// 4 choose 2 = 6
	check(t, segments, 6)
}

func TestSimple2x2Grid(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{0, 5}, P2: benott.Point{10, 5}},
		{P1: benott.Point{0, 6}, P2: benott.Point{10, 6}},
		{P1: benott.Point{5, 0}, P2: benott.Point{5, 10}},
		{P1: benott.Point{6, 0}, P2: benott.Point{6, 10}},
	}
	check(t, segments, 4)
}

func TestComplexCaseWithMultipleIntersections(t *testing.T) {
	segments := []benott.Segment{
		{P1: benott.Point{0, 0}, P2: benott.Point{10, 10}},
		{P1: benott.Point{0, 10}, P2: benott.Point{10, 0}},
		{P1: benott.Point{2, 0}, P2: benott.Point{8, 10}},
		{P1: benott.Point{0, 5}, P2: benott.Point{10, 5}},
	}
	check(t, segments, 6)
}
