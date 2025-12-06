package main

import (
	"fmt"
	"strings"

	"github.com/smort/aoc2025/util"
)

type problem struct {
	nums     []int
	operator string
}

func main() {
	part1("input.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0
	lines := util.GetLines(filename)

	problems := make([]problem, 0)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		i := 0
		nums := strings.Split(line, " ")
		for _, num := range nums {
			if num == "" {
				continue
			}

			if i >= len(problems) {
				problems = append(problems, problem{nums: make([]int, 0)})
			}
			p := &problems[i]

			if num == "+" || num == "*" {
				p.operator = num
			} else {
				p.nums = append(p.nums, util.MustConvAtoi(num))
			}
			i++
		}
	}

	for _, p := range problems {
		total := 0
		for _, n := range p.nums {
			switch p.operator {
			case "+":
				result += n
			case "*":
				if total == 0 {
					total = 1
				}
				total *= n
			}
		}
		result += total
	}

	fmt.Println(result)
}

func part2(filename string) {
	result := 0
	lines := util.GetLines(filename)

	numLines := len(lines)
	totalLen := len(lines[0])

	problems := make([]problem, 0)

	currCol := totalLen - 1
	p := problem{nums: make([]int, 0)}
	for currCol >= 0 {
		currRow := 0
		numStr := ""
		for currRow < numLines {
			currVal := rune(lines[currRow][currCol])
			if currVal == '+' || currVal == '*' {
				p.operator = string(currVal)
				if numStr != "" {
					num := util.MustConvAtoi(numStr)
					p.nums = append(p.nums, num)
					numStr = ""
				}
				problems = append(problems, p)
				p = problem{nums: make([]int, 0)}
				break
			}
			if currVal != ' ' {
				numStr += string(currVal)
			}

			currRow++
		}
		if numStr != "" {
			num := util.MustConvAtoi(numStr)
			p.nums = append(p.nums, num)
		}

		currCol--
	}

	for _, p := range problems {
		total := 0
		for _, n := range p.nums {
			switch p.operator {
			case "+":
				result += n
			case "*":
				if total == 0 {
					total = 1
				}
				total *= n
			}
		}
		result += total
	}

	fmt.Println(result)
}
