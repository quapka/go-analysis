package main

/*
 * This code is based on the example from
 * https://www.socketloop.com/tutorials/golang-example-for-ecdsa-elliptic-curve-digital-signature-algorithm-functions
 */

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	// "crypto/md5"
	"crypto/rand"
	// "hash"
	// "io"
	// "math/big"
	// "os"
)

func main() {

	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

	privatekey := new(ecdsa.PrivateKey)
	privatekey, _ = ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair

	// var pubkey ecdsa.PublicKey
	// pubkey = privatekey.PublicKey

	// fmt.Println("Private Key :")
	// fmt.Printf("%x \n", privatekey)

	// fmt.Println("Public Key :")
	// fmt.Printf("%x \n", pubkey)

	// Sign ecdsa style

	// var h hash.Hash
	// h = md5.New()
	// r := big.NewInt(0)
	// s := big.NewInt(0)

	// io.WriteString(h, "This is a message to be signed and verified by ECDSA!")
	// signhash := h.Sum(nil)

	ecdsa.Sign(rand.Reader, privatekey, []byte("Hello World!"))
	// if serr != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// signature := r.Bytes()
	// signature = append(signature, s.Bytes()...)

	// fmt.Printf("Signature : %x\n", signature)

	// // Verify
	// verifystatus := ecdsa.Verify(&pubkey, signhash, r, s)
	// fmt.Println(verifystatus) // should be true
}
