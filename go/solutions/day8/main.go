package main

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/smort/aoc2025/util"
)

type Point3D struct {
	X, Y, Z float64
}

type Pair struct {
	I, J     int // indexes into point slice
	Distance float64
}

type Graph struct {
	Adj [][]int
}

func main() {
	part1("example.txt", 10)
	part1("input.txt", 1000)

	part2("example.txt", 10)
	part2("input.txt", 1000)
}

func part1(filename string, numConns int) {
	lines := util.GetLines(filename)

	// make array of points
	points := make([]Point3D, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ",")
		x, y, z := util.MustConvAtoi(parts[0]), util.MustConvAtoi(parts[1]), util.MustConvAtoi(parts[2])
		points = append(points, Point3D{float64(x), float64(y), float64(z)})
	}

	// calculate the distance between each pair of points (oof)
	minHeap := allPairsMinHeap(points)

	// make graph
	g := buildGraph(len(points), minHeap, numConns)

	// traverse to get sizes - bfs
	n := len(g.Adj)
	visited := make([]bool, n)
	sizes := make([]int, 0, n)
	for start := range n {
		if visited[start] {
			continue
		}

		queue := []int{start}
		visited[start] = true
		size := 0

		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]
			size++

			for _, node := range g.Adj[curr] {
				if !visited[node] {
					visited[node] = true
					queue = append(queue, node)
				}
			}
		}

		sizes = append(sizes, size)
	}

	slices.SortFunc(sizes, func(a, b int) int {
		return b - a
	})

	var out [3]int
	for i := 0; i < 3 && i < len(sizes); i++ {
		out[i] = sizes[i]
	}

	fmt.Println(out[0] * out[1] * out[2])
}

func part2(filename string, numConns int) {
	lines := util.GetLines(filename)

	// make array of points
	points := make([]Point3D, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, ",")
		x, y, z := util.MustConvAtoi(parts[0]), util.MustConvAtoi(parts[1]), util.MustConvAtoi(parts[2])
		points = append(points, Point3D{float64(x), float64(y), float64(z)})
	}

	// calculate the distance between each pair of points (oof)
	minHeap := allPairsMinHeap(points)

	g := Graph{Adj: make([][]int, len(points))}
	for minHeap.Len() > 0 {
		p := minHeap.Pop()
		g.Adj[p.I] = append(g.Adj[p.I], p.J)
		g.Adj[p.J] = append(g.Adj[p.J], p.I)

		// check if everything connected together
		if isConnected(g) {
			fmt.Printf("%f\n", points[p.I].X*points[p.J].X)
			return
		}
	}

	fmt.Println("FAIL")
}

func distance(a, b Point3D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func allPairsMinHeap(points []Point3D) *util.MinHeap[Pair] {
	heap := util.NewMinHeap[Pair]()

	n := len(points)
	for i := range n {
		for j := i + 1; j < n; j++ {
			d := distance(points[i], points[j])
			heap.Push(Pair{I: i, J: j, Distance: d}, int(d))
		}
	}

	return heap
}

func buildGraph(n int, heap *util.MinHeap[Pair], k int) Graph {
	g := Graph{Adj: make([][]int, n)}

	// keep connecting until we reach k nodes
	for i := 0; i < k && heap.Len() > 0; i++ {
		p := heap.Pop()
		g.Adj[p.I] = append(g.Adj[p.I], p.J)
		g.Adj[p.J] = append(g.Adj[p.J], p.I)
	}

	return g
}

func isConnected(g Graph) bool {
	n := len(g.Adj)
	visited := make([]bool, n)
	queue := []int{0} // start at node 0
	visited[0] = true
	count := 1

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, node := range g.Adj[curr] {
			if !visited[node] {
				visited[node] = true
				queue = append(queue, node)
				count++
			}
		}
	}

	return count == n
}
