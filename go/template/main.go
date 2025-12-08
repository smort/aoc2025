package main

import (
	"fmt"

	"github.com/smort/aoc2025/util"
)

func main() {
	part1("example.txt")
	part1("input.txt")

	part2("example.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0
	lines := util.GetLines(filename)

	fmt.Println(result)
}

func part2(filename string) {
	result := 0
	lines := util.GetLines(filename)

	fmt.Println(result)
}
