package util

import "math"

func ManhattanDistance[T any](p0, p1 Point[T]) int {
	return int(math.Abs(float64(p0.Y-p1.Y)) + math.Abs(float64(p0.X-p1.X)))
}

type Directions[T any] struct {
	North Point[T]
	South Point[T]
	East  Point[T]
	West  Point[T]
}

func Neighbors[T any]() Directions[T] {
	return Directions[T]{
		West:  Point[T]{X: -1, Y: 0},
		North: Point[T]{X: 0, Y: -1},
		East:  Point[T]{X: 1, Y: 0},
		South: Point[T]{X: 0, Y: 1},
	}
}

type Point[T any] struct {
	X   int
	Y   int
	Val T
}

func (p Point[T]) Add(other Point[T]) Point[T] {
	return Point[T]{X: p.X + other.X, Y: p.Y + other.Y}
}

func (p Point[T]) Sub(other Point[T]) Point[T] {
	return Point[T]{X: p.X - other.X, Y: p.Y - other.Y}
}

func (p Point[T]) In(g Grid[any]) bool {
	return p.X > -1 && p.Y > -1 && p.X < g.Width && p.Y < g.Height
}
