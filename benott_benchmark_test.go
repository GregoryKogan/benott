package benott_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/GregoryKogan/benott"
)

// --- Helper Functions for Benchmark Data Generation ---

// generateRandomSegments creates n segments with random coordinates.
// This scenario typically has a low to moderate number of intersections.
func generateRandomSegments(n int, maxCoord float64) []benott.Segment {
	segments := make([]benott.Segment, n)
	// Use a fixed seed for reproducible benchmarks
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range n {
		segments[i] = benott.Segment{
			P1: benott.Point{X: rng.Float64() * maxCoord, Y: rng.Float64() * maxCoord},
			P2: benott.Point{X: rng.Float64() * maxCoord, Y: rng.Float64() * maxCoord},
		}
	}
	return segments
}

// generateGridSegments creates a grid of n horizontal and n vertical lines.
// This scenario is designed to generate a very high number of intersections (n*n).
// Total segments created will be 2*n.
func generateGridSegments(n int, maxCoord float64) []benott.Segment {
	segments := make([]benott.Segment, 2*n)
	step := maxCoord / float64(n+1)

	// Create n horizontal lines
	for i := range n {
		y := step * float64(i+1)
		segments[i] = benott.Segment{
			P1: benott.Point{X: 0, Y: y},
			P2: benott.Point{X: maxCoord, Y: y},
		}
	}

	// Create n vertical lines
	for i := range n {
		x := step * float64(i+1)
		segments[n+i] = benott.Segment{
			P1: benott.Point{X: x, Y: 0},
			P2: benott.Point{X: x, Y: maxCoord},
		}
	}
	return segments
}

// --- Benchmark Suites ---

// BenchmarkRandomSegments tests performance against N with a sparse number of intersections.
// This should highlight the O(n log n) characteristic.
func BenchmarkRandomSegments(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, n := range sizes {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			// Generate the test data once, outside the timed loop.
			segments := generateRandomSegments(n, 1000.0)
			b.ResetTimer() // Exclude generation time from the benchmark.

			for b.Loop() {
				// The function being benchmarked.
				benott.CountIntersections(segments)
			}
		})
	}
}

// BenchmarkGridSegments tests performance against N with a dense number of intersections (k).
// This should highlight the O((n+k) log n) characteristic, where k is dominant.
func BenchmarkGridSegments(b *testing.B) {
	// gridSizes defines the number of lines in one direction (e.g., 10x10, 50x50).
	// Total segments will be 2 * gridSize.
	gridSizes := []int{10, 50, 100, 200}

	for _, size := range gridSizes {
		numSegments := 2 * size
		numIntersections := size * size
		b.Run(fmt.Sprintf("Grid=%dx%d_Segments=%d_Intersections=%d", size, size, numSegments, numIntersections), func(b *testing.B) {
			segments := generateGridSegments(size, 1000.0)
			b.ResetTimer()

			for b.Loop() {
				benott.CountIntersections(segments)
			}
		})
	}
}
