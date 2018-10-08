package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"syscall/js"

	"github.com/bytom-community/wasm/blockchain/pseudohsm"

	"github.com/bytom-community/wasm/blockchain/txbuilder"
	"github.com/bytom-community/wasm/crypto/ed25519/chainkd"
)

const getKeyByXPub = "getKeyByXPub"

type signResp struct {
	Tx           *txbuilder.Template `json:"transaction"`
	SignComplete bool                `json:"sign_complete"`
}

func signTransaction(args []js.Value) {
	defer endFunc(args[1])
	transaction := args[0].Get("transaction").String()
	password := args[0].Get("password").String()
	if isEmpty(transaction) || isEmpty(password) {
		args[1].Set("error", "args empty")
		return
	}
	var tx txbuilder.Template
	err := json.Unmarshal([]byte(transaction), &tx)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	if err := txbuilder.Sign(nil, &tx, password, getSignFunc()); err != nil {
		args[1].Set("error", err.Error())
		return
	}
	sr := signResp{Tx: &tx, SignComplete: txbuilder.SignProgress(&tx)}
	resp, err := json.Marshal(sr)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	args[1].Set("data", string(resp))
}

func getSignFunc() txbuilder.SignFunc {
	return func(ctx context.Context, xpub chainkd.XPub, path [][]byte, data [32]byte, password string) ([]byte, error) {
		var (
			err  error
			key  *pseudohsm.XKey
			done = make(chan struct{})
		)

		js.Global().Get(getKeyByXPub).Invoke(xpub.String()).Call("then", js.NewCallback(func(args []js.Value) {
			key, err = pseudohsm.DecryptKey([]byte(args[0].String()), password)
			done <- struct{}{}
		})).Call("catch", js.NewCallback(func(args []js.Value) {
			err = errors.New("call getKeyByXPub error maybe not found")
			done <- struct{}{}
		}))

		if err != nil {
			return nil, err
		}

		xprv := key.XPrv
		if len(path) > 0 {
			xprv = key.XPrv.Derive(path)
		}
		return xprv.Sign(data[:]), nil
	}
}
