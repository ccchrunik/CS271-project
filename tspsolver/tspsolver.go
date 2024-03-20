package tspsolver

import (
	"fmt"
	"math"
	"os"
)

func New(dm [][]float64, seed int) *TspSolver {
	return &TspSolver{
		dm:   dm,
		seed: seed,
	}
}

func NewPath(dm [][]float64, nodes []int) *Path {
	return &Path{
		dm:    dm,
		nodes: nodes,
	}
}

type Path struct {
	nodes []int
	dm    [][]float64
}

func (p *Path) Write(f *os.File) {
	fmt.Fprintf(f, "%d\n%v\n\n", len(p.nodes), p.nodes)
}

func (p *Path) Print(prefix string) {
	fmt.Printf("%s Nodes: %v. Distance: %f\n", prefix, p.nodes, p.Distance())
}

func (p *Path) Len() int {
	return len(p.nodes)
}

func (p *Path) Distance() float64 {
	dist := 0.0
	n := len(p.nodes)
	for i := 0; i < n; i++ {
		dist += p.dm[p.nodes[i]][p.nodes[(i+1)%n]]
	}
	return dist
}

func (p *Path) cost(x, y int) float64 {
	return p.dm[p.nodes[x]][p.nodes[y]]
}

func (p *Path) calculateDiff2(i, j int) float64 {
	n := len(p.nodes)
	return -p.cost(i, i+1) - p.cost(j, (j+1)%n) + p.cost(i, j) + p.cost(i+1, (j+1)%n)
}

func (p *Path) calculateDiff3(l, i, j, k int) float64 {

	n := len(p.nodes)

	switch l {
	case 1: // 001
		return -p.cost(i, i+1) - p.cost(j, j+1) + p.cost(i, j) + p.cost(i+1, j+1)
	case 2: // 010
		return -p.cost(j, j+1) - p.cost(k, (k+1)%n) + p.cost(j, k) + p.cost(j+1, (k+1)%n)
	case 4: // 100
		return -p.cost(k, (k+1)%n) - p.cost(i, i+1) + p.cost(k, i) + p.cost((k+1)%n, i+1)
	case 3: // 011
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, (k+1)%n) + p.cost(i, j) + p.cost(i+1, k) + p.cost(j+1, (k+1)%n)
	case 6: // 110
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, (k+1)%n) + p.cost(j, k) + p.cost(j+1, i) + p.cost((k+1)%n, i+1)
	case 5: // 101
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, (k+1)%n) + p.cost(k, i) + p.cost((k+1)%n, j) + p.cost(i+1, j+1)
	case 7: // 111
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, (k+1)%n) + p.cost(i+1, k) + p.cost(j+1, i) + p.cost((k+1)%n, j)
	default: // 000
		// original path
		return 0.0
	}
}

// swap portion of the list given two index
func reverse(list []int) {
	n := len(list)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		list[i], list[j] = list[j], list[i]
	}
}

func (p *Path) swap2(i, j int) *Path {
	newPath := &Path{}
	newPath.dm = p.dm
	newPath.nodes = append([]int{}, p.nodes...)
	reverse(newPath.nodes[i+1 : j+1])
	return newPath
}

func (p *Path) swap3(l, i, j, k int) *Path {
	newPath := &Path{}
	newPath.dm = p.dm
	newPath.nodes = []int{}

	// no matter what different, the first idx of the new path will be node i
	// ex: [0, 1, 2, 3, 4, 5, 6, 7]
	// (i = 0, j = 4) -> [0, 4, 3, 2, 1, 5, 6, 7]
	// append [i+1, j] segment
	newPath.nodes = append(newPath.nodes, p.nodes[i+1:j+1]...)

	// append [j+1, k] segment
	newPath.nodes = append(newPath.nodes, p.nodes[j+1:k+1]...)

	// append [k+1, i] segment
	newPath.nodes = append(newPath.nodes, p.nodes[k+1:]...)
	newPath.nodes = append(newPath.nodes, p.nodes[:i+1]...)

	// ex: (i, j, k) = (2, 4, 7) -> (3, 4), (5, 6, 7), (0, 1, 2)
	if l&1 != 0 {
		reverse(newPath.nodes[0 : j-i])
	}
	if l&2 != 0 {
		reverse(newPath.nodes[j-i : k-i])
	}
	if l&4 != 0 {
		reverse(newPath.nodes[k-i:])
	}

	return newPath
}

type TspSolver struct {
	dm   [][]float64
	seed int
}

func (ts *TspSolver) Len() int {
	return len(ts.dm)
}

func (ts *TspSolver) SolveSLS(n int, method string) *Path {
	paths := getInitialPaths(ts.dm, n, ts.seed)
	optPaths := []*Path{}

	for i, p := range paths {
		if method == "nn" {
			if i < len(ts.dm) {
				optPaths = append(optPaths, p)
			}
		} else if method == "2opt" {
			optPaths = append(optPaths, Opt2(p))
		} else {
			optPaths = append(optPaths, Opt3(p))
		}
	}

	var minPath *Path
	minDist := math.MaxFloat64
	for _, p := range optPaths {
		dist := p.Distance()
		if dist < minDist {
			minDist = dist
			minPath = p
		}
	}

	return minPath
}
