package types

import (
	"github.com/miekg/pkcs11"
)

type PrivateKeyHSM struct {
	privateKey pkcs11.ObjectHandle
	publicKey  pkcs11.ObjectHandle
}

// TODO is PublicKeyHSM necessary?
type PublicKeyHSM struct {
}
