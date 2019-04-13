package main

import (
	"github.com/miekg/pkcs11"
	"types"
)

// Implement the Signer
func (privKey *PrivateKeyHSM) Public() crypto.PublicKey {
	// get the publick key from the token
}

func (privKey *PrivateKeyHSM) Sign(rand io.Reader, digest []byte, opts SignerOpts) (signature []byte, err error) {
	// implement signing a message
}

// TODO: do we need/can use the random on the token?
func GenerateKey(random io.Reader, bits int) (*PrivateKeyHSM, error) {

}
