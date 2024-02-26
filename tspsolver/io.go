package tspsolver

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

func ReadMatrix(fileName string) ([][]float64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	dm := [][]float64{}
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		if lineNum == 0 {
			n, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}
			lineNum = n
			continue
		}

		dv := []float64{}
		for _, v := range strings.Fields(scanner.Text()) {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			dv = append(dv, f)
		}
		dm = append(dm, dv)
	}

	if lineNum != len(dm) {
		return nil, errors.New("line number doesn't match distance matrix")
	}

	return dm, nil
}

func ReadPath(fileName string) ([]int, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	path := []int{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}
		path = append(path, i)
	}

	return path, nil
}
