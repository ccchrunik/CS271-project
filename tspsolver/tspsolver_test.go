package tspsolver

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadInput(t *testing.T) {
	matrixName1 := "../data/test_d5.txt"
	pathName1 := "../data/test_s5.txt"

	dm5 := [][]float64{
		{0.0, 3.0, 4.0, 2.0, 7.0},
		{3.0, 0.0, 4.0, 6.0, 3.0},
		{4.0, 4.0, 0.0, 5.0, 8.0},
		{2.0, 6.0, 5.0, 0.0, 6.0},
		{7.0, 3.0, 8.0, 6.0, 0.0},
	}

	dm, err := ReadMatrix(matrixName1)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(dm))
	for i := 0; i < 5; i++ {
		assert.Equal(t, dm5[i], dm[i])
	}

	path, err := ReadPath(pathName1)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(path))
	assert.Equal(t, []int{1, 3, 2, 5, 4}, path)

	for _, v := range []int{15, 17, 26, 48} {
		matrixName := fmt.Sprintf("../data/test_d%d.txt", v)
		pathName := fmt.Sprintf("../data/test_s%d.txt", v)

		dm, err = ReadMatrix(matrixName)
		assert.Nil(t, err)
		assert.Equal(t, v, len(dm))
		for i := 0; i < v; i++ {
			assert.Equal(t, v, len(dm[i]))
		}

		path, err = ReadPath(pathName)
		assert.Nil(t, err)
		assert.Equal(t, v, len(path))
	}
}

func TestRandomPath(t *testing.T) {
	matrixName := "../data/test_d48.txt"
	dm, err := ReadMatrix(matrixName)
	assert.Nil(t, err)

	r := rand.New(rand.NewSource(42))
	for i := 0; i < 10000; i++ {
		path := RandomPath(dm, r)
		nodes := map[int]bool{}
		assert.Equal(t, dm, path.dm)
		assert.Equal(t, 48, len(path.nodes))
		for _, v := range path.nodes {
			if _, ok := nodes[v]; !ok {
				nodes[v] = true
			} else {
				t.Fatal("non-unique node")
			}
		}
	}
}

func TestNearestNeighbor(t *testing.T) {
	matrixName := "../data/test_d15.txt"
	dm, err := ReadMatrix(matrixName)
	assert.Nil(t, err)

	// start from 0
	nn1 := []int{0, 12, 1, 14, 8, 4, 6, 2, 11, 13, 9, 7, 5, 3, 10}
	nnPath1 := NearestNeighborPath(dm, 0)
	assert.Equal(t, len(nn1), len(nnPath1.nodes))
	assert.Equal(t, nn1, nnPath1.nodes)
}

func TestGetInitialPaths(t *testing.T) {
	matrixName := "../data/test_d48.txt"
	dm, err := ReadMatrix(matrixName)
	assert.Nil(t, err)

	paths1 := getInitialPaths(dm, 48, 42)
	assert.Equal(t, 48, len(paths1))
	for _, path := range paths1 {
		assert.Equal(t, len(dm), len(path.nodes))
		// log.Println(path.nodes)
	}

	paths2 := getInitialPaths(dm, 100, 42)
	assert.Equal(t, 100, len(paths2))
	for _, path := range paths2 {
		assert.Equal(t, len(dm), len(path.nodes))
		// log.Println(path.nodes)
	}
}

func TestPathDistanceAndCost(t *testing.T) {
	matrixName := "../data/test_d5.txt"
	dm, err := ReadMatrix(matrixName)
	assert.Nil(t, err)

	path := &Path{}
	// 0.0  3.0  4.0  2.0  7.0
	// 3.0  0.0  4.0  6.0  3.0
	// 4.0  4.0  0.0  5.0  8.0
	// 2.0  6.0  5.0  0.0  6.0
	// 7.0  3.0  8.0  6.0  0.0
	path.dm = dm

	// 3.0 + 4.0 + 5.0 + 6.0 + 7.0 = 25.0
	path.nodes = []int{0, 1, 2, 3, 4}
	assert.Equal(t, 25.0, path.Distance())
	assert.Equal(t, 3.0, path.cost(0, 1))
	assert.Equal(t, 4.0, path.cost(1, 2))
	assert.Equal(t, 5.0, path.cost(2, 3))
	assert.Equal(t, 6.0, path.cost(3, 4))
	assert.Equal(t, 7.0, path.cost(4, 0))
	for i := 0; i < len(dm); i++ {
		for j := i; j < len(dm); j++ {
			if i == j {
				assert.Equal(t, 0.0, path.cost(i, j))
				assert.Equal(t, 0.0, path.cost(j, i))
			} else {
				assert.Equal(t, path.cost(i, j), path.cost(j, i))
			}
		}
	}

	// 2.0 + 6.0 + 3.0 + 8.0 + 4.0 = 23.0
	path.nodes = []int{0, 3, 1, 4, 2}
	assert.Equal(t, 23.0, path.Distance())
	assert.Equal(t, 2.0, path.cost(0, 1))
	assert.Equal(t, 6.0, path.cost(1, 2))
	assert.Equal(t, 3.0, path.cost(2, 3))
	assert.Equal(t, 8.0, path.cost(3, 4))
	assert.Equal(t, 4.0, path.cost(4, 0))
	for i := 0; i < len(dm); i++ {
		for j := i; j < len(dm); j++ {
			assert.Equal(t, path.cost(i, j), path.cost(j, i))
		}
	}

}

func TestDiffAndSwap3(t *testing.T) {
	matrixName := "../data/test_d8.txt"
	dm, err := ReadMatrix(matrixName)
	assert.Nil(t, err)

	before := &Path{}
	before.dm = dm
	after := &Path{}
	after.dm = dm

	// 000: original
	before.nodes = []int{0, 1, 2, 3, 4, 5, 6, 7}
	assert.Equal(t, 0.0, before.calculateDiff3(0, 0, 2, 4))
	assert.Equal(t, 0.0, before.calculateDiff3(0, 0, 3, 6))

	// 001 (0, 1, 5): [0 x 1 x 2, 3, 4, 5 x 6, 7] -> [1, 2, 3, 4, 5, 6, 7, 0]
	after.nodes = []int{1, 2, 3, 4, 5, 6, 7, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(1, 0, 1, 4))
	assert.Equal(t, after.nodes, before.swap3(1, 0, 1, 4).nodes)

	// 001 (1, 2, 6): [0, 1 x 2 x 3, 4, 5, 6 x 7] -> [2, 3, 4, 5, 6, 7, 0, 1]
	after.nodes = []int{2, 3, 4, 5, 6, 7, 0, 1}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(1, 1, 2, 6))
	assert.Equal(t, after.nodes, before.swap3(1, 1, 2, 6).nodes)

	// 001 (0, 4, 7): [0 x 1, 2, 3, 4 x 5, 6, 7 x] -> [4, 3, 2, 1, 5, 6, 7, 0]
	after.nodes = []int{4, 3, 2, 1, 5, 6, 7, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(1, 0, 4, 7))
	assert.Equal(t, after.nodes, before.swap3(1, 0, 4, 7).nodes)

	// 001 (1, 4, 6): [0, 1 x 2, 3, 4 x 5, 6 x 7] -> [4, 3, 2, 5, 6, 7, 0, 1]
	after.nodes = []int{4, 3, 2, 5, 6, 7, 0, 1}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(1, 1, 4, 6))
	assert.Equal(t, after.nodes, before.swap3(1, 1, 4, 6).nodes)

	// 010 (0, 4, 7): [0 x 1, 2, 3, 4 x 5, 6, 7 x] -> [1, 2, 3, 4, 7, 6, 5, 0]
	after.nodes = []int{1, 2, 3, 4, 7, 6, 5, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(2, 0, 4, 7))
	assert.Equal(t, after.nodes, before.swap3(2, 0, 4, 7).nodes)

	// 010 (1, 4, 6): [0, 1 x 2, 3, 4 x 5, 6 x 7] -> [2, 3, 4, 6, 5, 7, 0, 1]
	after.nodes = []int{2, 3, 4, 6, 5, 7, 0, 1}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(2, 1, 4, 6))
	assert.Equal(t, after.nodes, before.swap3(2, 1, 4, 6).nodes)

	// 100 (0, 4, 7): [0 x 1, 2, 3, 4 x 5, 6, 7 x] -> [1, 2, 3, 4, 5, 6, 7, 0]
	after.nodes = []int{1, 2, 3, 4, 5, 6, 7, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(4, 0, 4, 7))
	assert.Equal(t, after.nodes, before.swap3(4, 0, 4, 7).nodes)

	// 100 (1, 4, 6): [0, 1 x 2, 3, 4 x 5, 6 x 7] -> [2, 3, 4, 5, 6, 1, 0, 7]
	after.nodes = []int{2, 3, 4, 5, 6, 1, 0, 7}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(4, 1, 4, 6))
	assert.Equal(t, after.nodes, before.swap3(4, 1, 4, 6).nodes)

	// 011 (0, 4, 7): [0 x 1, 2, 3, 4 x 5, 6, 7 x] -> [4, 3, 2, 1, 7, 6, 5, 0]
	after.nodes = []int{4, 3, 2, 1, 7, 6, 5, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(3, 0, 4, 7))
	assert.Equal(t, after.nodes, before.swap3(3, 0, 4, 7).nodes)

	// 011 (1, 2, 6): [0, 1 x 2 x 3, 4, 5, 6 x 7] -> [2, 6, 5, 4, 3, 7, 0, 1]
	after.nodes = []int{2, 6, 5, 4, 3, 7, 0, 1}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(3, 1, 2, 6))
	assert.Equal(t, after.nodes, before.swap3(3, 1, 2, 6).nodes)

	// 101 (0, 4, 7): [0 x 1, 2, 3, 4 x 5, 6, 7 x] -> [4, 3, 2, 1, 5, 6, 7, 0]
	after.nodes = []int{4, 3, 2, 1, 5, 6, 7, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(5, 0, 4, 7))
	assert.Equal(t, after.nodes, before.swap3(5, 0, 4, 7).nodes)

	// 101 (1, 2, 6): [0, 1 x 2 x 3, 4, 5, 6 x 7] -> [2, 3, 4, 5, 6, 1, 0, 7]
	after.nodes = []int{2, 3, 4, 5, 6, 1, 0, 7}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(5, 1, 2, 6))
	assert.Equal(t, after.nodes, before.swap3(5, 1, 2, 6).nodes)

	// 110 (0, 4, 7): [0 x 1, 2, 3, 4 x 5, 6, 7 x] -> [1, 2, 3, 4, 7, 6, 5, 0]
	after.nodes = []int{1, 2, 3, 4, 7, 6, 5, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(6, 0, 4, 7))
	assert.Equal(t, after.nodes, before.swap3(6, 0, 4, 7).nodes)

	// 110 (1, 2, 6): [0, 1 x 2 x 3, 4, 5, 6 x 7] -> [2, 6, 5, 4, 3, 1, 0, 7]
	after.nodes = []int{2, 6, 5, 4, 3, 1, 0, 7}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(6, 1, 2, 6))
	assert.Equal(t, after.nodes, before.swap3(6, 1, 2, 6).nodes)

	// 111 (0, 4, 7): [0 x 1, 2, 3, 4 x 5, 6, 7 x] -> [4, 3, 2, 1, 7, 6, 5, 0]
	after.nodes = []int{4, 3, 2, 1, 7, 6, 5, 0}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(7, 0, 4, 7))
	assert.Equal(t, after.nodes, before.swap3(7, 0, 4, 7).nodes)

	// 111 (1, 2, 6): [0, 1 x 2 x 3, 4, 5, 6 x 7] -> [2, 6, 5, 4, 3, 1, 0, 7]
	after.nodes = []int{2, 6, 5, 4, 3, 1, 0, 7}
	assert.Equal(t, after.Distance()-before.Distance(), before.calculateDiff3(7, 1, 2, 6))
	assert.Equal(t, after.nodes, before.swap3(7, 1, 2, 6).nodes)
}

func TestSLS(t *testing.T) {
	// matrixName := "../data/tsp-problem-100-1000-50-25-1.txt"
	// matrixName := "../data/tsp-problem-200-1000-500-100-1.txt"
	matrixName := "../data/tsp-problem-1000-1000-500-100-1.txt"

	dm, _ := ReadMatrix(matrixName)

	r := rand.New(rand.NewSource(0))
	solver := New(dm, 0)
	before := RandomPath(dm, r)
	before.Print("Before  ")
	after2 := solver.SolveSLS(1, "2opt")
	after2.Print("After2  ")
	after3 := solver.SolveSLS(1, "3opt")
	after3.Print("After3  ")
}
