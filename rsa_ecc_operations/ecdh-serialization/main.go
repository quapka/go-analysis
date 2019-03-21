package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func main() {
	p256Curve := elliptic.P256()
	privateKey, _ := ecdsa.GenerateKey(p256Curve, rand.Reader)

	_ = elliptic.Marshal(p256Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
}
