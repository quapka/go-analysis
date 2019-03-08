/*
 * Genarate rsa keys.
 */

package main

import (
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

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Write header
	if _, err := f.Write([]byte("id;n;e;d;p;q;t1\n")); err != nil {
		log.Fatal(err)
	}

	// TODO rand.Redaer is global and maybe we don't need to pass it along
	// but just call it, when it is needed
	reader := rand.Reader

	for id := 0; id < *keyCount; id++ {
		var dataRow string = ""
		var data string = ""

		idPrefix := fmt.Sprintf("%d;", id)

		if rsaCmd.Happened() {
			data, err = getRSAData(reader, *bitSize)
			if err != nil {
				// fmt.Printf("Got errorneous data: %s\n", data)
				log.Fatal(err)
			}
		} else if eccCmd.Happened() {
			log.Fatal("Generating ECC keys is not implemented yet.")
		} else {
			log.Fatal("Please, specify a 'rsa' or 'ecc' key type.")
		}

		dataRow = idPrefix + data + "\n"

		// Write keys
		if _, err := f.Write([]byte(dataRow)); err != nil {
			log.Fatal(err)
		}

	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func getRSAData(reader io.Reader, bitSize int) (data string, err error) {
	start := time.Now()
	key, err := rsa.GenerateKey(reader, bitSize)
	end := time.Now()
	t1 := end.Sub(start)

	if err != nil {
		fmt.Println("We got a problem")
		return "", err
	}

	n := key.PublicKey.N
	e := key.PublicKey.E
	d := key.D
	p := key.Primes[0]
	q := key.Primes[1]

	return fmt.Sprintf("%x;%x;%x;%x;%x;%d;", n, e, d, p, q, t1), nil
}
