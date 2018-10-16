package standard

import (
	"context"
	"encoding/json"
	"errors"
	"syscall/js"

	"github.com/bytom-community/wasm/sdk/lib"

	"github.com/bytom-community/wasm/blockchain/pseudohsm"
	"github.com/bytom-community/wasm/blockchain/txbuilder"
	"github.com/bytom-community/wasm/crypto/ed25519/chainkd"
)

type signResp struct {
	Tx           *txbuilder.Template `json:"transaction"`
	SignComplete bool                `json:"sign_complete"`
}

func getSignFunc(keys map[string]string) txbuilder.SignFunc {
	return func(ctx context.Context, xpub chainkd.XPub, path [][]byte, data [32]byte, password string) ([]byte, error) {
		var (
			err error
			key *pseudohsm.XKey
		)
		if keys, ok := keys[xpub.String()]; ok {
			key, err = pseudohsm.DecryptKey([]byte(keys), password)
			if err != nil {
				return nil, err
			}

			xprv := key.XPrv
			if len(path) > 0 {
				xprv = key.XPrv.Derive(path)
			}
			return xprv.Sign(data[:]), nil
		}

		return nil, errors.New("not found keys")
	}
}

//SignTransaction sign transaction
func SignTransaction(args []js.Value) {
	defer lib.EndFunc(args[1])
	transaction := args[0].Get("transaction").String()
	password := args[0].Get("password").String()
	keys := args[0].Get("keys").String()
	if lib.IsEmpty(transaction) || lib.IsEmpty(password) || lib.IsEmpty(keys) {
		args[1].Set("error", "args empty")
		return
	}
	var tx txbuilder.Template
	err := json.Unmarshal([]byte(transaction), &tx)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	keysMap := make(map[string]string)
	err = json.Unmarshal([]byte(keys), &keysMap)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	if err := txbuilder.Sign(nil, &tx, password, getSignFunc(keysMap)); err != nil {
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
