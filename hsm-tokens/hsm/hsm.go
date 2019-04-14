package hsm

import (
	"errors"
	"github.com/miekg/pkcs11"
)

type Hsm struct {
	hsmInfo       hsmInfo
	Ctx           *pkcs11.Ctx
	SessionHandle pkcs11.SessionHandle
}

type hsmInfo struct {
	library    string
	tokenLabel string // FIXME some other identifier of a token?
	pin        string // FIXME make pin *string for better security?
}

func New(library string, tokenLabel string, pin string) *Hsm {
	newHsm := new(Hsm)

	newHsm.hsmInfo.library = library
	newHsm.hsmInfo.tokenLabel = tokenLabel
	newHsm.hsmInfo.pin = pin

	return newHsm
}

func (hsm *Hsm) Initialize() error {

	ctx := pkcs11.New(hsm.hsmInfo.library)
	err := ctx.Initialize()
	if err != nil {
		return err
	}

	slot, err := hsm.findSlot()
	if err != nil {
		return err
	}

	sessionHandle, err := ctx.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION) // FIXME flags??
	if err != nil {
		return err
	}

	err = ctx.Login(sessionHandle, pkcs11.CKU_USER, hsm.hsmInfo.pin) // FIXME usertype??
	if err != nil {
		return err
	}

	hsm.Ctx = ctx
	hsm.SessionHandle = sessionHandle

	return nil
}

func (hsm *Hsm) Finalize() error {

	if hsm.Ctx == nil {
		return errors.New("hsm has already been finalized")
	}

	defer hsm.Ctx.Destroy()
	defer hsm.Ctx.Finalize()
	defer hsm.Ctx.CloseSession(hsm.SessionHandle)
	defer hsm.Ctx.Logout(hsm.SessionHandle)

	hsm.Ctx = nil
	hsm.SessionHandle = 0

	return nil
}

func (hsm *Hsm) findSlot() (uint, error) {
	return 0, nil // TODO actually use tokenLabel to find the slot id
}
