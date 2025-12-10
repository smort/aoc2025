package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/smort/aoc2025/util"
)

type Machine struct {
	indicators []int
	ops        [][]int
	joltage    []int
}

func main() {
	part1("example.txt")
	part1("input.txt")

	part2("example.txt")
	part2("input.txt")
}

func part1(filename string) {
	result := 0
	lines := util.GetLines(filename)
	machines := makeMachines(lines)

	for _, m := range machines {
		start := initIndicators(len(m.indicators))
		pressed := make(map[string]int)
		pressed[serialize(start)] = 0
		queue := make([][]int, 0)
		queue = append(queue, start)

		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]

			if slices.Equal(curr, m.indicators) {
				result += pressed[serialize(curr)]
				break
			}

			// generate next states
			steps := pressed[serialize(curr)]
			for _, op := range m.ops {
				newIndicators := make([]int, len(curr))
				copy(newIndicators, curr)
				for _, idx := range op {
					newIndicators[idx] *= -1
				}

				ser := serialize(newIndicators)

				if _, ok := pressed[ser]; !ok {
					pressed[ser] = steps + 1
					queue = append(queue, newIndicators)
				} else if pressed[ser] > steps+1 {
					pressed[ser] = steps + 1
					queue = append(queue, newIndicators)
				}
			}
		}
	}

	fmt.Println(result)
}

func part2(filename string) {
	result := 0
	lines := util.GetLines(filename)
	machines := makeMachines(lines)
	fmt.Println("total machines:", len(machines))

	for i, m := range machines {
		fmt.Printf("processing machine %d\n", i+1)
		steps := solveMachineZ3(m)
		if steps >= 0 {
			result += steps
		}
	}

	fmt.Println("total presses:", result)
}

func initIndicators(num int) []int {
	indicators := make([]int, num)
	for i := 0; i < num; i++ {
		indicators[i] = -1
	}
	return indicators
}

func serialize(indicators []int) string {
	var sb strings.Builder
	for _, ind := range indicators {
		if ind < 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

// solveMachineZ3 uses an ilp constraint solver to find minimal presses
func solveMachineZ3(m Machine) int {
	numOps := len(m.ops)
	numJoltage := len(m.joltage)

	// build coefficient matrix A where A[j][i] = number of times operation i affects joltage j
	A := make([][]int, numJoltage)
	for j := 0; j < numJoltage; j++ {
		A[j] = make([]int, numOps)
		for i, op := range m.ops {
			for _, pos := range op {
				if pos == j {
					A[j][i]++
				}
			}
		}
	}

	return solveCBC(A, m.joltage, numOps)
}

func solveCBC(A [][]int, targets []int, numVars int) int {
	// generate LP format - thanks AI
	lpContent := generateLP(A, targets, numVars)

	// write to temp file
	tmpFile := "/tmp/problem.lp"
	err := os.WriteFile(tmpFile, []byte(lpContent), 0644)
	if err != nil {
		fmt.Printf("error writing LP file: %v\n", err)
		return -1
	}

	// call cbc
	cmd := exec.Command("cbc", tmpFile, "solve", "solution", "/tmp/solution.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("cbc error. output: %s\n %w", string(output), err))
	}

	// parse - thanks AI
	return parseCBCSolution("/tmp/solution.txt", numVars)
}

// generateLP creates LP format for the constraint system
func generateLP(A [][]int, targets []int, numVars int) string {
	var sb strings.Builder

	// Objective function: minimize sum of all variables
	sb.WriteString("Minimize\n")
	sb.WriteString("obj: ")
	for i := 0; i < numVars; i++ {
		if i > 0 {
			sb.WriteString(" + ")
		}
		sb.WriteString(fmt.Sprintf("x%d", i))
	}
	sb.WriteString("\n\n")

	// Constraints
	sb.WriteString("Subject To\n")
	for i := 0; i < len(A); i++ {
		sb.WriteString(fmt.Sprintf("c%d: ", i))
		first := true
		for j := 0; j < numVars; j++ {
			if A[i][j] != 0 {
				if !first && A[i][j] > 0 {
					sb.WriteString(" + ")
				} else if A[i][j] < 0 {
					sb.WriteString(" - ")
				}

				if abs(A[i][j]) == 1 {
					sb.WriteString(fmt.Sprintf("x%d", j))
				} else {
					sb.WriteString(fmt.Sprintf("%dx%d", abs(A[i][j]), j))
				}
				first = false
			}
		}
		sb.WriteString(fmt.Sprintf(" = %d\n", targets[i]))
	}
	sb.WriteString("\n")

	// Variable bounds (non-negative integers)
	sb.WriteString("Bounds\n")
	for i := 0; i < numVars; i++ {
		sb.WriteString(fmt.Sprintf("x%d >= 0\n", i))
	}
	sb.WriteString("\n")

	// Integer variables
	sb.WriteString("Integer\n")
	for i := 0; i < numVars; i++ {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(fmt.Sprintf("x%d", i))
	}
	sb.WriteString("\n\n")

	sb.WriteString("End\n")

	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// parseCBCSolution reads the CBC solution file and extracts the objective value
func parseCBCSolution(filename string, numVars int) int {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading solution file: %v\n", err)
		return -1
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for the objective value line: "Optimal - objective value 50.00000000"
		if strings.Contains(line, "objective value") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "value" && i+1 < len(parts) {
					if value, err := strconv.ParseFloat(parts[i+1], 64); err == nil {
						return int(value + 0.5) // round to nearest integer
					}
				}
			}
		}
	}

	// If no objective value found, try summing individual variables (fallback)
	total := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Variable lines: "      0 x0                     2                       1"
		if len(line) > 10 && line[8:10] == "x" {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				if value, err := strconv.ParseFloat(parts[2], 64); err == nil {
					total += int(value + 0.5)
				}
			}
		}
	}

	return total
}

func makeMachines(lines []string) []Machine {
	machines := make([]Machine, 0, len(lines))
	for _, line := range lines {
		// example line: [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
		// square brackets are indicators, parentheses are ops, curly braces are joltage
		m := Machine{}

		// parse indicators
		indicatorsEnd := strings.LastIndex(line, "]")
		for i := range line[1:indicatorsEnd] {
			if line[i+1] == '#' {
				m.indicators = append(m.indicators, 1)
			} else {
				m.indicators = append(m.indicators, -1)
			}
		}

		// parse ops
		joltageStart := strings.Index(line, "{")
		opsParts := strings.Split(line[indicatorsEnd+1:joltageStart], " ")
		for _, opPart := range opsParts {
			if opPart == "" {
				continue
			}
			nums := strings.Split(opPart[1:len(opPart)-1], ",")
			ops := make([]int, 0, len(nums))
			for _, num := range nums {
				ops = append(ops, util.MustConvAtoi(num))
			}
			m.ops = append(m.ops, ops)
		}

		// parse joltage
		joltageNums := strings.Split(line[joltageStart+1:len(line)-1], ",")
		nums := make([]int, 0, len(joltageNums))
		for _, num := range joltageNums {
			nums = append(nums, util.MustConvAtoi(num))
		}
		m.joltage = nums

		machines = append(machines, m)
	}

	return machines
}
