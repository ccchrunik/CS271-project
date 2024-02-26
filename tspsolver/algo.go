package tspsolver

import (
	"math"
	"math/rand"
)

func getInitialPaths(dm [][]float64, num int, seed int) []*Path {
	paths := []*Path{}

	r := rand.New(rand.NewSource(int64(seed)))
	for i := 0; i < num; i++ {
		if i < len(dm) {
			paths = append(paths, nearestNeighborPath(dm, i))
		} else {
			paths = append(paths, randomPath(dm, r))
		}
	}

	return paths
}

func nearestNeighborPath(dm [][]float64, startNode int) *Path {
	notVisited := map[int]bool{}
	nodes := []int{}
	for i := 0; i < len(dm); i++ {
		notVisited[i] = true
	}

	nodes = append(nodes, startNode)
	delete(notVisited, startNode)
	node := startNode
	for i := 0; i < len(dm)-1; i++ {
		nextNode := 0
		minDist := math.MaxFloat64
		for k := range notVisited {
			if dm[node][k] < minDist {
				minDist = dm[node][k]
				nextNode = k
			}
		}
		nodes = append(nodes, nextNode)
		delete(notVisited, nextNode)
		node = nextNode
	}

	return &Path{
		nodes: nodes,
		dm:    dm,
	}
}

func randomPath(dm [][]float64, r *rand.Rand) *Path {
	n := len(dm)
	nodes := []int{}
	for i := 0; i < n; i++ {
		nodes = append(nodes, i)
	}
	r.Shuffle(n, func(i, j int) { nodes[i], nodes[j] = nodes[j], nodes[i] })
	// fmt.Println(nodes)
	return &Path{
		nodes: nodes,
		dm:    dm,
	}
}

// func _2opt(path *Path) *Path {

// }

func _3opt(path *Path) *Path {
	improved := true
	n := path.Len()

OUTER:
	for improved {
		for i := 0; i < n-2; i++ {
			for j := i + 1; j < n-1; j++ {
				for k := j + 1; k < n; k++ {
					diff := 0.0
					conf := 0
					for l := 1; l < 8; l++ {
						localDiff := path.calculateDiff3(l, i, j, k)
						if localDiff < diff {
							diff = localDiff
							conf = l
						}
					}

					if diff < 0 {
						path.swap3(conf, i, j, k)
						continue OUTER
					}
				}
			}
		}
		improved = false
	}

	return nil
}
