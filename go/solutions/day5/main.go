package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/smort/aoc2025/util"
)

type FreshRange struct {
	start int
	end   int
}

func main() {
	part1("input.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0
	lines := util.GetLines(filename)

	fresh := true
	freshRanges := []FreshRange{}
	for _, line := range lines {
		if line == "" {
			fresh = false
			continue
		}

		if fresh {
			min, max := util.MustConvAtoi2(strings.Split(line, "-"))
			freshRanges = append(freshRanges, FreshRange{start: min, end: max})
			continue
		}

		num := util.MustConvAtoi(line)
		for _, fr := range freshRanges {
			if num >= fr.start && num <= fr.end {
				result++
				break
			}
		}
	}

	fmt.Println(result)
}

func part2(filename string) {
	result := 0
	lines := util.GetLines(filename)

	freshRanges := []FreshRange{}
	for _, line := range lines {
		if line == "" {
			break
		}

		min, max := util.MustConvAtoi2(strings.Split(line, "-"))
		freshRanges = append(freshRanges, FreshRange{start: min, end: max})
	}

	freshRanges = consolidateRanges(freshRanges)

	for _, fr := range freshRanges {
		result += fr.end - fr.start + 1
	}

	fmt.Println(result)
}

func consolidateRanges(ranges []FreshRange) []FreshRange {
	if len(ranges) == 0 {
		return ranges
	}
	// sort by min so we can do single pass
	sort.Slice(ranges, func(i, j int) bool { return ranges[i].start < ranges[j].start })

	out := []FreshRange{ranges[0]}
	for _, r := range ranges[1:] {
		curr := &out[len(out)-1]

		// figure out if we overlap. handle one range inside the other
		if r.start <= curr.end {
			if r.end > curr.end {
				// do in-place update
				curr.end = r.end
			}
		} else {
			out = append(out, r) // no overlap, so just add it
		}
	}
	return out
}
