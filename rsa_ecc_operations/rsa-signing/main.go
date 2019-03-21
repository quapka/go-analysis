package main

/*
 * Partially taken and changed from:
 * https://golang.org/src/crypto/rsa/example_test.go
 */

import (
	"crypto/rand"
	"crypto/rsa"
)

func main() {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	signature, _ := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, []byte("Hello World!"))
	_ = signature

}
