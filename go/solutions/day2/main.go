package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/smort/aoc2025/util"
)

func main() {
	part1("input.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0
	ranges := parseInput(filename)

	for _, r := range ranges {
		curr := r.start
		for curr <= r.end {
			if detectRepeat(curr) {
				result += curr
			}
			curr++
		}
	}

	fmt.Println(result)
}

func part2(filename string) {
	result := 0
	ranges := parseInput(filename)

	for _, r := range ranges {
		curr := r.start
		for curr <= r.end {
			if detectRepeat(curr) {
				result += curr
			}
			curr++
		}
	}

	fmt.Println(result)
}

type Range struct {
	start int
	end   int
}

func parseInput(filename string) []Range {
	result := []Range{}
	lines := util.GetLines(filename)
	for _, line := range lines {
		ranges := strings.Split(line, ",")
		for r := range ranges {
			rangeParts := strings.Split(ranges[r], "-")
			start := util.MustConvAtoi(rangeParts[0])
			end := util.MustConvAtoi(rangeParts[1])

			result = append(result, Range{start, end})
		}
	}
	return result
}

func detectRepeat(num int) bool {
	numStr := strconv.Itoa(num)
	numStrLen := len(numStr)
	middle := numStrLen / 2

	startChar := string(numStr[0])

	for j := 1; j < numStrLen; j++ {
		// if we haven't found it by this point, there is no sequence
		if j > middle {
			return false
		}

		if string(numStr[j]) == startChar {
			possibleSequence := numStr[:j]
			if testSequence(numStr[j:], possibleSequence) {
				return true
			}
		}
	}
	return false
}

func testSequence(str string, seq string) bool {
	j := 0 // sequence position
	for i := 0; i < len(str); i++ {
		if j >= len(seq) {
			j = 0
		}
		if string(str[i]) != string(seq[j]) {
			return false
		}
		j++
	}

	return j == len(seq)
}
