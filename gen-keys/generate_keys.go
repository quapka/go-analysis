/*
 * Generate rsa keys.
 */

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/akamensky/argparse"
	"io"
	"log"
	"os"
	"time"
)

const typeRSA string = "rsa"
const typeECC string = "ecc"

func main() {
	// I'm having troubles with declaration/definition and scope in Go,
	// therefore I'm declaring err here, so later I can use
	// if <var>, err = [...] <-- no colon here
	// that is, definition without declaration
	var err error
	// os.Args[0] will be the name of the executable and therefore give
	// more meaningful help message
	parser := argparse.NewParser(os.Args[0], "Generates RSA and ECC key for further analysis")

	rsaCmd := parser.NewCommand(typeRSA, "Generate a RSA private key")
	// bit-size is valid option only for RSA
	bitSize := rsaCmd.Int("b", "bit-size", &argparse.Options{Required: true, Help: "Specify the bit size of the private key"})

	eccCmd := parser.NewCommand(typeECC, "Generate a ECC private key")

	keyCount := parser.Int("c", "key-count", &argparse.Options{Required: true, Help: "Specify the number of keys to be generated"})

	err = parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// If the file doesn't exist, create it, or append to the file
	filename := ""
	if rsaCmd.Happened() {
		filename = fmt.Sprintf("%s_%db_%d.csv", typeRSA, *bitSize, *keyCount)
	} else if eccCmd.Happened() {
		filename = fmt.Sprintf("%s_%d.csv", typeECC, *keyCount)
	} else {
		log.Fatal("Please, specify a 'rsa' or 'ecc' key type.")
	}

	outFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// TODO rand.Reader is global and maybe we don't need to pass it along
	// but just call it, when it is needed
	reader := rand.Reader

	if rsaCmd.Happened() {
		// add RSA header to the output csv
		writeRow(outFile, "id;n;e;d;p;q;t1;")
		// generate <keyCount> RSA keys and save them to the output csv
		for id := 0; id < *keyCount; id++ {
			data := getRSAData(reader, *bitSize)
			writeRow(outFile, fmt.Sprintf("%d;%s", id, data))
		}

	} else if eccCmd.Happened() {
		// in contrast to the project assignment
		// both coordinates (x, y) of the curve are saved
		writeRow(outFile, "id;e;d;t1;")
		// generate <keyCount> ECC keys and save them to the output csv
		for id := 0; id < *keyCount; id++ {
			data := getECCData(reader)
			writeRow(outFile, fmt.Sprintf("%d;%s", id, data))
		}
	} else {
		log.Fatal("Please, specify a 'rsa' or 'ecc' key type.")
	}

	if err := outFile.Close(); err != nil {
		log.Fatal(err)
	}
}

func getRSAData(reader io.Reader, bitSize int) (data string) {
	start := time.Now()
	key, err := rsa.GenerateKey(reader, bitSize)
	end := time.Now()
	elapsed := end.Sub(start)

	if err != nil {
		log.Fatal(err)
	}

	n := key.PublicKey.N
	e := key.PublicKey.E
	d := key.D
	p := key.Primes[0]
	q := key.Primes[1]
	t1 := elapsed.Nanoseconds()

	return fmt.Sprintf("%x;%x;%x;%x;%x;%d;", n, e, d, p, q, t1)
}

func getECCData(reader io.Reader) (data string) {
	curve := elliptic.P256()
	start := time.Now()
	key, err := ecdsa.GenerateKey(curve, reader)
	end := time.Now()
	elapsed := end.Sub(start)

	if err != nil {
		log.Fatal(err)
	}

	//we will store both coordinates of the point corresponding to the public key
	x := key.PublicKey.X
	y := key.PublicKey.Y
	d := key.D
	t1 := elapsed.Nanoseconds()

	return fmt.Sprintf("04%x%x;%x;%d;", x, y, d, t1)
}

func writeRow(file *os.File, row string) {
	if _, err := file.Write([]byte(row + "\n")); err != nil {
		log.Fatal(err)
	}
}
