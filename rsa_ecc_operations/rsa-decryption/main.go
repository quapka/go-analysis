package main

/*
 * Partially taken and changed from:
 * https://golang.org/src/crypto/rsa/example_test.go
 */

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	key := make([]byte, 32)
	rng := rand.Reader
	rsaCiphertext, _ := hex.DecodeString("aabbccddeeff")
	rsaPrivateKey, _ := rsa.GenerateKey(rng, 2048)

	rsa.DecryptPKCS1v15SessionKey(rng, rsaPrivateKey, rsaCiphertext, key)
}
