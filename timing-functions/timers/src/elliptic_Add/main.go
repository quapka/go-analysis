package main

import (
	"../util"
	"crypto/elliptic"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"log"
	"math/big"
	"os"
	"time"
)

func main() {
	// setup parser
	parser := argparse.NewParser(os.Args[0], "Generate time measurements for Add function on p256.")
	numIterationArg := parser.Int("n", "num-iteration", &argparse.Options{Required: true, Help: "Specify the number of iteration for each input"})
	inputFileName := parser.String("f", "filename", &argparse.Options{Required: true, Help: "The path to the file that contains the inputs."})

	// parse argument
	// if parsing failed print usage and return
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// load input data as [][]string
	data, err := util.LoadInputData(*inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	// open outfile
	outFile, err := os.OpenFile("data.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	p256 := elliptic.P256()

	// new csv writer, separator = ;
	w := csv.NewWriter(outFile)
	w.Comma = ';'

	// get row of column headers
	var headers []string
	for _, row := range data {
		inputLabel := row[0]
		headers = append(headers, inputLabel)
	}
	// write row of column headers
	if err := w.Write(headers); err != nil {
		log.Fatal(err)
	}
	w.Flush()

	// iterate over rows
	for i := 0; i < *numIterationArg; i++ {
		// get row of values
		var values []string
		// iterate over columns
		for _, row := range data {
			// input for the measured function
			x1 := new(big.Int)
			x1, ok := x1.SetString(row[1], 16)
			if !ok {
				log.Fatal(errors.New("Cannot convert \"" + row[1] + "\" into a number."))
			}

			y1 := new(big.Int)
			y1, ok = y1.SetString(row[2], 16)
			if !ok {
				log.Fatal(errors.New("Cannot convert \"" + row[2] + "\" into a number."))
			}

			x2 := new(big.Int)
			x2, ok = x2.SetString(row[3], 16)
			if !ok {
				log.Fatal(errors.New("Cannot convert \"" + row[3] + "\" into a number."))
			}

			y2 := new(big.Int)
			y2, ok = y2.SetString(row[4], 16)
			if !ok {
				log.Fatal(errors.New("Cannot convert \"" + row[4] + "\" into a number."))
			}

			// timing
			start := time.Now()

			_, _ = p256.Params().Add(x1, y1, x2, y2) // measured function

			end := time.Now()
			elapsed := end.Sub(start)
			t1 := elapsed.Nanoseconds()

			values = append(values, fmt.Sprintf("%d", t1))

		}

		// write row of values
		if err := w.Write(values); err != nil {
			log.Fatal(err)
		}
		w.Flush()
	}

	// close outfile
	if err := outFile.Close(); err != nil {
		log.Fatal(err)
	}
}
