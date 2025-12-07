package main

import (
	"fmt"
	"time"

	"github.com/smort/aoc2025/util"
)

func main() {
	start := time.Now()
	part1("input.txt")
	fmt.Printf("part1 elapsed: %v\n", time.Since(start))

	start = time.Now()
	part2("input.txt")
	fmt.Printf("part2 elapsed: %v\n", time.Since(start))
}

type waterGrid struct {
	*util.DenseGrid
}

func (wg *waterGrid) GetNeighbors(pos util.Coordinate) []util.Coordinate {
	if wg.At(pos) == nil {
		return nil
	}

	val := wg.At(pos)
	if val != nil && *val == '^' {
		right := pos.Add(util.Coordinate{X: 1, Y: 0})
		left := pos.Add(util.Coordinate{X: -1, Y: 0})
		return []util.Coordinate{
			left,
			right,
		}
	}

	return []util.Coordinate{
		pos.Add(util.Coordinate{X: 0, Y: 1}),
	}
}

func part1(filename string) {
	lines := util.GetLines(filename)
	grid := util.NewDenseGridFromLines(lines)
	start := grid.FindCoordinates('S')[0]

	result := simulateWater(grid, start, map[util.Coordinate]struct{}{})
	fmt.Println(result)
}

func part2(filename string) {
	lines := util.GetLines(filename)
	grid := waterGrid{
		DenseGrid: util.NewDenseGridFromLines(lines),
	}
	start := grid.FindCoordinates('S')[0]
	lastRow := grid.Height - 1

	memo := make(map[util.Coordinate]int, grid.Height*grid.Width)
	var dfs func(util.Coordinate) int
	dfs = func(pos util.Coordinate) int {
		// out of bounds or empty cell. fail
		if grid.At(pos) == nil {
			return 0
		}

		// reached last row yay
		if pos.Y == lastRow {
			return 1
		}

		if v, ok := memo[pos]; ok {
			return v // use memoized result
		}

		total := 0
		for _, n := range grid.GetNeighbors(pos) {
			total += dfs(n)
		}
		memo[pos] = total
		return total
	}

	result := dfs(start)
	fmt.Println(result)
}

func simulateWater(grid *util.DenseGrid, curr util.Coordinate, visited map[util.Coordinate]struct{}) int {
	split := fallTilSplit(grid, curr)
	if split == nil {
		return 0
	}

	if _, ok := visited[*split]; ok {
		return 0
	}

	visited[*split] = struct{}{}

	left := split.Add(util.Coordinate{X: -1, Y: 0})
	right := split.Add(util.Coordinate{X: 1, Y: 0})

	res := 1
	if val := grid.At(left); val != nil {
		res += simulateWater(grid, left, visited)
	}

	if val := grid.At(right); val != nil {
		res += simulateWater(grid, right, visited)
	}

	return res
}

func fallTilSplit(grid *util.DenseGrid, start util.Coordinate) *util.Coordinate {
	down := start.Add(util.Coordinate{X: 0, Y: 1})
	val := grid.At(down)
	for val != nil && *val == '.' {
		down = down.Add(util.Coordinate{X: 0, Y: 1})
		val = grid.At(down)
	}

	// reached bottom
	if grid.At(down) == nil {
		return nil
	}

	return &down // reached splitter
}
