package util

import (
	"fmt"
	"math"
)

// Grid
type Grid[T any] struct {
	Width  int
	Height int
	Points []T
}

func (g *Grid[T]) At(x, y int) *T {
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return nil
	}
	index := y*g.Width + x
	if index < 0 || index >= len(g.Points) {
		return nil
	}
	return &g.Points[index]
}

func (g *Grid[T]) Set(x, y int, val T) error {
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return fmt.Errorf("coordinates out of bounds")
	}
	index := y*g.Width + x
	if index < 0 || index >= len(g.Points) {
		return fmt.Errorf("coordinates out of bounds")
	}
	g.Points[index] = val
	return nil
}

func (g *Grid[T]) AddPoint(x, y int, val T) error {
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return fmt.Errorf("coordinates out of bounds")
	}
	index := y*g.Width + x
	if index >= len(g.Points) {
		g.Points = append(g.Points, val)
		return nil
	}
	return fmt.Errorf("coordinates already occupied")
}

func (g *Grid[T]) Print() {
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			fmt.Print(*g.At(x, y))
		}
		fmt.Println()
	}
}

func (g *Grid[T]) Neighbors(x, y int) struct {
	North *T
	South *T
	East  *T
	West  *T
} {
	return struct {
		North *T
		South *T
		East  *T
		West  *T
	}{
		North: g.At(x, y-1),
		South: g.At(x, y+1),
		East:  g.At(x+1, y),
		West:  g.At(x-1, y),
	}
}

func (g *Grid[T]) ManhattanDistance(x1 int, y1 int, x2 int, y2 int) int {
	return int(math.Abs(float64(x1-x2))) + int(math.Abs(float64(y1-y2)))
}
