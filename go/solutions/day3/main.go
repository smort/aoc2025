package main

import (
	"fmt"
	"math"

	"github.com/smort/aoc2025/util"
)

func main() {
	part1("input.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0
	lines := util.GetLines(filename)
	for _, line := range lines {
		result += maxJoltage(line)
	}

	fmt.Println(result)
}

func part2(filename string) {
	result := 0
	lines := util.GetLines(filename)
	for _, line := range lines {
		result += maxJoltageArbitrary(line, 12)
	}

	fmt.Println(result)
}

func maxJoltage(str string) int {
	max := 0
	foundAt := -1
	for idx, r := range str[:len(str)-1] {
		if util.MustConvAtoi(string(r)) > max {
			max = util.MustConvAtoi(string(r))
			foundAt = idx
		}
	}

	secondMax := 0
	for _, r := range str[foundAt+1:] {
		if util.MustConvAtoi(string(r)) > secondMax {
			secondMax = util.MustConvAtoi(string(r))
		}
	}

	return max*10 + secondMax
}

func maxJoltageArbitrary(str string, numDigits int) int {
	numToRemove := len(str) - numDigits // how many digits to remove
	stack := &util.Stack[int]{}

	for _, s := range str {
		curr := util.MustConvAtoi(string(s))
		for !stack.IsEmpty() && numToRemove > 0 {
			top, _ := stack.Peek()
			if curr <= top {
				break // skip since this is a smaller number than we already have
			}
			// we've got a bigger number, pop off the smaller one
			stack.Pop()
			numToRemove--
		}
		stack.Push(curr)
	}

	result := 0
	for i := range numDigits {
		val := stack.Items[i]
		result += val * int(math.Pow(10, float64(numDigits-i-1)))
	}

	println("Result:", result)
	return result
}
