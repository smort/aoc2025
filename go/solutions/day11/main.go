package main

import (
	"fmt"
	"strings"

	"github.com/smort/aoc2025/util"
)

func main() {
	part1("example.txt")
	part1("input.txt")

	part2("example2.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := int64(0)
	lines := util.GetLines(filename)

	graph := util.AdjList[string]{}
	for _, line := range lines {
		working := strings.ReplaceAll(line, ":", "")
		parts := strings.Split(working, " ")
		graph[parts[0]] = parts[1:]
	}

	result = util.CountAllPaths(graph, "you", "out", make(map[string]int64))

	fmt.Println(result)
}

func part2(filename string) {
	lines := util.GetLines(filename)

	graph := util.AdjList[string]{}
	for _, line := range lines {
		working := strings.ReplaceAll(line, ":", "")
		parts := strings.Split(working, " ")
		graph[parts[0]] = parts[1:]
	}

	// possibilities
	// svr -> fft -> dac -> out
	// svr -> dac -> fft -> out
	poss1 := []string{"svr", "fft", "dac", "out"}
	poss2 := []string{"svr", "dac", "fft", "out"}

	results := []int64{}
	for _, poss := range [][]string{poss1, poss2} {
		steps := [3]int64{}
		for i := 0; i < len(poss)-1; i++ {
			curr := poss[i]
			next := poss[i+1]
			steps[i] = util.CountAllPaths(graph, curr, next, make(map[string]int64))
		}

		results = append(results, steps[0]*steps[1]*steps[2])
	}

	fmt.Println(results[0] + results[1])
}
