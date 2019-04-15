package hsm

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/miekg/pkcs11"
	"io"
	"log"
)

type keyHSM struct {
	*PublicKey
	*PrivateKeyHSM
}

// FIXME rename PublicKeyHSM
type PublicKey struct {
	*Hsm            // contains PIN, so it's not really public
	KeyLabel []byte // FIXME some other identifier of a key?
	handle   pkcs11.ObjectHandle
}

type PrivateKeyHSM struct {
	*Hsm
	*PublicKey
	KeyLabel []byte
	handle   pkcs11.ObjectHandle
}

type Bla struct {
}

func (key *PublicKey) FindKeyHandle() (pkcs11.ObjectHandle, error) {

	if key.Ctx == nil {
		return 0, errors.New("hsm has not been initialized")
	}

	err := key.Ctx.FindObjectsInit(
		// FIXME key is missing the SessionHandle
		key.SessionHandle,
		[]*pkcs11.Attribute{pkcs11.NewAttribute(pkcs11.CKA_LABEL, key.KeyLabel)})

	if err != nil {
		return 0, err
	}

	objs, _, err := key.Ctx.FindObjects(key.SessionHandle, 1)
	if len(objs) == 0 {
		return 0, errors.New("no keys found")
	}
	if len(objs) > 1 {
		return objs[0], errors.New(fmt.Sprintf("%d keys found", len(objs)))
	}

	return objs[0], err
}

func GenerateKeyHSM(bitSize uint, hsmInstance *Hsm) (key keyHSM, err error) {
	labelSize := 64
	tokenPersistent := true
	// tokenLabel := []byte(hsmInstance.hsmInfo.tokenLabel)
	publicKeyLabel := make([]byte, labelSize)
	_, err = rand.Read(publicKeyLabel)
	if err != nil {
		return key, err
	}
	privateKeyLabel := make([]byte, labelSize)
	_, err = rand.Read(privateKeyLabel)
	if err != nil {
		return key, err
	}

	// TODO reason about the attributes we use - which we need and why
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, tokenPersistent),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		// TODO do not fix public exponent
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, []byte{1, 0, 0, 0, 1}),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		// TODO use tokenLabel - to link a key to a token, but slightly redundant
		// pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, publicKeyLabel),
	}
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, tokenPersistent),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		// pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, privateKeyLabel),
	}
	publicObjHandle, privateObjHandle, err := hsmInstance.Ctx.GenerateKeyPair(
		hsmInstance.SessionHandle,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)},
		publicKeyTemplate,
		privateKeyTemplate)
	if err != nil {
		return key, err
	}

	// FIXME rename PublicKeyHSM
	key.PublicKey = &PublicKey{hsmInstance, publicKeyLabel, publicObjHandle}
	key.PrivateKeyHSM = &PrivateKeyHSM{
		hsmInstance,
		key.PublicKey,
		privateKeyLabel,
		privateObjHandle,
	}
	return key, nil
}
