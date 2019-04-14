package hsm

import (
	"errors"
	"fmt"
	"github.com/miekg/pkcs11"
)

type PublicKey struct {
	*Hsm            // contains PIN, so it's not really public
	KeyLabel string // FIXME some other identifier of a key?
}

type Bla struct {
}

func (key *PublicKey) FindKeyHandle() (pkcs11.ObjectHandle, error) {

	if key.Ctx == nil {
		return 0, errors.New("hsm has not been initialized")
	}

	err := key.Ctx.FindObjectsInit(
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
