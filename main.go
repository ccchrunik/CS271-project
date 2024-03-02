package main

import (
	"fmt"
	"math/rand"
	"time"
	"tsp/tspsolver"
)

func main() {
	for i := 100; i < 1001; i += 100 {
		matrixName := fmt.Sprintf("data/tsp-problem-%d-1000-500-100-1.txt", i)
		dm, _ := tspsolver.ReadMatrix(matrixName)

		fmt.Printf("%-8s: %d.\n", "Size", i)
		r := rand.New(rand.NewSource(0))
		solver := tspsolver.New(dm, 0)
		before := tspsolver.RandomPath(dm, r)
		fmt.Printf("%-8s: %f.\n", "Before", before.Distance())
		start2 := time.Now()
		after2 := solver.SolveSLS(100, "2opt")
		fmt.Printf("%-8s: %f. (%vms)\n", "After2", after2.Distance(), time.Since(start2).Milliseconds())
		start3 := time.Now()
		after3 := solver.SolveSLS(1, "3opt")
		fmt.Printf("%-8s: %f. (%vms)\n", "After3", after3.Distance(), time.Since(start3).Milliseconds())
		fmt.Println()
	}

}
