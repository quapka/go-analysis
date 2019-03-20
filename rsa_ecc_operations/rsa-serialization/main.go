package main

/* Parts of code based or taken from
 * https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
 */

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func main() {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	privateBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	// using so called blank identifier '_' allows to skip the
	// declared and not used error
	_ = privateBytes
}
