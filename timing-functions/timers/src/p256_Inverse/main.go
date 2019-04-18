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

type invertible interface {
	// Inverse returns the inverse of k in GF(P)
	Inverse(k *big.Int) *big.Int
}

func main() {
	// setup parser
	parser := argparse.NewParser(os.Args[0], "Generate time measurements for Inverse function on p256.")
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
			k := new(big.Int)
			k, ok := k.SetString(row[1], 16)
			if !ok {
				log.Fatal(errors.New("Cannot convert \"" + row[1] + "\" into a number."))
			}

			in, ok := p256.(invertible)
			if !ok {
				log.Fatal(errors.New("the curve does not implement invertible"))
			}

			// timing
			start := time.Now()

			_ = in.Inverse(k) // measured function

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
