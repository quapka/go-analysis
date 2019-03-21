package main

import (
	"crypto/elliptic"
	"fmt"
	"github.com/akamensky/argparse"
	"log"
	"os"
	"time"
)

func main() {
	parser := argparse.NewParser(os.Args[0], "Generate time measurements for ScalarBaseMult function")
	numIteration := parser.Int("n", "num-iteration", &argparse.Options{Required: true, Help: "Specify the number of iteration for each input"})
	_ := parser.String("f", "filename", &argparse.Options{Required: true, Help: "The path to the file that contains the inputs."})

	p256 := elliptic.P256()
	for id := 0; id < 10; id++ {
		outFile, err := os.OpenFile(fmt.Sprintf("data_%d.csv", id), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		for j := 0; j < 256; j++ {
			if _, err := outFile.Write([]byte(fmt.Sprintf("%d;", j))); err != nil {
				log.Fatal(err)
			}
		}
		if _, err := outFile.Write([]byte("\n")); err != nil {
			log.Fatal(err)
		}
		for i := 0; i < numIteration; i++ {
			for j := 0; j < 256; j++ {
				scalar := []byte{byte(Pow(2, j))}
				start := time.Now()
				_, _ = p256.ScalarBaseMult(scalar)
				end := time.Now()
				elapsed := end.Sub(start)
				t1 := elapsed.Nanoseconds()

				if _, err := outFile.Write([]byte(fmt.Sprintf("%d;", t1))); err != nil {
					log.Fatal(err)
				}

			}
			if _, err := outFile.Write([]byte("\n")); err != nil {
				log.Fatal(err)
			}
		}
		if err := outFile.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
