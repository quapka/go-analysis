package hsm_crypto

import (
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/miekg/pkcs11"
	"io"
)

type PublicKey struct {
	*Hsm            // contains PIN, so it's not really public
	KeyLabel []byte // FIXME some other identifier of a key?
	handle   pkcs11.ObjectHandle
}

type PrivateKey struct {
	*Hsm
	PublicKey
	KeyLabel []byte
	handle   pkcs11.ObjectHandle
}

func (privKey *PrivateKey) Public() PublicKey {
	return privKey.PublicKey
}

func (key *PublicKey) FindKeyHandle() (pkcs11.ObjectHandle, error) {

	if !key.isInitialized() {
		return 0, errors.New("hsm has not been initialized")
	}

	err := key.Ctx.FindObjectsInit(
		key.SessionHandle, // key has SessionHandle from Hsm
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

// FIXME unused bitsize
func GenerateRsaKey(bitSize uint, hsmInstance *Hsm) (privKey PrivateKey, err error) {

	if !hsmInstance.isInitialized() {
		return privKey, errors.New("hsm has not been initialized")
	}

	labelSize := 64
	// tokenLabel := []byte(hsmInstance.hsmInfo.tokenLabel)
	publicKeyLabel := make([]byte, labelSize)
	_, err = rand.Read(publicKeyLabel)
	if err != nil {
		return privKey, err
	}
	privateKeyLabel := make([]byte, labelSize)
	_, err = rand.Read(privateKeyLabel)
	if err != nil {
		return privKey, err
	}

	// TODO reason about the attributes we use - which we need and why
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		// TODO do not fix public exponent
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, []byte{1, 0, 0, 0, 1}),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		// TODO use tokenLabel - to link a key to a token, but slightly redundant
		// pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, publicKeyLabel),
	}
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
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
		return privKey, err
	}

	privKey.PublicKey = PublicKey{hsmInstance, publicKeyLabel, publicObjHandle}
	privKey = PrivateKey{
		hsmInstance,
		privKey.PublicKey,
		privateKeyLabel,
		privateObjHandle,
	}
	return privKey, nil
}

func (privKey *PrivateKey) sign(digest []byte, m []*pkcs11.Mechanism) (signature []byte, err error) {

	if !privKey.isInitialized() {
		return nil, errors.New("hsm has not been initialized")
	}

	ctx := privKey.Hsm.Ctx
	sessionHandle := privKey.Hsm.SessionHandle

	err = ctx.SignInit(sessionHandle, m, privKey.handle)
	if err != nil {
		return nil, err
	}

	return ctx.Sign(sessionHandle, digest)
}

// does not use rand nor opts
func (privKey *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	// FIXME correct mechanism?
	return privKey.sign(digest, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)})
}

func (pubKey *PublicKey) verify(digest []byte, signature []byte, m []*pkcs11.Mechanism) (bool, error) {

	if !pubKey.isInitialized() {
		return false, errors.New("hsm has not been initialized")
	}

	ctx := pubKey.Hsm.Ctx
	sessionHandle := pubKey.Hsm.SessionHandle

	err := ctx.VerifyInit(sessionHandle, m, pubKey.handle)
	if err != nil {
		return false, err
	}

	err = ctx.Verify(sessionHandle, digest, signature)
	return err == nil, err
}

func (pubKey *PublicKey) Verify(digest []byte, signature []byte) (bool, error) {
	// FIXME correct mechanism?
	return pubKey.verify(digest, signature, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)})
}

func (pubKey *PublicKey) encrypt(plaintext []byte, m []*pkcs11.Mechanism) ([]byte, error) {

	if !pubKey.isInitialized() {
		return nil, errors.New("hsm has not been initialized")
	}

	ctx := pubKey.Ctx
	sessionHandle := pubKey.SessionHandle

	err := ctx.EncryptInit(sessionHandle, m, pubKey.handle)
	if err != nil {
		return nil, err
	}

	return ctx.Encrypt(sessionHandle, plaintext)
}

func (pubKey *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	// FIXME correct mechanism?
	return pubKey.encrypt(msg, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)})
}

func (privKey *PrivateKey) decrypt(ciphertext []byte, m []*pkcs11.Mechanism) (plaintext []byte, err error) {

	if !privKey.isInitialized() {
		return nil, errors.New("hsm has not been initialized")
	}

	ctx := privKey.Hsm.Ctx
	sessionHandle := privKey.Hsm.SessionHandle

	err = ctx.DecryptInit(sessionHandle, m, privKey.handle)
	if err != nil {
		return nil, err
	}

	return ctx.Decrypt(sessionHandle, ciphertext)
}

// does not use rand nor opts
func (privKey *PrivateKey) Decrypt(rand io.Reader, msg []byte, opts crypto.DecrypterOpts) (plaintext []byte, err error) {
	// FIXME correct mechanism?
	return privKey.decrypt(msg, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS, nil)})
}
