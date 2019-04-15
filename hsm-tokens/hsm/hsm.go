package hsm

import (
	"errors"
	"fmt"
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
	// serialNumber string
	pin string // FIXME make pin *string for better security?
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

	hsm.Ctx = ctx
	slot, err := hsm.findSlot()
	if err != nil {
		return err
	}

	sessionHandle, err := ctx.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION) // FIXME flags??
	if err != nil {
		return err
	}
	fmt.Printf("Session: %d.\n", sessionHandle)

	sessionInfo, err := ctx.GetSessionInfo(sessionHandle) // FIXME flags??
	if err != nil {
		return err
	}
	fmt.Printf("Session state: %x.\n", sessionInfo.State)
	fmt.Printf("Session state: %d.\n", sessionInfo.State)

	err = ctx.Login(sessionHandle, pkcs11.CKU_USER, hsm.hsmInfo.pin) // FIXME usertype??
	if err != nil {
		return err
	}

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

func (hsm *Hsm) findSlot() (slotID uint, err error) {
	slots, err := hsm.Ctx.GetSlotList(true)
	if err != nil {
		return 0, err
	}
	// var pkcs11.SlotInfo slotInfo
	// var pkcs11.TokenInfo tokenInfo

	for _, slot := range slots {
		if slot == 0 {
			continue
		}
		fmt.Printf("slot: %d", slot)
		fmt.Println()
		// slotInfo, err := GetSlotInfo(slot)
		tokenInfo, err := hsm.Ctx.GetTokenInfo(slot)
		if err != nil {
			return 0, err
		}
		if tokenInfo.Label == hsm.hsmInfo.tokenLabel {
			slotID = slot
			break
		}
	}
	// return 0, nil // TODO actually use tokenLabel to find the slot id
	return slotID, nil
}
