package main

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/smort/aoc2025/util"
)

func main() {
	part1("example.txt")
	part1("input.txt")

	part2("example.txt")
	part2("input.txt")
}

type Point struct {
	X, Y int
}

func part1(filename string) {
	// result := 0
	lines := util.GetLines(filename)

	points := make([]Point, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ",")
		points = append(points, Point{X: util.MustConvAtoi(parts[0]), Y: util.MustConvAtoi(parts[1])})
	}

	maxSize := 0
	winningPoints := [2]Point{}
	for _, point := range points {
		for _, other := range points {
			if point == other {
				continue
			}

			x := math.Abs(float64(point.X-other.X)) + 1
			y := math.Abs(float64(point.Y-other.Y)) + 1
			area := int(x * y)
			if area > maxSize {
				maxSize = area
				winningPoints[0] = point
				winningPoints[1] = other
			}
		}
	}

	fmt.Println(maxSize)
	fmt.Println(winningPoints)
}

func part2(filename string) {
	lines := util.GetLines(filename)

	points := make([]Point, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ",")
		points = append(points, Point{X: util.MustConvAtoi(parts[0]), Y: util.MustConvAtoi(parts[1])})
	}

	edgePoints := collectAllEdgePoints(points)
	maxArea := 0
	for _, point := range points {
		for _, other := range points {
			if point == other {
				continue
			}

			x := math.Abs(float64(point.X-other.X)) + 1
			y := math.Abs(float64(point.Y-other.Y)) + 1
			area := int(x * y)
			if area <= maxArea {
				continue // skip ones we've already beaten
			}

			// check if any points are inside. if so we've intersected the shape
			isInvalid := false
			for _, e := range edgePoints {
				if isIn(e, point, other) {
					isInvalid = true
					break
				}
			}
			if isInvalid {
				continue
			}

			maxArea = area
		}
	}

	fmt.Println(maxArea)
}

func isIn(p, corner1, corner2 Point) bool {
	minX := min(corner1.X, corner2.X)
	maxX := max(corner1.X, corner2.X)
	minY := min(corner1.Y, corner2.Y)
	maxY := max(corner1.Y, corner2.Y)

	return p.X > minX && p.X < maxX && p.Y > minY && p.Y < maxY
}

func collectAllEdgePoints(points []Point) []Point {
	edgePoints := make([]Point, 0, len(points))
	for i := range len(points) - 1 {
		newEdgePoints := getEdgePoints(points[i], points[i+1])
		edgePoints = slices.Concat(edgePoints, newEdgePoints)
	}

	newEdgePoints := getEdgePoints(points[len(points)-1], points[0]) // connect first and last
	edgePoints = slices.Concat(edgePoints, newEdgePoints)

	return edgePoints
}

func getEdgePoints(point1 Point, point2 Point) []Point {
	points := make([]Point, 0)
	if point1.X == point2.X { // vertical line
		minY := min(point1.Y, point2.Y)
		maxY := max(point1.Y, point2.Y)
		for y := minY; y < maxY; y++ {
			points = append(points, Point{X: point1.X, Y: y})
		}
	} else if point1.Y == point2.Y { // horizontal line
		minX := min(point1.X, point2.X)
		maxX := max(point1.X, point2.X)
		for x := minX; x < maxX; x++ {
			points = append(points, Point{X: x, Y: point1.Y})
		}
	}

	return points
}
