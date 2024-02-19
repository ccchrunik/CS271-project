package tspsolver

import "math"

func New(dm [][]float64) *TspSolver {
	return &TspSolver{
		dm: dm,
	}
}

type Path struct {
	nodes []int
	dm    [][]float64
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

func (p *Path) calculateDiff3(l, i, j, k int) float64 {
	switch l {
	case 1: // 001
		return -p.cost(i, i+1) - p.cost(j, j+1) + p.cost(i, j) + p.cost(i+1, j+1)
	case 2: // 010
		return -p.cost(j, j+1) - p.cost(k, k+1) + p.cost(j, k) + p.cost(j+1, k+1)
	case 4: // 100
		return -p.cost(k, k+1) - p.cost(i, i+1) + p.cost(k, i) + p.cost(k+1, i+1)
	case 3: // 011
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, k+1) + p.cost(i, j) + p.cost(i+1, k) + p.cost(j+1, k+1)
	case 6: // 110
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, k+1) + p.cost(j, k) + p.cost(j+1, i) + p.cost(k+1, i+1)
	case 5: // 101
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, k+1) + p.cost(k, i) + p.cost(k+1, j) + p.cost(i+1, j+1)
	case 7: // 111
		return -p.cost(i, i+1) - p.cost(j, j+1) - p.cost(k, k+1) + p.cost(i+1, k) + p.cost(j+1, i) + p.cost(k+1, j)
	default: // 000
		// original path
		return 0.0
	}
}

func (p *Path) swap3(l, i, j, k int) {
	n := p.Len()
	newNodes := make([]int, n)

	len1 := (j - (i + 1) + 1 + n) % n
	len2 := (k - (j + 1) + 1 + n) % n
	len3 := (i - (k + 1) + 1 + n) % n
	for pos := 0; pos < n; pos++ {
		newNodes = append(newNodes, p.nodes[(pos+i+1)%n])
	}

	// reverse first segment
	if l&1 != 0 {
		for pos := 0; pos < len1/2; pos++ {
			newNodes[pos], newNodes[len1-1-pos] = newNodes[len1-1-pos], newNodes[pos]
		}
	}

	// reverse second segment
	if l&2 != 0 {
		for pos := len1; pos < len1+len2/2; pos++ {
			newNodes[len1+pos], newNodes[len1+len2-1-pos] = newNodes[len1+len2-1-pos], newNodes[len1+pos]
		}
	}

	// reverse third segment
	if l&4 != 0 {
		for pos := len2; pos < len1+len2+len3/2; pos++ {
			newNodes[len1+len2+pos], newNodes[len1+len2+len3-1-pos] = newNodes[len1+len2+len3-1-pos], newNodes[len1+len2+pos]
		}
	}

	p.nodes = newNodes
}

type TspSolver struct {
	dm [][]float64
}

func (ts *TspSolver) solveSLS(n int) *Path {
	paths := getInitialPaths(ts.dm, n)

	for _, p := range paths {
		// or _2opt(p)
		_3opt(p)
	}

	var minPath *Path
	minDist := math.MaxFloat64
	for _, p := range paths {
		dist := p.Distance()
		if dist < minDist {
			minDist = dist
			minPath = p
		}
	}

	return minPath
}
