package Cykle

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/stacks/arraystack"
	"golang.org/x/exp/slices"
	"math"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

type Graph struct {
	VerticesList []Vertex
}

type Vertex struct {
	Root     int
	Adjacent *LinkedList
}

func (v Vertex) GetAllNextUnvisited(visited []int) []int {
	var res []int

	for _, edge := range v.Adjacent.Values() {
		if !slices.Contains(visited, edge) {
			res = append(res, edge)
		}
	}

	return res
}

func NewVertex(value int) Vertex {
	v := Vertex{Root: value, Adjacent: NewLinkedList()}
	return v
}

func NewGraph(size int) Graph {
	g := Graph{VerticesList: make([]Vertex, size)}
	for i := 0; i < size; i++ {
		g.VerticesList[i] = NewVertex(i + 1)
	}

	return g
}

func (g Graph) GetVertex(needle int) *Vertex {
	i, _ := slices.BinarySearchFunc(g.VerticesList, Vertex{Root: needle, Adjacent: nil}, comp)

	return &g.VerticesList[i]
}

func comp(v1, v2 Vertex) int {
	if v1.Root == v2.Root {
		return 0
	} else if v1.Root > v2.Root {
		return 1
	} else {
		return -1
	}
}

func (g Graph) EdgeCount() (res int) {
	for _, v := range g.VerticesList {
		res += v.Adjacent.Size()
	}
	return res / 2
}

func (g Graph) AddEdge(vertexOne, vertexTwo int) {
	g.GetVertex(vertexOne).Adjacent.AddSingle(vertexTwo, vertexOne, g.GetVertex(vertexTwo).Adjacent)
}

func (g Graph) CheckEdge(vertexOne, vertexTwo int) bool {
	return g.GetVertex(vertexOne).Adjacent.Contains(vertexTwo)
}

func NewGraphWithDensity(size int, density float64) Graph {
	g := NewGraph(size)

	initial := hashset.New()
	for _, v := range g.VerticesList {
		initial.Add(v.Root)
	}
	// make sure the graph is connected
	// [0, size) => [1, size+1)
	curr := rand.Intn(size) + 1
	initial.Remove(curr)

	for initial.Size() > 0 {
		values := initial.Values()
		adj := values[rand.Intn(len(values))].(int)
		g.AddEdge(curr, adj)
		initial.Remove(adj)
		curr = adj
	}

	numberOfEdges := int(math.Floor(density*float64(size*(size-1))) / 2)
	missingEdges := numberOfEdges - g.EdgeCount()

	for i := 0; i < missingEdges; i++ {
		var vertexOne, vertexTwo int
		for ok := true; ok; ok = vertexOne == vertexTwo || g.CheckEdge(vertexOne, vertexTwo) {
			vertexOne = rand.Intn(size) + 1
			vertexTwo = rand.Intn(size) + 1
		}
		g.AddEdge(vertexOne, vertexTwo)
	}

	uneven := g.getAllUneven()
	for i := 0; i < len(uneven); i += 2 {
		g.AddEdge(uneven[i], uneven[i+1])
	}

	return g
}

func (g Graph) getAllUneven() (uneven []int) {
	for _, v := range g.VerticesList {
		if v.GetVertexDegree()%2 != 0 {
			uneven = append(uneven, v.Root)
		}
	}
	return uneven
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func (g Graph) CheckIfGraphIsConnected() bool {
	values := []int{g.VerticesList[0].Root}
	stack := arraystack.New()
	stack.Push(values[0])

	numberOfVertices := len(g.VerticesList)

	for len(values) != numberOfVertices {
		if stack.Size() == 0 {
			return false
		}

		stackLast, _ := stack.Peek()
		node := g.GetVertex(stackLast.(int))
		allNext := node.GetAllNextUnvisited(values)

		if len(allNext) == 0 {
			stack.Pop()
			continue
		}

		values = append(values, allNext[0])
		stack.Push(allNext[0])
	}

	return true
}

func (g Graph) GetDensity() float64 {
	numberOfEdges := 0
	for _, v := range g.VerticesList {
		numberOfEdges += v.Adjacent.Size()
	}
	numberOfVertices := len(g.VerticesList)
	return float64(numberOfEdges) / float64(numberOfVertices*(numberOfVertices-1))
}

func (g Graph) String() string {
	var res strings.Builder
	res.WriteRune('\n')
	for _, v := range g.VerticesList {
		res.WriteString(fmt.Sprintf("%d: %s\n", v.Root, v.Adjacent.String()))
	}

	return res.String()
}

func (g Graph) RemoveEdge(vertexOne, vertexTwo int) {
	g.GetVertex(vertexOne).Adjacent.Remove(vertexTwo, g.GetVertex(vertexTwo).Adjacent)
}

func (v Vertex) GetVertexDegree() (degree int) {
	degree += v.Adjacent.Size()

	return degree
}

func (g Graph) CheckIfEulerian() bool {
	for _, v := range g.VerticesList {
		if v.GetVertexDegree()%2 != 0 {
			return false
		}
	}
	return true
}

func (g Graph) GetEulerianCircuit() []int {
	var res []int
	stack := arraystack.New()
	stack.Push(g.VerticesList[0].Root)

	for stack.Size() > 0 {
		stackLast, _ := stack.Peek()
		node := g.GetVertex(stackLast.(int))

		if node.Adjacent.Size() > 0 {
			stack.Push(node.Adjacent.first.value)
			g.RemoveEdge(stackLast.(int), node.Adjacent.first.value)
		} else {
			val, _ := stack.Pop()
			res = append(res, val.(int))
		}
	}

	return res
}

var nResults int
var timeStart time.Time
var timeFirst time.Duration
var timeFirstTaken bool

func (g Graph) MeasureHamilton(all bool) (int, float64, float64) {
	nResults = 0
	timeFirstTaken = false

	timeStart = time.Now()
	hamilton(g, len(g.VerticesList), g.VerticesList[0].Root, []int{}, []int{}, all)
	timeForAll := time.Since(timeStart)

	return nResults, timeFirst.Seconds(), timeForAll.Seconds()
}

func hamilton(g Graph, size, current int, res, visited []int, all bool) {
	if !all && timeFirstTaken {
		return
	}

	res = append(res, current)
	if len(res) != size {
		visited = append(visited, current)
		for _, adj := range g.GetVertex(current).Adjacent.Values() {
			if !slices.Contains(visited, adj) {
				hamilton(g, size, adj, res, visited, all)
			}
		}
		idx := slices.Index(visited, current)
		visited = slices.Delete(visited, idx, idx)
	} else if g.CheckEdge(current, g.VerticesList[0].Root) {
		nResults += 1
		if timeFirstTaken == false {
			timeFirst = time.Since(timeStart)
			timeFirstTaken = true
		}
	}

	res = res[:len(res)-1]
}