package main

import (
	"crypto/elliptic"
	"fmt"
	"github.com/akamensky/argparse"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// setup parser
	parser := argparse.NewParser(os.Args[0], "Generate time measurements for ScalarBaseMult function on p256.")
	numIterationArg := parser.Int("n", "num-iteration", &argparse.Options{Required: true, Help: "Specify the number of iteration for each input"})
	_ = parser.String("f", "filename", &argparse.Options{Required: true, Help: "[NOT IMPLEMENTED YET] The path to the file that contains the inputs."})

	// parse argument
	err := parser.Parse(os.Args)
	// if parsing failed print usage and return
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// open outfile
	outFile, err := os.OpenFile("data.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// write column headers
	for j := 0; j < 256; j++ {
		if _, err := outFile.WriteString(fmt.Sprintf("%d;", j)); err != nil {
			log.Fatal(err)
		}
	}
	// write endline
	if _, err := outFile.WriteString("\n"); err != nil {
		log.Fatal(err)
	}

	p256 := elliptic.P256()

	for i := 0; i < *numIterationArg; i++ {
		for j := 0; j < 256; j++ {
			// measured function input parameters
			scalar := []byte{byte(math.Pow(2, float64(j)))}

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
