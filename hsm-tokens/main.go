package main

import (
	"./hsm_crypto"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	// "encoding/asn1"
	"fmt"
	"log"
)

func main() {

	pin := "1234"
	hsmInstance := hsm_crypto.New("/usr/lib/softhsm/libsofthsm2.so", "pv204", &pin)
	err := hsmInstance.Initialize()
	defer hsmInstance.Finalize()
	if err != nil {
		log.Fatal(err)
	}

	// hsmExampleRSA(hsmInstance)

	hsmExampleECDSA(hsmInstance)
}

func hsmExampleRSA(hsmInstance *hsm_crypto.Hsm) {
	bitSize := uint(1024)
	privKey, _ := hsm_crypto.GenerateRsaKey(bitSize, hsmInstance)
	message := []byte("Hello World")

	// test sign
	signature, err := privKey.Sign(nil, message, nil) //, nil)
	fmt.Println(fmt.Sprintf("%s", signature))
	fmt.Println(err)

	// test verify
	ver, err := privKey.Verify(message, signature)
	fmt.Println(ver)
	fmt.Println(err)

	// test encrypt
	cryptotext, err := privKey.Encrypt(message)
	fmt.Println(fmt.Sprintf("%s", cryptotext))
	fmt.Println(err)

	// test decrypt
	plaintext, err := privKey.Decrypt(nil, cryptotext, nil)
	fmt.Println(fmt.Sprintf("%s", plaintext))
	fmt.Println(err)

	// test export
	keyExp, err := privKey.Export()
	fmt.Println(keyExp)
	fmt.Println(err)

	// test encrypt with rsa package
	cryptotext, err = rsa.EncryptPKCS1v15(rand.Reader, &keyExp, message)
	fmt.Println(fmt.Sprintf("%s", cryptotext))
	fmt.Println(err)

	// test decrypt ciphertext encrypted with rsa package
	plaintext, err = privKey.Decrypt(nil, cryptotext, nil)
	fmt.Println(fmt.Sprintf("%s", plaintext))
	fmt.Println(err)
}

func hsmExampleECDSA(hsmInstance *hsm_crypto.Hsm) {
	// message := []byte("Hello World")
	curve := elliptic.P256()
	privKey, err := hsm_crypto.GenerateECDSAKey(curve, rand.Reader, hsmInstance)
	if err != nil {
		fmt.Println(err)
	}
	_ = privKey
	// fmt.Printf("\n%s\n", privKey.KeyLabel)
	// _ = privKey
	// data, _ := asn1.Marshal(curve.Params())
	// fmt.Printf("\n%d\n", len(data))
	// fmt.Println(fmt.Sprintf("%s", data))
}
