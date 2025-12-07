package util

import (
	"container/heap"
	"math"
)

type CostFunc func(grid GridInterface, from, to Coordinate) int

type ValidFunc func(grid GridInterface, pos Coordinate) bool

type Coordinate struct {
	X, Y int
}

func (c Coordinate) Add(other Coordinate) Coordinate {
	return Coordinate{X: c.X + other.X, Y: c.Y + other.Y}
}

func (c Coordinate) ManhattanDistance(other Coordinate) int {
	return int(math.Abs(float64(c.X-other.X)) + math.Abs(float64(c.Y-other.Y)))
}

// Standard 4-directional movement
var Directions4 = []Coordinate{
	{0, -1}, // North
	{1, 0},  // East
	{0, 1},  // South
	{-1, 0}, // West
}

// 8-directional movement (including diagonals)
var Directions8 = []Coordinate{
	{0, -1}, {1, -1}, {1, 0}, {1, 1},
	{0, 1}, {-1, 1}, {-1, 0}, {-1, -1},
}

// GridInterface defines the interface for both dense and sparse grids
type GridInterface interface {
	// IsValid checks if a coordinate is valid (within bounds and not blocked)
	IsValid(pos Coordinate) bool
	// GetCost returns the cost to move to this position (1 for unweighted, custom for weighted)
	GetCost(from, to Coordinate) int
	// GetNeighbors returns valid neighboring coordinates
	GetNeighbors(pos Coordinate) []Coordinate
	// FindCoordinates returns all coordinates containing the target rune
	FindCoordinates(target rune) []Coordinate
	// At returns the rune at the given coordinate
	At(pos Coordinate) *rune
	// Set a cell at the given coordinate
	SetCell(pos Coordinate, cell rune)
}

// TODO: use functional options pattern for cost/valid funcs
type DenseGrid struct {
	Width, Height int
	Grid          [][]rune
	Directions    []Coordinate // Directions4 or Directions8
	costFunc      CostFunc     // Optional
	validFunc     ValidFunc    // Optional
}

func NewDenseGrid(width, height int, grid [][]rune) *DenseGrid {
	return &DenseGrid{
		Width:      width,
		Height:     height,
		Grid:       grid,
		Directions: Directions4,
	}
}

func NewDenseGridWithCost(width, height int, grid [][]rune, costFunc CostFunc) *DenseGrid {
	return &DenseGrid{
		Width:      width,
		Height:     height,
		Grid:       grid,
		Directions: Directions4,
		costFunc:   costFunc,
	}
}

func NewDenseGridWithOptions(width, height int, grid [][]rune, costFunc CostFunc, validFunc ValidFunc) *DenseGrid {
	return &DenseGrid{
		Width:      width,
		Height:     height,
		Grid:       grid,
		Directions: Directions4,
		costFunc:   costFunc,
		validFunc:  validFunc,
	}
}

func NewDenseGridFromLines(lines []string) *DenseGrid {
	height := len(lines)
	if height == 0 {
		return nil
	}
	width := len(lines[0])
	grid := make([][]rune, height)
	for y, line := range lines {
		grid[y] = []rune(line)
	}
	return &DenseGrid{
		Width:      width,
		Height:     height,
		Grid:       grid,
		Directions: Directions4,
	}
}

func (g *DenseGrid) IsValid(pos Coordinate) bool {
	// check bounds first
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.Width || pos.Y >= g.Height {
		return false
	}

	if g.validFunc != nil {
		return g.validFunc(g, pos)
	}

	cell := g.Grid[pos.Y][pos.X]
	return cell != '#' // may want to make this configurable?
}

func (g *DenseGrid) GetCost(from, to Coordinate) int {
	if !g.IsValid(to) {
		return -1
	}

	if g.costFunc != nil {
		return g.costFunc(g, from, to)
	}

	// default cost
	return 1
}

func (g *DenseGrid) GetNeighbors(pos Coordinate) []Coordinate {
	neighbors := make([]Coordinate, 0, len(g.Directions))
	for _, dir := range g.Directions {
		next := pos.Add(dir)
		if g.IsValid(next) {
			neighbors = append(neighbors, next)
		}
	}
	return neighbors
}

func (g *DenseGrid) FindCoordinates(target rune) []Coordinate {
	coords := make([]Coordinate, 0)
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Grid[y][x] == target {
				coords = append(coords, Coordinate{X: x, Y: y})
			}
		}
	}
	return coords
}

func (g *DenseGrid) At(pos Coordinate) *rune {
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.Width || pos.Y >= g.Height {
		return nil
	}
	return &g.Grid[pos.Y][pos.X]
}

func (g *DenseGrid) SetCell(pos Coordinate, cell rune) {
	if pos.X >= 0 && pos.Y >= 0 && pos.X < g.Width && pos.Y < g.Height {
		g.Grid[pos.Y][pos.X] = cell
	}
}

// TODO: implement a constructor func that takes a slice of strings
type SparseGrid struct {
	Cells      map[Coordinate]rune
	Directions []Coordinate
	MinX, MaxX int
	MinY, MaxY int
	costFunc   CostFunc  // Optional
	validFunc  ValidFunc // Optional
}

func NewSparseGrid() *SparseGrid {
	return &SparseGrid{
		Cells:      make(map[Coordinate]rune),
		Directions: Directions4,
		MinX:       math.MaxInt32,
		MaxX:       math.MinInt32,
		MinY:       math.MaxInt32,
		MaxY:       math.MinInt32,
	}
}

func NewSparseGridWithCost(costFunc CostFunc) *SparseGrid {
	return &SparseGrid{
		Cells:      make(map[Coordinate]rune),
		Directions: Directions4,
		MinX:       math.MaxInt32,
		MaxX:       math.MinInt32,
		MinY:       math.MaxInt32,
		MaxY:       math.MinInt32,
		costFunc:   costFunc,
	}
}

func NewSparseGridWithOptions(costFunc CostFunc, validFunc ValidFunc) *SparseGrid {
	return &SparseGrid{
		Cells:      make(map[Coordinate]rune),
		Directions: Directions4,
		MinX:       math.MaxInt32,
		MaxX:       math.MinInt32,
		MinY:       math.MaxInt32,
		MaxY:       math.MinInt32,
		costFunc:   costFunc,
		validFunc:  validFunc,
	}
}

func (g *SparseGrid) SetCell(pos Coordinate, cell rune) {
	g.Cells[pos] = cell
	if pos.X < g.MinX {
		g.MinX = pos.X
	}
	if pos.X > g.MaxX {
		g.MaxX = pos.X
	}
	if pos.Y < g.MinY {
		g.MinY = pos.Y
	}
	if pos.Y > g.MaxY {
		g.MaxY = pos.Y
	}
}

func (g *SparseGrid) IsValid(pos Coordinate) bool {
	if g.validFunc != nil {
		return g.validFunc(g, pos)
	}

	// cells not in the map are considered empty/valid by default
	cell, exists := g.Cells[pos]
	if !exists {
		return true
	}
	return cell != '#' // may want to make this configurable?
}

func (g *SparseGrid) GetCost(from, to Coordinate) int {
	if !g.IsValid(to) {
		return -1
	}

	if g.costFunc != nil {
		return g.costFunc(g, from, to)
	}

	// default cost
	return 1
}

func (g *SparseGrid) GetNeighbors(pos Coordinate) []Coordinate {
	neighbors := make([]Coordinate, 0, len(g.Directions))
	for _, dir := range g.Directions {
		next := pos.Add(dir)
		if g.IsValid(next) {
			neighbors = append(neighbors, next)
		}
	}
	return neighbors
}

func (g *SparseGrid) FindCoordinates(target rune) []Coordinate {
	coords := make([]Coordinate, 0)
	for coord, cell := range g.Cells {
		if cell == target {
			coords = append(coords, coord)
		}
	}
	return coords
}

func (g *SparseGrid) At(pos Coordinate) *rune {
	if cell, exists := g.Cells[pos]; exists {
		return &cell
	}
	return nil
}

// PathResult represents the result of a pathfinding operation
type PathResult struct {
	Found    bool
	Path     []Coordinate
	Cost     int
	Visited  int
	Distance map[Coordinate]int
}

// BFS
func BFS(grid GridInterface, start, goal Coordinate) PathResult {
	if start == goal {
		return PathResult{
			Found: true,
			Path:  []Coordinate{start},
			Cost:  0,
		}
	}

	queue := []Coordinate{start}
	visited := make(map[Coordinate]bool)
	parent := make(map[Coordinate]Coordinate)
	distance := make(map[Coordinate]int)

	visited[start] = true
	distance[start] = 0
	visitedCount := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visitedCount++

		// get path if reached goal
		if current == goal {
			path := []Coordinate{}
			for pos := goal; pos != start; pos = parent[pos] {
				path = append([]Coordinate{pos}, path...)
			}
			path = append([]Coordinate{start}, path...)

			return PathResult{
				Found:    true,
				Path:     path,
				Cost:     distance[goal],
				Visited:  visitedCount,
				Distance: distance,
			}
		}

		for _, neighbor := range grid.GetNeighbors(current) {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				distance[neighbor] = distance[current] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return PathResult{
		Found:    false,
		Visited:  visitedCount,
		Distance: distance,
	}
}

type DijkstraNode struct {
	Position Coordinate
	Cost     int
	Index    int // For heap operations
}

type DijkstraHeap []*DijkstraNode

func (h DijkstraHeap) Len() int           { return len(h) }
func (h DijkstraHeap) Less(i, j int) bool { return h[i].Cost < h[j].Cost }
func (h DijkstraHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *DijkstraHeap) Push(x interface{}) {
	n := len(*h)
	node := x.(*DijkstraNode)
	node.Index = n
	*h = append(*h, node)
}

func (h *DijkstraHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	node.Index = -1
	*h = old[0 : n-1]
	return node
}

func Dijkstra(grid GridInterface, start, goal Coordinate) PathResult {
	if start == goal {
		return PathResult{
			Found: true,
			Path:  []Coordinate{start},
			Cost:  0,
		}
	}

	pq := &DijkstraHeap{}
	heap.Init(pq)
	heap.Push(pq, &DijkstraNode{Position: start, Cost: 0})

	distance := make(map[Coordinate]int)
	parent := make(map[Coordinate]Coordinate)
	visited := make(map[Coordinate]bool)
	visitedCount := 0

	distance[start] = 0

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*DijkstraNode)

		if visited[current.Position] {
			continue
		}

		visited[current.Position] = true
		visitedCount++

		// reconstruct path if goal is reached
		if current.Position == goal {
			path := []Coordinate{}
			for pos := goal; pos != start; pos = parent[pos] {
				path = append([]Coordinate{pos}, path...)
			}
			path = append([]Coordinate{start}, path...)

			return PathResult{
				Found:    true,
				Path:     path,
				Cost:     current.Cost,
				Visited:  visitedCount,
				Distance: distance,
			}
		}

		for _, neighbor := range grid.GetNeighbors(current.Position) {
			if visited[neighbor] {
				continue
			}

			cost := grid.GetCost(current.Position, neighbor)
			if cost < 0 {
				continue // Invalid move
			}

			newCost := current.Cost + cost
			if oldCost, exists := distance[neighbor]; !exists || newCost < oldCost {
				distance[neighbor] = newCost
				parent[neighbor] = current.Position
				heap.Push(pq, &DijkstraNode{Position: neighbor, Cost: newCost})
			}
		}
	}

	return PathResult{
		Found:    false,
		Visited:  visitedCount,
		Distance: distance,
	}
}

type AStarNode struct {
	Position Coordinate
	GCost    int // cost from start
	HCost    int // heuristic cost to goal
	FCost    int // total cost (GCost + HCost)
	Index    int // for heap operations
}

type AStarHeap []*AStarNode

func (h AStarHeap) Len() int { return len(h) }
func (h AStarHeap) Less(i, j int) bool {
	if h[i].FCost == h[j].FCost {
		return h[i].HCost < h[j].HCost
	}
	return h[i].FCost < h[j].FCost
}
func (h AStarHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *AStarHeap) Push(x interface{}) {
	n := len(*h)
	node := x.(*AStarNode)
	node.Index = n
	*h = append(*h, node)
}

func (h *AStarHeap) Pop() interface{} {
	old := *h
	n := len(old)
	node := old[n-1]
	node.Index = -1
	*h = old[0 : n-1]
	return node
}

// interface for heuristic functions
type HeuristicFunc func(from, to Coordinate) int

func ManhattanHeuristic(from, to Coordinate) int {
	return from.ManhattanDistance(to)
}

func EuclideanHeuristic(from, to Coordinate) int {
	dx := float64(from.X - to.X)
	dy := float64(from.Y - to.Y)
	return int(math.Sqrt(dx*dx + dy*dy))
}

func AStar(grid GridInterface, start, goal Coordinate, heuristic HeuristicFunc) PathResult {
	if start == goal {
		return PathResult{
			Found: true,
			Path:  []Coordinate{start},
			Cost:  0,
		}
	}

	if heuristic == nil {
		heuristic = ManhattanHeuristic
	}

	pq := &AStarHeap{}
	heap.Init(pq)

	startH := heuristic(start, goal)
	heap.Push(pq, &AStarNode{
		Position: start,
		GCost:    0,
		HCost:    startH,
		FCost:    startH,
	})

	gScore := make(map[Coordinate]int)
	parent := make(map[Coordinate]Coordinate)
	visited := make(map[Coordinate]bool)
	visitedCount := 0

	gScore[start] = 0

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*AStarNode)

		if visited[current.Position] {
			continue
		}

		visited[current.Position] = true
		visitedCount++

		// get the path if goal is reached
		if current.Position == goal {
			path := []Coordinate{}
			for pos := goal; pos != start; pos = parent[pos] {
				path = append([]Coordinate{pos}, path...)
			}
			path = append([]Coordinate{start}, path...)

			return PathResult{
				Found:    true,
				Path:     path,
				Cost:     current.GCost,
				Visited:  visitedCount,
				Distance: gScore,
			}
		}

		for _, neighbor := range grid.GetNeighbors(current.Position) {
			if visited[neighbor] {
				continue
			}

			cost := grid.GetCost(current.Position, neighbor)
			if cost < 0 {
				continue // invalid
			}

			tentativeG := current.GCost + cost
			if oldG, exists := gScore[neighbor]; !exists || tentativeG < oldG {
				gScore[neighbor] = tentativeG
				parent[neighbor] = current.Position

				h := heuristic(neighbor, goal)
				heap.Push(pq, &AStarNode{
					Position: neighbor,
					GCost:    tentativeG,
					HCost:    h,
					FCost:    tentativeG + h,
				})
			}
		}
	}

	return PathResult{
		Found:    false,
		Visited:  visitedCount,
		Distance: gScore,
	}
}

// FloodFill performs a flood fill to find all reachable positions from start
func FloodFill(grid GridInterface, start Coordinate) map[Coordinate]int {
	queue := []Coordinate{start}
	distance := make(map[Coordinate]int)
	visited := make(map[Coordinate]bool)

	distance[start] = 0
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range grid.GetNeighbors(current) {
			if !visited[neighbor] {
				visited[neighbor] = true
				distance[neighbor] = distance[current] + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return distance
}

// BFS but finds all shortest paths from start to goal
func FindAllShortestPaths(grid GridInterface, start, goal Coordinate) [][]Coordinate {
	if start == goal {
		return [][]Coordinate{{start}}
	}

	queue := [][]Coordinate{{start}}
	visited := make(map[Coordinate]int)
	visited[start] = 0
	var allPaths [][]Coordinate
	shortestLength := -1

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		current := path[len(path)-1]

		// if we found a shorter path before, skip longer ones
		if shortestLength != -1 && len(path) > shortestLength {
			continue
		}

		if current == goal {
			if shortestLength == -1 {
				shortestLength = len(path)
			}
			if len(path) == shortestLength {
				pathCopy := make([]Coordinate, len(path))
				copy(pathCopy, path)
				allPaths = append(allPaths, pathCopy)
			}
			continue
		}

		for _, neighbor := range grid.GetNeighbors(current) {
			if dist, exists := visited[neighbor]; !exists || dist >= len(path) {
				visited[neighbor] = len(path)
				newPath := make([]Coordinate, len(path)+1)
				copy(newPath, path)
				newPath[len(path)] = neighbor
				queue = append(queue, newPath)
			}
		}
	}

	return allPaths
}

// Common cost function constructors for convenience

// creates a cost function based on terrain types
func TerrainCostFunc(terrainCosts map[rune]int) CostFunc {
	return func(grid GridInterface, from, to Coordinate) int {
		// Cast to access the underlying grid data
		switch g := grid.(type) {
		case *DenseGrid:
			cell := g.Grid[to.Y][to.X]
			if cost, exists := terrainCosts[cell]; exists {
				return cost
			}
		case *SparseGrid:
			if cell, exists := g.Cells[to]; exists {
				if cost, exists := terrainCosts[cell]; exists {
					return cost
				}
			}
		}
		return 1 // Default cost
	}
}

// creates a cost function that varies by movement direction (north/south)
func DirectionalCostFunc(baseCostFunc CostFunc, uphill, downhill float64) CostFunc {
	return func(grid GridInterface, from, to Coordinate) int {
		baseCost := 1
		if baseCostFunc != nil {
			baseCost = baseCostFunc(grid, from, to)
		}

		// Apply directional modifier
		modifier := 1.0
		if to.Y > from.Y { // Moving south (downhill in many coordinate systems)
			modifier = downhill
		} else if to.Y < from.Y { // Moving north (uphill)
			modifier = uphill
		}

		result := int(float64(baseCost) * modifier)
		// Ensure minimum cost of 1
		if result < 1 {
			result = 1
		}
		return result
	}
}

// creates a cost function that varies by time or other context
func TimeCostFunc(baseCostFunc CostFunc, timeModifier func() float64) CostFunc {
	return func(grid GridInterface, from, to Coordinate) int {
		baseCost := 1
		if baseCostFunc != nil {
			baseCost = baseCostFunc(grid, from, to)
		}

		return int(float64(baseCost) * timeModifier())
	}
}
