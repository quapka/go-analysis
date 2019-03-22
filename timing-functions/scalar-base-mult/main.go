package main

import (
	"crypto/elliptic"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

func main() {
	// setup parser
	parser := argparse.NewParser(os.Args[0], "Generate time measurements for ScalarBaseMult function on p256.")
	numIterationArg := parser.Int("n", "num-iteration", &argparse.Options{Required: true, Help: "Specify the number of iteration for each input"})
	inputFileName := parser.String("f", "filename", &argparse.Options{Required: true, Help: "The path to the file that contains the inputs."})

	// parse argument
	// if parsing failed print usage and return
	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// load input data as [][]string
	data, err := loadInputData(*inputFileName)
	if err != nil {
		log.Fatal(err)
	}

	// open outfile
	outFile, err := os.OpenFile("data.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	p256 := elliptic.P256()

	// write column headers
	for _, row := range data {
		inputLabel := row[0]
		// write input label (column header)
		if _, err := outFile.WriteString(inputLabel + ";"); err != nil {
			log.Fatal(err)
		}
	}

	// write endline
	if _, err := outFile.WriteString("\n"); err != nil {
		log.Fatal(err)
	}

	// iterate over rows
	for i := 0; i < *numIterationArg; i++ {
		// iterate over columns
		for _, row := range data {
			// input for the measured function
			scalar, err := stringToIntBytes(row[1])
			if err != nil {
				log.Fatal(err)
			}

			// timing
			start := time.Now()
			_, _ = p256.ScalarBaseMult(scalar) // measured function
			end := time.Now()
			elapsed := end.Sub(start)
			t1 := elapsed.Nanoseconds()

			// write t1
			if _, err := outFile.WriteString(fmt.Sprintf("%d;", t1)); err != nil {
				log.Fatal(err)
			}

		}

		// write endline
		if _, err := outFile.WriteString("\n"); err != nil {
			log.Fatal(err)
		}
	}

	// close outfile
	if err := outFile.Close(); err != nil {
		log.Fatal(err)
	}
}

// Reads file as csv separated with semicolon.
func loadInputData(fileName string) ([][]string, error) {

	// read whole file into data
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// new csv reader, separator = ;
	r := csv.NewReader(strings.NewReader(string(data)))
	r.Comma = ';'

	// read all
	result, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	return result, err
}

// Converts a number represented as decimal string into bytes.
// Error if the string does not represent a decimal number.
func stringToIntBytes(input string) ([]byte, error) {

	value := new(big.Int)

	value, ok := value.SetString(input, 10)
	if !ok {
		return nil, errors.New("Cannot convert \"" + input + "\" into a number.")
	}

	return value.Bytes(), nil
}
