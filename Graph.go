package Cykle

import (
	"fmt"
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

func init() {
	//rand.Seed(42)
	rand.Seed(time.Now().UnixNano())
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

	numberOfEdges := int(math.Floor((density * float64(size*(size-1))) / 2))
	for numberOfEdges > 0 {
		vertexOne := rand.Intn(size) + 1
		vertexTwo := rand.Intn(size) + 1

		for !g.CheckEdge(vertexOne, vertexTwo) && vertexOne != vertexTwo {
			g.AddEdge(vertexOne, vertexTwo)
			numberOfEdges--
		}
	}

	for i := 0; i < size-1; i++ {
		if g.VerticesList[i].GetVertexDegree()%2 == 1 {
			vertexTwo := rand.Intn(size-i) + i + 1
			if g.CheckEdge(g.VerticesList[i].Root, vertexTwo) {
				g.RemoveEdge(g.VerticesList[i].Root, vertexTwo)

			} else {
				g.AddEdge(g.VerticesList[i].Root, vertexTwo)
			}
		}
	}

	if !g.CheckIfEulerian() || g.CheckIfForIsolated() || !g.CheckIfGraphIsConnected() {
		return NewGraphWithDensity(size, density)
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

func (g Graph) CheckIfForIsolated() bool {
	for _, v := range g.VerticesList {
		if v.Adjacent.Size() == 0 {
			return true
		}
	}
	return false
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

func (g *Graph) MeasureHamilton(all bool) (int, float64, float64) {
	t := TimingStuff{
		timeStart: time.Now(),
		timeFirst: 0,
		nResults:  0,
		all:       all,
	}

	//hamilton(g, len(g.VerticesList), g.VerticesList[0].Root, []int{}, []int{}, all, &nResults, timeStart, &timeFirst)
	first := g.VerticesList[0].Root
	NewHamilton(g, []int{first}, []int{first}, first, len(g.VerticesList), &t)
	timeForAll := time.Since(t.timeStart).Seconds()

	if t.timeFirst == float64(0) {
		t.timeFirst = timeForAll
	}

	return t.nResults, t.timeFirst, timeForAll
}

type TimingStuff struct {
	timeStart time.Time
	timeFirst float64
	nResults  int
	all       bool
}

func NewHamilton(g *Graph, res, visited []int, current, size int, t *TimingStuff) {
	v := g.GetVertex(current)

	if len(res) == size {
		if v.Adjacent.Contains(g.VerticesList[0].Root) {
			if t.timeFirst == float64(0) {
				t.timeFirst = time.Since(t.timeStart).Seconds()
			}
			t.nResults++
			return
		}
	}

	for _, adj := range v.Adjacent.Values() {
		if t.timeFirst != float64(0) && !t.all {
			return
		}
		if !slices.Contains(visited, adj) {
			NewHamilton(g, append(res, adj), append(visited, adj), adj, size, t)
		}
	}
}
