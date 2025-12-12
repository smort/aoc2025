package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/smort/aoc2025/util"
)

func main() {
	part1("example.txt")
	part1("input.txt")
}

var dimensionRegex = regexp.MustCompile(`\d+x.+`)

type Shape struct {
	Width     int
	Height    int
	Mask      [][]bool // . or not
	Rotations []*Shape
}

type Region struct {
	Width  int
	Height int
	Counts []int
}

func part1(filename string) {
	result := 0

	shapes, regions := parseAll(filename)

	fmt.Printf("parsed %d shapes and %d regions\n", len(shapes), len(regions))
	for i, reg := range regions {
		feasible := regionFeasibleByArea(reg, shapes)
		if feasible {
			result++
		}
		fmt.Printf("region %d: %dx%d counts=%v feasible_by_area=%v\n", i, reg.Width, reg.Height, reg.Counts, feasible)
	}

	fmt.Println(result)
}

func pack(r Region, shapes map[int]*Shape) {
}

// 0:
// ###
// ##.
// ##.
//
// ... followed by lines like:
// 12x5: 1 0 1 0 3 2
func parseAll(filename string) (map[int]*Shape, []Region) {
	lines := util.GetLines(filename)

	shapes := map[int]*Shape{}
	regions := []Region{}

	// parse shapes until we hit a dimension line
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		if dimensionRegex.MatchString(line) {
			break
		}

		idStr := strings.TrimSuffix(line, ":")
		id := util.MustConvAtoi(idStr)

		// read rows until blank line or next header
		rows := [][]bool{}
		j := i + 1
		for j < len(lines) {
			l := strings.TrimSpace(lines[j])
			if l == "" {
				break
			}

			parts := strings.Split(l, "")
			row := make([]bool, len(parts))
			for k, char := range parts {
				row[k] = (char == "#")
			}
			rows = append(rows, row)
			j++
		}

		// make shapes
		if len(rows) > 0 {
			h := len(rows)
			w := len(rows[0])
			for r := range rows {
				if len(rows[r]) < w {
					nr := make([]bool, w)
					copy(nr, rows[r])
					rows[r] = nr
				}
			}

			// precalculate rotations to save time
			s := &Shape{Width: w, Height: h, Mask: rows}
			r90 := rotateShape(s)
			r180 := rotateShape(r90)
			r270 := rotateShape(r180)
			s.Rotations = []*Shape{r90, r180, r270}
			shapes[id] = s
		}

		i = j
	}

	// rest are regions
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++
		if line == "" {
			continue
		}

		// 4x4: 0 0 0 0 2 0
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		dimension := strings.TrimSpace(parts[0])
		countsStr := strings.Split(strings.TrimSpace(parts[1]), " ")
		dimParts := strings.Split(dimension, "x")

		w := util.MustConvAtoi(dimParts[0])
		h := util.MustConvAtoi(dimParts[1])

		counts := make([]int, 0, len(countsStr))
		total := 0
		for _, cs := range countsStr {
			// how much of each present
			v := util.MustConvAtoi(cs)
			counts = append(counts, v)
			total += v
		}

		regions = append(regions, Region{Width: w, Height: h, Counts: counts})
	}

	return shapes, regions
}

func rotateShape(s *Shape) *Shape {
	h := s.Height
	w := s.Width
	newW := h
	newH := w
	newMask := make([][]bool, newH)
	for r := range newH {
		newMask[r] = make([]bool, newW)
	}

	for r := range h {
		for c := range w {
			newMask[c][h-1-r] = s.Mask[r][c]
		}
	}

	return &Shape{Width: newW, Height: newH, Mask: newMask}
}

// just check whether its even possible a little bit
func regionFeasibleByArea(reg Region, shapes map[int]*Shape) bool {
	total := 0
	for idx, count := range reg.Counts {
		if count <= 0 {
			continue
		}
		s := shapes[idx]
		total += shapeArea(s) * count
	}
	return total <= reg.Width*reg.Height
}

// shapeArea returns the number of actually occupied cells
func shapeArea(s *Shape) int {
	count := 0
	for r := range s.Height {
		for c := range s.Width {
			if s.Mask[r][c] {
				count++
			}
		}
	}
	return count
}
