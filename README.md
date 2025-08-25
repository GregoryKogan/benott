# benott

[![Go CI](https://github.com/GregoryKogan/benott/actions/workflows/go.yml/badge.svg)](https://github.com/GregoryKogan/benott/actions)
[![codecov](https://codecov.io/gh/GregoryKogan/benott/graph/badge.svg)](https://codecov.io/gh/GregoryKogan/benott)
[![Go Report Card](https://goreportcard.com/badge/github.com/GregoryKogan/benott)](https://goreportcard.com/report/github.com/GregoryKogan/benott)
[![Go Reference](https://pkg.go.dev/badge/github.com/GregoryKogan/benott.svg)](https://pkg.go.dev/github.com/GregoryKogan/benott)

A blazingly fast, production-ready Go implementation of the Bentley-Ottmann algorithm for counting line segment intersections.

## Overview

This library provides an efficient and robust solution for a classic computational geometry problem: finding the number of intersections in a set of 2D line segments. It uses a sweep-line algorithm powered by a Red-Black Tree for its status structure, ensuring optimal performance.

The implementation is carefully designed to handle complex edge cases, including vertical segments and multiple segments intersecting at a single point, making it suitable for demanding, production-level applications.

## Features

- **High Performance**: Achieves the optimal `O((n+k) log n)` time complexity, where `n` is the number of segments and `k` is the number of intersections.
- **Robust and Accurate**: Correctly handles edge cases like vertical lines, collinear points, and multiple segments intersecting at the same point.
- **Extensively Tested**: Near-perfect test coverage ensures reliability and correctness.
- **Simple API**: A single, clear entry point (`CountIntersections`) makes the library easy to integrate.
- **Memory Efficient**: Predictable, linear memory scaling with the number of input segments.

## Installation

```sh
go get github.com/GregoryKogan/benott
```

## Usage

Here is a simple example of how to use the library to count intersections.

```go
package main

import (
 "fmt"
 "github.com/GregoryKogan/benott"
)

func main() {
 // Define a set of line segments.
 segments := []benott.Segment{
  {P1: benott.Point{X: 0, Y: 0}, P2: benott.Point{X: 10, Y: 10}},
  {P1: benott.Point{X: 0, Y: 10}, P2: benott.Point{X: 10, Y: 0}},
  {P1: benott.Point{X: 5, Y: 0}, P2: benott.Point{X: 5, Y: 10}},
 }

 // Count the number of intersections.
 // For this "star" pattern, the 3 segments intersect at one point,
 // creating 3 unique intersection pairs: (s1,s2), (s1,s3), (s2,s3).
 intersections := benott.CountIntersections(segments)

 fmt.Printf("Found %d intersections.\n", intersections)
 // Output: Found 3 intersections.
}
```

## Performance

Benchmarks confirm the library's theoretical time complexity. The tests were run on an **Apple M1 Pro**.

The `BenchmarkRandomSegments` suite highlights the `O(n log n)` scaling with a sparse number of intersections. The `BenchmarkGridSegments` suite creates a high-contention scenario to highlight the `O(k log n)` scaling, where `k` is the dominant factor.

| Benchmark                                                        | Operations | Time/Op      | Memory/Op    | Allocs/Op  |
| ---------------------------------------------------------------- | ---------- | ------------ | ------------ | ---------- |
| `BenchmarkRandomSegments/N=10`                                   | 234319     | 4372 ns/op   | 3496 B/op    | 85 allocs/op   |
| `BenchmarkRandomSegments/N=100`                                  | 9804       | 115627 ns/op | 44056 B/op   | 1030 allocs/op |
| `BenchmarkRandomSegments/N=1000`                                 | 645        | 1851251 ns/op| 433798 B/op  | 9852 allocs/op |
| `BenchmarkRandomSegments/N=10000`                                | 43         | 29025479 ns/op| 4694687 B/op | 100505 allocs/op|
| `BenchmarkGridSegments/Grid=10x10_Segments=20_Intersections=100`   | 21450      | 55818 ns/op  | 31624 B/op   | 811 allocs/op  |
| `BenchmarkGridSegments/Grid=50x50_Segments=100_Intersections=2500` | 745        | 1604035 ns/op| 685324 B/op  | 18013 allocs/op|
| `BenchmarkGridSegments/Grid=100x100_Segments=200_Intersections=10000`| 165        | 7506986 ns/op| 2691138 B/op | 71014 allocs/op|
| `BenchmarkGridSegments/Grid=200x200_Segments=400_Intersections=40000`| 36         | 31361110 ns/op| 10661605 B/op| 282015 allocs/op|

The results show excellent, predictable scaling in line with the algorithm's optimal complexity.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for bugs, feature requests, or suggestions.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

This project also includes third-party dependencies. Their licenses are collected in the [LICENSES-3RD-PARTY.md](LICENSES-3RD-PARTY.md) file.
