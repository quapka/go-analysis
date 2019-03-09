/*
 * Genarate rsa keys.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/elliptic"
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	var err error

	algorithm := os.Args[1]
	arg2 := os.Args[2]
	arg3 := os.Args[3]

	bitSize, err := strconv.Atoi(arg2)
	if err != nil {
		fmt.Printf("bitSize must be an integer and not %T!\n", arg2)
		os.Exit(1)
	}

	keyCount, err := strconv.Atoi(arg3)
	if err != nil {
		fmt.Printf("keyCount must be an integer and not %T!\n", arg3)
		os.Exit(1)
	}

	// If the file doesn't exist, create it, or append to the file
	filename := fmt.Sprintf("%s%d.csv", algorithm, bitSize)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Write header
	if _, err := f.Write([]byte("id;n;e;d;p;q;t1\n")); err != nil {
		log.Fatal(err)
	}

	// TODO rand.Reader is global and maybe we don't need to pass it along
	// but just call it, when it is needed
	reader := rand.Reader

	for id := 0; id < keyCount; id++ {
		var dataRow string = ""
		var data string = ""

		idPrefix := fmt.Sprintf("%d;", id)

		if algorithm == "rsa" {
			data, err = getRSAData(reader, bitSize)
			if err != nil {
				// fmt.Printf("Got erroneous data: %s\n", data)
				log.Fatal(err)
			}
		} else if algorithm == "ecc" {
			data, err = getECCData(reader)
			// log.Fatal("Generating ECC keys is not implemented yet.")
			if err != nil {
				// fmt.Printf("Got erroneous data: %s\n", data)
				log.Fatal(err)
			}
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

func getECCData(reader io.Reader) (data string, err error) {
	curve := elliptic.P256()
	start := time.Now()
	key, err := ecdsa.GenerateKey(curve, reader)
	end := time.Now()
	elapsed := end.Sub(start)

	if err != nil {
		fmt.Println("We got a problem")
		return "", err
	}

	//we will store both coordinates of the point corresponding to the public key
	x := key.PublicKey.X
	y := key.PublicKey.Y
	d := key.D
	t1 := elapsed.Nanoseconds()

	return fmt.Sprintf("%x;%x;%x;%d;", x, y, d, t1), nil
}