/*
 * Genarate rsa keys.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	algorithm := os.Args[1]
	arg2 := os.Args[2]

	bitSize, err := strconv.Atoi(arg2)
	if err != nil {
		fmt.Printf("bitSize must be an integer and not %T!\n", bitSize)
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

	reader := rand.Reader

	for id := 0; id < 10; id++ {

		start := time.Now()
		key, rsaerr := rsa.GenerateKey(reader, bitSize)
		end := time.Now()
		t1 := end.Sub(start)

		checkError(rsaerr)

		publicKey := key.PublicKey

		var n = publicKey.N
		var e = publicKey.E
		var d = key.D
		var p = key.Primes[0]
		var q = key.Primes[1]

		// Write keys
		if _, err := f.Write([]byte(fmt.Sprintf("%d;%x;%x;%x;%x;%x;%d;\n", id, n, e, d, p, q, t1))); err != nil {
			log.Fatal(err)
		}

	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
