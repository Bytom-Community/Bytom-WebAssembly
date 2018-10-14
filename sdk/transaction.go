package sdk

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"syscall/js"

	"github.com/bytom-community/wasm/blockchain/pseudohsm"
	"github.com/bytom-community/wasm/blockchain/txbuilder"
	"github.com/bytom-community/wasm/crypto/ed25519/chainkd"
	chainjson "github.com/bytom-community/wasm/encoding/json"
	"github.com/bytom-community/wasm/protocol/bc/types"
)

const getKeyByXPub = "getKeyByXPub"

type signResp struct {
	Tx           *txbuilder.Template `json:"transaction"`
	SignComplete bool                `json:"sign_complete"`
}

//Template server build struct
type Template struct {
	Transaction         *types.Tx `json:"raw_transaction"`
	SigningInstructions []struct {
		DerivationPath []chainjson.HexBytes `json:"derivation_path"`
		SignData       []string             `json:"sign_data"`
	} `json:"signing_instructions"`
}

func signTransactionServer(args []js.Value) {
	defer endFunc(args[1])
	transaction := args[0].Get("transaction").String()
	password := args[0].Get("password").String()
	xpub := args[0].Get("xpub").String()
	if isEmpty(transaction) || isEmpty(password) || isEmpty(xpub) {
		args[1].Set("error", "args empty")
		return
	}
	var tx Template
	err := json.Unmarshal([]byte(transaction), &tx)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	signRet := make([]string, len(tx.SigningInstructions))
	for _, v := range tx.SigningInstructions {
		path := make([][]byte, len(v.DerivationPath))
		for i, p := range v.DerivationPath {
			path[i] = p
		}
		for _, d := range v.SignData {
			var h [32]byte
			t, err := hex.DecodeString(d)
			if err != nil {
				args[1].Set("error", err.Error())
				return
			}
			copy(h[:], t)
			signData, err := signServer(xpub, path, h, password)
			if err != nil {
				args[1].Set("error", err.Error())
				return
			}
			signRet = append(signRet, hex.EncodeToString(signData))
		}
	}
	j, err := json.Marshal(signRet)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	args[1].Set("data", string(j))
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
	if err := txbuilder.Sign(nil, &tx, password, sign); err != nil {
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

func signServer(xpub string, path [][]byte, data [32]byte, password string) ([]byte, error) {
	var (
		err  error
		key  *pseudohsm.XKey
		done = make(chan struct{})
	)
	js.Global().Get(getKeyByXPub).Invoke(xpub).Call("then", js.NewCallback(func(args []js.Value) {
		key, err = pseudohsm.DecryptKey([]byte(args[0].String()), password)
		done <- struct{}{}
	})).Call("catch", js.NewCallback(func(args []js.Value) {
		err = errors.New("call getKeyByXPub error maybe not found")
		done <- struct{}{}
	}))
	<-done
	if err != nil {
		return nil, err
	}

	xprv := key.XPrv
	if len(path) > 0 {
		xprv = key.XPrv.Derive(path)
	}
	return xprv.Sign(data[:]), nil
}

func sign(ctx context.Context, xpub chainkd.XPub, path [][]byte, data [32]byte, password string) ([]byte, error) {
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
	<-done
	if err != nil {
		return nil, err
	}

	xprv := key.XPrv
	if len(path) > 0 {
		xprv = key.XPrv.Derive(path)
	}
	return xprv.Sign(data[:]), nil
}
