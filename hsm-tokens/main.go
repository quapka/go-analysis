package main

import (
	"./hsm_crypto"
	"fmt"
	"github.com/miekg/pkcs11"
	"log"
)

func main() {
	p := pkcs11.New("/usr/lib/softhsm/libsofthsm2.so")
	err := p.Initialize()
	if err != nil {
		panic(err)
	}

	defer p.Destroy()
	defer p.Finalize()

	slots, err := p.GetSlotList(true)
	if err != nil {
		panic(err)
	}

	session, err := p.OpenSession(slots[0], pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		panic(err)
	}
	defer p.CloseSession(session)

	fmt.Println(slots)
	fmt.Println(session)

	err = p.Login(session, pkcs11.CKU_USER, "1234")
	if err != nil {
		panic(err)
	}
	defer p.Logout(session)

	obj, _, _ := p.FindObjects(session, 10)
	fmt.Println(obj)

	tokenPersistent := true
	tokenLabel := "label"
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_RSA),
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, tokenPersistent),
		pkcs11.NewAttribute(pkcs11.CKA_VERIFY, true),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, []byte{1, 0, 0, 0, 1}),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
	}
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, tokenPersistent),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
		pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
		pkcs11.NewAttribute(pkcs11.CKA_EXTRACTABLE, true),
	}
	pbk, pvk, e := p.GenerateKeyPair(
		session,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)},
		publicKeyTemplate,
		privateKeyTemplate)
	if e != nil {
		log.Fatal("failed to generate keypair: %s\n", e)
	}
	fmt.Printf("Session: %d.\n", session)
	fmt.Println(fmt.Sprintf("%x", pbk))
	fmt.Println(fmt.Sprintf("%x", pvk))
	fmt.Println(fmt.Sprintf("%x", publicKeyTemplate[4].Value))
	// p.GenerateKeyPair(
	//     session,
	//     []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)},
	//     []*pkcs11.Attribute)
	fmt.Println("Print")
	// key, err := p.
	// p.DigestInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA_1, nil)})
	// hash, err := p.Digest(session, []byte("this is a string"))
	// if err != nil {
	// 	panic(err)
	// }

	// for _, d := range hash {
	// 	fmt.Printf("%x", d)
	// }
	// fmt.Println()

	publicKeyTemplate = []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, tokenLabel),
	}

	p.FindObjectsInit(session, publicKeyTemplate)
	obj, _, _ = p.FindObjects(session, 100)
	fmt.Println(obj)

	p.Destroy()
	p.Finalize()

	p.CloseSession(session)

	p.Logout(session)

	fmt.Println("----------------------------------")

	hsmInstance := hsm_crypto.New("/usr/lib/softhsm/libsofthsm2.so", "pv204", "1234")
	hsmInstance.Initialize()
	_ = hsmInstance

	// var pub hsm.PublicKey
	// pub.Hsm = hsmInstance
	// pub.KeyLabel = tokenLabel

	// err = pub.Initialize()
	// fmt.Println(err)
	// defer pub.Finalize()

	// h, err := pub.FindKeyHandle()
	// fmt.Println(err)
	// fmt.Println(h)

	bitSize := uint(1024)
	key, _ := hsm_crypto.GenerateKeyHSM(bitSize, hsmInstance)
	// _ = key
	privKey := key.PrivateKeyHSM
	_ = privKey

	message := []byte("Hello World")
	signature, err := privKey.Sign(nil, message) //, nil)
	if err != nil {
		log.Fatal("Failed to create a signature", err)
	}
	fmt.Println(fmt.Sprintf("%s", signature))

	// _ = key
}
