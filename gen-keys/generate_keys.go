
/*
 * Genarate rsa keys.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"log"
	"time"
)

func main() {

	// If the file doesn't exist, create it, or append to the file
    f, err := os.OpenFile("rsa2048.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }

    // Write header
    if _, err := f.Write([]byte("id;n;e;d;p;q;t1\n")); err != nil {
        log.Fatal(err)
    }

	reader := rand.Reader
	bitSize := 2048

	for id := 0; id < 100; id++ {

		start := time.Now()
		key, rsaerr := rsa.GenerateKey(reader, bitSize)
		end := time.Now()
		// t1 = end - start
		t1 := end.Sub(start)

		checkError(rsaerr)

		publicKey := key.PublicKey

		var n = publicKey.N
		var e = publicKey.E
		var d = key.D
		var p = key.Primes[0]
		var q = key.Primes[1]

	    // Write keys
	    if _, err := f.Write([]byte(fmt.Sprintf("%d;%x;%x;%x;%x;%x;%d;\n", id,n,e,d,p,q,t1))); err != nil {
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
