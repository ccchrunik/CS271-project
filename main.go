package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"tsp/tspsolver"
)

func exp() {
	for i := 100; i < 1001; i += 100 {
		matrixName := fmt.Sprintf("data/tsp-problem-%d-1000-500-100-1.txt", i)
		dm, _ := tspsolver.ReadMatrix(matrixName)

		fmt.Printf("%-8s: %d.\n", "Size", i)
		// r := rand.New(rand.NewSource(0))
		solver := tspsolver.New(dm, 0)
		// random := tspsolver.RandomPath(dm, r)
		// fmt.Printf("%-8s: %f.\n", "Random", random.Distance())

		start1 := time.Now()
		nn := solver.SolveSLS(1, "nn")
		fmt.Printf("%-8s: %f. (%vms)\n", "NN", nn.Distance(), time.Since(start1).Milliseconds())

		start2 := time.Now()
		after2 := solver.SolveSLS(1, "2opt")
		fmt.Printf("%-8s: %f. (%vms)\n", "After2", after2.Distance(), time.Since(start2).Milliseconds())

		start3 := time.Now()
		after3 := solver.SolveSLS(1, "3opt")
		fmt.Printf("%-8s: %f. (%vms)\n", "After3", after3.Distance(), time.Since(start3).Milliseconds())

		fmt.Println()
	}
}

type Output struct {
	fileName  string
	iteration int
	result    *tspsolver.Path
	duration  time.Duration
}

type Input struct {
	fileName string
	order    []int
	outCh    chan *Output
}

func NewInput(dirName string, fileName string) *Input {
	order := []int{}

	s := strings.TrimSuffix(fileName, ".txt")
	fields := strings.Split(s, "-")[2:]
	for _, v := range fields {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		order = append(order, i)
	}

	return &Input{
		fileName: path.Join(dirName, fileName),
		order:    order,
	}
}

func main() {

	dir := "./data/Competion"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}

	inputs := []*Input{}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "tsp-problem") {
			continue
		}
		inputs = append(inputs, NewInput(dir, file.Name()))
	}
	sort.Slice(inputs, func(i, j int) bool {
		// return inputFiles[i].order[0] < inputFiles[j].order[0]
		for k := 0; k < 5; k++ {
			if inputs[i].order[k] < inputs[j].order[k] {
				return true
			} else if inputs[i].order[k] > inputs[j].order[k] {
				return false
			}
		}
		return false
	})

	worker := func(id int, jobs <-chan *Input) {
		for input := range jobs {
			dm, err := tspsolver.ReadMatrix(input.fileName)
			if err != nil {
				log.Fatal(err)
			}
			solver := tspsolver.New(dm, 0)
			start3 := time.Now()
			factor := 1000.0 / float64(solver.Len())
			iteration := int(math.Min(1000.0, math.Ceil(3.0*factor*factor*factor)))
			// iteration := 1

			after3 := solver.SolveSLS(iteration, "3opt")
			input.outCh <- &Output{
				fileName:  input.fileName,
				iteration: iteration,
				duration:  time.Since(start3),
				result:    after3,
			}
		}
	}

	totalStart := time.Now()

	resultChs := []chan *Output{}
	workerNum := runtime.NumCPU()
	jobsCh := make(chan *Input, len(inputs))

	for i := 0; i < workerNum; i++ {
		go worker(i, jobsCh)
	}

	for _, input := range inputs {
		outCh := make(chan *Output, 1)
		resultChs = append(resultChs, outCh)
		input.outCh = outCh
		jobsCh <- input
	}
	close(jobsCh)

	for _, ch := range resultChs {
		output := <-ch
		fmt.Println(output.fileName)
		fmt.Fprintln(f, output.fileName)
		fmt.Printf("%-8s: %f. (n:%d) (%vms)\n\n", "3-opt", output.result.Distance(), output.iteration, output.duration.Milliseconds())
		fmt.Fprintf(f, "%-8s: %f. (n:%d) (%vms)\n", "3-opt", output.result.Distance(), output.iteration, output.duration.Milliseconds())
		output.result.Write(f)
	}
	f.Close()

	fmt.Printf("Total Time: %vs\n", time.Since(totalStart).Seconds())
}
