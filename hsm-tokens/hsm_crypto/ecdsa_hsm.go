package hsm_crypto

import (
	"crypto/elliptic"
	"errors"
	// "fmt"
	"encoding/hex"
	"github.com/miekg/pkcs11"
	"io"
)

// values copied from the RFC #6637 section 11.
// https://www.ietf.org/rfc/rfc6637.txt
const P_256_DER = "06082A8648CE3D030107"
const P_384_DER = "06052B81040022"
const P_521_DER = "06052B81040023"

// FIXME what is priv if it has not been initialized?
// maybe return pointer and return nil in case of a error
func GenerateECDSAKey(c elliptic.Curve, rand io.Reader, hsmInstance *Hsm) (priv PrivateKey, err error) {

	if !hsmInstance.isInitialized() {
		return priv, errors.New("hsm has not been initialized")
	}

	labelSize := 64
	// tokenLabel := []byte(hsmInstance.hsmInfo.tokenLabel)
	publicKeyLabel := make([]byte, labelSize)
	_, err = rand.Read(publicKeyLabel)
	if err != nil {
		return priv, err
	}
	privateKeyLabel := make([]byte, labelSize)
	_, err = rand.Read(privateKeyLabel)
	if err != nil {
		return priv, err
	}

	ecdsaParams, err := getCurveParamsInDER(c.Params().Name)
	if err != nil {
		return priv, err
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
	}

	publicObjHandle, privateObjHandle, err := hsmInstance.Ctx.GenerateKeyPair(
		hsmInstance.SessionHandle,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA_KEY_PAIR_GEN, nil)},
		publicKeyTemplate,
		privateKeyTemplate)
	if err != nil {
		return priv, err
	}

	priv.PublicKey = PublicKey{
		hsmInstance,
		publicKeyLabel,
		publicObjHandle,
	}
	priv = PrivateKey{
		priv.PublicKey,
		privateKeyLabel,
		privateObjHandle,
	}
	return priv, nil

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
