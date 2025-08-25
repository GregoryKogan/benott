package benott_test

import (
	"testing"

	"github.com/GregoryKogan/benott"
)

func TestCountIntersections(t *testing.T) {
	testCases := []struct {
		name     string
		segments []benott.Segment
		expected int
	}{
		{
			name: "Simple single intersection",
			segments: []benott.Segment{
				{P1: benott.Point{X: 0, Y: 0}, P2: benott.Point{X: 10, Y: 10}},
				{P1: benott.Point{X: 0, Y: 10}, P2: benott.Point{X: 10, Y: 0}},
			},
			expected: 1,
		},
		{
			name: "No intersection",
			segments: []benott.Segment{
				{P1: benott.Point{X: 0, Y: 0}, P2: benott.Point{X: 10, Y: 10}},
				{P1: benott.Point{X: 0, Y: 1}, P2: benott.Point{X: 10, Y: 11}}, // Parallel
			},
			expected: 0,
		},
		{
			name: "Three segments, one intersection point",
			segments: []benott.Segment{
				{P1: benott.Point{X: 0, Y: 0}, P2: benott.Point{X: 10, Y: 10}},
				{P1: benott.Point{X: 0, Y: 10}, P2: benott.Point{X: 10, Y: 0}},
				{P1: benott.Point{X: 5, Y: 0}, P2: benott.Point{X: 5, Y: 10}}, // Vertical
			},
			expected: 1, // The vertical line intersects the other two
		},
		{
			name: "Multiple separate intersections",
			segments: []benott.Segment{
				{P1: benott.Point{X: 0, Y: 0}, P2: benott.Point{X: 10, Y: 10}},
				{P1: benott.Point{X: 0, Y: 10}, P2: benott.Point{X: 10, Y: 0}},
				{P1: benott.Point{X: 1, Y: 2}, P2: benott.Point{X: 5, Y: 0}},
			},
			expected: 2,
		},
		{
			name: "Segments sharing an endpoint, no crossing",
			segments: []benott.Segment{
				{P1: benott.Point{X: 0, Y: 0}, P2: benott.Point{X: 5, Y: 5}},
				{P1: benott.Point{X: 5, Y: 5}, P2: benott.Point{X: 10, Y: 0}},
			},
			expected: 0,
		},
		{
			name:     "Empty set",
			segments: []benott.Segment{},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := benott.CountIntersections(tc.segments)
			if actual != tc.expected {
				t.Errorf("Expected %d intersections, but got %d", tc.expected, actual)
			}
		})
	}
}
