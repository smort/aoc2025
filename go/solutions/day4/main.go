package main

import (
	"fmt"

	"github.com/smort/aoc2025/util"
)

func main() {
	part1("input.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0

	lines := util.GetLines(filename)
	grid := makeGrid(lines)

	visitMovableRolls(grid, func(coord util.Coordinate) {
		result++
	})

	fmt.Println(result)
}

func part2(filename string) {
	result := 0

	lines := util.GetLines(filename)
	grid := makeGrid(lines)
	lastMoved := -1
	for lastMoved != 0 {
		lastMoved = 0
		visitMovableRolls(grid, func(coord util.Coordinate) {
			result++
			lastMoved++
			grid.SetCell(coord, 'X')
		})
	}

	fmt.Println(result)
}

func makeGrid(lines []string) *util.DenseGrid {
	grid := make([][]rune, 0)
	for _, line := range lines {
		grid = append(grid, []rune(line))
	}
	denseGrid := util.NewDenseGrid(len(grid[0]), len(grid), grid)
	denseGrid.Directions = util.Directions8
	return denseGrid
}

func visitMovableRolls(grid *util.DenseGrid, fn func(util.Coordinate)) {
	rolls := grid.FindCoordinates('@')
	for _, r := range rolls {
		neighbors := grid.GetNeighbors(r)
		adjacentRolls := 0
		for _, n := range neighbors {
			if *grid.At(n) == '@' {
				adjacentRolls++
			}
			if adjacentRolls >= 4 {
				break
			}
		}
		if adjacentRolls < 4 {
			fn(r)
		}
	}
}
