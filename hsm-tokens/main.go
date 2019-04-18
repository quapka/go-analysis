package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/quapka/go-analysis/hsm-tokens/hsm_crypto"
	ecdsa_hsm "github.com/quapka/go-analysis/hsm-tokens/hsm_crypto/ecdsa"
	rsa_hsm "github.com/quapka/go-analysis/hsm-tokens/hsm_crypto/rsa"
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

	hsmExampleRSA(hsmInstance)
	hsmExampleECDSA(hsmInstance)
}

func hsmExampleRSA(hsmInstance *hsm_crypto.Hsm) {
	bitSize := uint(1024)
	privKey, _ := rsa_hsm.GenerateKey(bitSize, hsmInstance)
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
	message := []byte("Hello World")
	curve := elliptic.P256()
	privKey, err := ecdsa_hsm.GenerateKeyPair(curve, rand.Reader, hsmInstance)
	if err != nil {
		fmt.Println(err)
	}
	signature, err := privKey.Sign(nil, message, nil) //, nil)
	fmt.Println(fmt.Sprintf("%s", signature))
	fmt.Println(err)

	// test verify
	ver, err := privKey.Verify(message, signature)
	fmt.Println(ver)
	fmt.Println(err)
	// _ = privKey

	keyExp, err := privKey.PublicKey.Export()
	fmt.Println(keyExp)
	fmt.Println(err)

	// fmt.Println(keyExp)
	// curve = elliptic.P384()
	// privKey, err = ecdsa_hsm.GenerateECDSAKey(curve, rand.Reader, hsmInstance)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// _ = privKey

	// curve = elliptic.P521()
	// privKey, err = ecdsa_hsm.GenerateECDSAKey(curve, rand.Reader, hsmInstance)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// _ = privKey
}

func getGoRSAKey() (key *rsa.PrivateKey, err error) {
	reader := rand.Reader
	bitSize := 2048

	key, err = rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func getHSMRSAKey(hsmInstance *hsm_crypto.Hsm) (key rsa_hsm.PrivateKey, err error) {
	bitSize := uint(2048)

	key, err = rsa_hsm.GenerateKey(bitSize, hsmInstance)
	if err != nil {
		return key, err
	}
	return key, nil
}
