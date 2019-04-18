package ecdsa

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/miekg/pkcs11"
	"github.com/quapka/go-analysis/hsm-tokens/hsm_crypto"
	"io"
	"math/big"
	"reflect"
	"strings"
)

// values copied from the RFC #6637 section 11.
// https://www.ietf.org/rfc/rfc6637.txt
const P_256_DER = "06082A8648CE3D030107"
const P_384_DER = "06052B81040022"
const P_521_DER = "06052B81040023"

func curveFromDER(DER string) elliptic.Curve {
	// var curveNamesFromDer = make(
	curveNamesFromDER := map[string]elliptic.Curve{
		P_256_DER: elliptic.P256(),
		P_384_DER: elliptic.P384(),
		P_521_DER: elliptic.P521(),
	}

	return curveNamesFromDER[strings.ToUpper(DER)]
}

type PublicKey struct {
	*hsm_crypto.Hsm        // contains PIN, so it's not really public
	KeyLabel        []byte // FIXME some other identifier of a key?
	handle          pkcs11.ObjectHandle
}

type PrivateKey struct {
	PublicKey
	KeyLabel []byte
	handle   pkcs11.ObjectHandle
}

func (privKey *PrivateKey) Public() PublicKey {
	return privKey.PublicKey
}

func (key *PublicKey) FindKeyHandle() (pkcs11.ObjectHandle, error) {

	if !key.Hsm.IsInitialized() {
		return 0, errors.New("hsm has not been initialized")
	}

	err := key.Ctx.FindObjectsInit(
		key.SessionHandle, // key has SessionHandle from Hsm
		[]*pkcs11.Attribute{pkcs11.NewAttribute(pkcs11.CKA_LABEL, key.KeyLabel)})

	if err != nil {
		return 0, err
	}

	objs, _, err := key.Ctx.FindObjects(key.SessionHandle, 1)
	defer key.Ctx.FindObjectsFinal(key.SessionHandle)
	if len(objs) == 0 {
		return 0, errors.New("no keys found")
	}
	if len(objs) > 1 {
		return objs[0], errors.New(fmt.Sprintf("%d keys found", len(objs)))
	}

	return objs[0], err
}

// FIXME what is priv if it has not been initialized?
// maybe return pointer and return nil in case of a error
func GenerateKeyPair(c elliptic.Curve, rand io.Reader, hsmInstance *hsm_crypto.Hsm) (privKey PrivateKey, err error) {

	if !hsmInstance.IsInitialized() {
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

	ecdsaParams, err := getCurveParamsInDER(c.Params().Name)
	if err != nil {
		return privKey, err
	}

	// TODO reason about the attributes we use - which we need and why
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_ECDSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_ECDSA_PARAMS, ecdsaParams),
		// TODO use tokenLabel - to link a key to a token, but slightly redundant
		// pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, publicKeyLabel),
		// pkcs11.NewAttribute(pkcs11.CKF_EC_UNCOMPRESS, true),
	}
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_ECDSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		// pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, false),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, privateKeyLabel),
		// pkcs11.NewAttribute(pkcs11.CKF_EC_UNCOMPRESS, true),
	}

	publicObjHandle, privateObjHandle, err := hsmInstance.Ctx.GenerateKeyPair(
		hsmInstance.SessionHandle,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA_KEY_PAIR_GEN, nil)},
		publicKeyTemplate,
		privateKeyTemplate)
	if err != nil {
		return privKey, err
	}

	privKey.PublicKey = PublicKey{
		hsmInstance,
		publicKeyLabel,
		publicObjHandle,
	}
	privKey = PrivateKey{
		privKey.PublicKey,
		privateKeyLabel,
		privateObjHandle,
	}
	return privKey, nil

}

func getCurveParamsInDER(curveName string) (params []byte, err error) {

	var curveDER string
	switch curveName {
	case "P-256":
		curveDER = P_256_DER
	case "P-384":
		curveDER = P_384_DER
	case "P-521":
		curveDER = P_521_DER
	default:
		return []byte{}, errors.New("unknown curve")
	}
	params, err = hex.DecodeString(curveDER)
	if err != nil {
		return []byte{}, err
	}
	return params, nil
}

func (pubKey *PublicKey) Export() (key ecdsa.PublicKey, err error) {

	if !pubKey.IsInitialized() {
		return key, errors.New("hsm has not been initialized")
	}

	sessionHandle := pubKey.SessionHandle

	keyHandle, err := pubKey.FindKeyHandle()
	if err != nil {
		return key, err
	}

	ecdsaParams, err := pubKey.Ctx.GetAttributeValue(
		sessionHandle,
		keyHandle,
		[]*pkcs11.Attribute{pkcs11.NewAttribute(pkcs11.CKA_ECDSA_PARAMS, nil)},
	)
	if err != nil {
		return key, err
	}

	ecdsaPoint, err := pubKey.Ctx.GetAttributeValue(
		sessionHandle,
		keyHandle,
		[]*pkcs11.Attribute{pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, nil)},
	)

	fmt.Println(hex.EncodeToString(ecdsaPoint[0].Value))
	if err != nil {
		return key, err
	}

	fmt.Println("Point")
	fmt.Println(reflect.TypeOf(ecdsaPoint[0].Value).String())
	points := ecdsaPoint[0].Value[1:]
	// xHEX := ecdsaPoint[0].Value[1 : len(ecdsaPoint[0].Value)/2]
	// yHEX := ecdsaPoint[0].Value[len(ecdsaPoint[0].Value)/2:]
	xHEX := points[:len(points)/2]
	yHEX := points[len(points)/2:]
	// fmt.Println(len(ecdsaPoint[0].Value))
	// fmt.Println(len(xHEX))
	// fmt.Println(len(yHEX))
	// fmt.Println(hex.EncodeToString(ecdsaPoint[0].Value))
	// fmt.Println(hex.EncodeToString(xHEX))
	// fmt.Println(hex.EncodeToString(yHEX))
	// fmt.Println(len(ecdsaPoint[0].Value))
	curveDER := hex.EncodeToString(ecdsaParams[0].Value)
	fmt.Println(curveDER)
	curve := curveFromDER(curveDER)

	X := new(big.Int)
	X.SetString(hex.EncodeToString(xHEX), 16)

	Y := new(big.Int)
	Y.SetString(hex.EncodeToString(yHEX), 16)
	key = ecdsa.PublicKey{curve, X, Y}
	return key, err
}

func (privKey *PrivateKey) sign(digest []byte, m []*pkcs11.Mechanism) (signature []byte, err error) {
	if !privKey.Hsm.IsInitialized() {
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
	return privKey.sign(digest, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)})
}

func (pubKey *PublicKey) verify(digest []byte, signature []byte, m []*pkcs11.Mechanism) (bool, error) {

	if !pubKey.Hsm.IsInitialized() {
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
	return pubKey.verify(digest, signature, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)})
}

// func (privKey *PrivateKey) DeriveSharedSecret(pubKey ecdsa.PublicKey, rand io.Reader) (secret []byte, err error) {
// 	params := pkcs11.NewECDH1DeriveParams(pkcs11.CKD_SHA1_KDF, []byte{}, []byte{})
// 	session := privKey.Hsm.SessionHandle
// 	ctx := privKey.Hsm.Ctx

// 	labelSize := 64
// 	// tokenLabel := []byte(hsmInstance.hsmInfo.tokenLabel)
// 	secretLabel := make([]byte, labelSize)
// 	_, err = rand.Read(secretLabel)
// 	if err != nil {
// 		return secret, err
// 	}

// 	template := []*pkcs11.Attribute{
// 		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
// 		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_GENERIC_SECRET),
// 		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, false),
// 		// pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
// 		// pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
// 		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, false),
// 		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true),
// 		pkcs11.NewAttribute(pkcs11.CKA_LABEL, secretLabel),
// 		pkcs11.NewAttribute(pkcs11., secretLabel),
// 	}

// 	secretHandle, err := ctx.DeriveKey(
// 		session,
// 		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDH1_DERIVE, nil)},
// 		privKey.handle,
// 		template)
// 	if err != nil {
// 		return []byte{}, errors.New("could not derive a new key")
// 	}
// 	return secret, err
// }
