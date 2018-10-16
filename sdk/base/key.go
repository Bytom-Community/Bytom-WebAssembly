package base

import (
	"syscall/js"

	"github.com/bytom-community/wasm/sdk/lib"

	"github.com/bytom-community/wasm/blockchain/pseudohsm"

	"github.com/bytom-community/wasm/crypto/ed25519/chainkd"
	"github.com/pborman/uuid"
)

// XKey struct type for keystore file
type XKey struct {
	ID      uuid.UUID
	KeyType string
	Alias   string
	XPrv    chainkd.XPrv
	XPub    chainkd.XPub
}

//createKey create bytom key
func CreateKey(args []js.Value) {
	defer lib.EndFunc(args[1]) //end func call
	auth := args[0].Get("auth").String()
	if lib.IsEmpty(auth) {
		args[1].Set("error", "auth empty")
		return
	}
	alias := args[0].Get("alias").String()
	if lib.IsEmpty(alias) {
		args[1].Set("error", "alias empty")
		return
	}

	xprv, xpub, err := chainkd.NewXKeys(nil)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	id := uuid.NewRandom()
	key := &pseudohsm.XKey{
		ID:      id,
		KeyType: "bytom_kd",
		XPub:    xpub,
		XPrv:    xprv,
		Alias:   alias,
	}
	keyjson, err := pseudohsm.EncryptKey(key, auth, pseudohsm.LightScryptN, pseudohsm.LightScryptP)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	args[1].Set("data", string(keyjson))
}

func ResetKeyPassword(args []js.Value) {
	rootXPub := args[0].Get("rootXPub").String()
	oldPassword := args[0].Get("oldPassword").String()
	newPassword := args[0].Get("newPassword").String()
	if lib.IsEmpty(rootXPub) || lib.IsEmpty(oldPassword) || lib.IsEmpty(newPassword) {
		args[1].Set("error", "empty pm")
		lib.EndFunc(args[1])
		return
	}
	xpub := new(chainkd.XPub)
	xpub.UnmarshalText([]byte(rootXPub))
	jsv := js.Global().Get(getKeyByXPub).Invoke(xpub.String())
	var then, catch js.Callback
	then = js.NewCallback(func(a []js.Value) {
		defer then.Release()
		defer lib.EndFunc(args[1])
		key, err := pseudohsm.DecryptKey([]byte(a[0].String()), oldPassword)
		if err != nil {
			args[1].Set("error", err.Error())
			return
		}
		keyjson, err := pseudohsm.EncryptKey(key, newPassword, pseudohsm.LightScryptN, pseudohsm.LightScryptP)
		if err != nil {
			args[1].Set("error", err.Error())
			return
		}
		args[1].Set("data", string(keyjson))
	})
	catch = js.NewCallback(func(a []js.Value) {
		defer catch.Release()
		defer lib.EndFunc(args[1])
		args[1].Set("error", a[0])
	})
	jsv.Call("then", then).Call("catch", catch)
}
