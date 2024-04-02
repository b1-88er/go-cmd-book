package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

func sum(data []float64) float64 {
	var sum float64
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

type statsFunc func([]float64) float64

func csv2float(r io.Reader, column int) ([]float64, error) {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true

	// what a sad interface
	column -= 1

	var data []float64
	for i := 0; ; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row %d: %w", i, err)
		}
		// skip headers
		if i == 0 {
			continue
		}

		if len(row) <= column {
			return nil, fmt.Errorf("%w: %s", ErrInvalidColum, row)
		}

		val, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, val)
	}
	return data, nil
}

func run(filenames []string, op string, column int, out io.Writer) error {
	var opFunc statsFunc
	if len(filenames) == 0 {
		return ErrNoFiles
	}

	if column < 1 {
		return ErrInvalidColum
	}

	switch op {
	case "sum":
		opFunc = sum
	case "avg":
		opFunc = avg
	default:
		return ErrInvalidOperation
	}

	results := make([]float64, 0)
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", filename, err)
		}

		data, err := csv2float(f, column)
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("error closing file %s: %v", filename, err)
		}

		results = append(results, data...)
	}

	_, err := fmt.Fprintln(out, opFunc(results))
	return err
}

func main() {
	op := flag.String("op", "sum", "operation to perform")
	column := flag.Int("col", 1, "column to process")
	flag.Parse()

	if err := run(flag.Args(), *op, *column, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
