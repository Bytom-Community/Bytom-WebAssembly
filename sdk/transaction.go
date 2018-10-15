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
)

const getKeyByXPub = "getKeyByXPub"

type signResp struct {
	Tx           *txbuilder.Template `json:"transaction"`
	SignComplete bool                `json:"sign_complete"`
}

//Template server build struct
type Template struct {
	Transaction         string `json:"raw_transaction"`
	SigningInstructions []struct {
		DerivationPath []chainjson.HexBytes `json:"derivation_path"`
		SignData       []string             `json:"sign_data"`
	} `json:"signing_instructions"`
}

//RespSign result sign
type RespSign struct {
	Transaction string   `json:"raw_transaction"`
	Signatures  []string `json:"signatures"`
}

func signTransaction1(args []js.Value) {
	defer endFunc(args[1])
	transaction := args[0].Get("transaction").String()
	password := args[0].Get("password").String()
	keyJSON := args[0].Get("key").String()
	if isEmpty(transaction) || isEmpty(password) || isEmpty(keyJSON) {
		args[1].Set("error", "args empty")
		return
	}
	var tx Template
	err := json.Unmarshal([]byte(transaction), &tx)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	signRet := make([]string, 0, len(tx.SigningInstructions))
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
			signData, err := signServer(keyJSON, path, h, password)
			if err != nil {
				args[1].Set("error", err.Error())
				return
			}
			signRet = append(signRet, hex.EncodeToString(signData))
		}
	}
	var ret RespSign
	ret.Transaction = tx.Transaction
	ret.Signatures = signRet
	j, err := json.Marshal(ret)
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
	keys := args[0].Get("keys").String()
	if isEmpty(transaction) || isEmpty(password) || isEmpty(keys) {
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

func signServer(keyJSON string, path [][]byte, data [32]byte, password string) ([]byte, error) {
	var (
		err error
		key *pseudohsm.XKey
	)

	key, err = pseudohsm.DecryptKey([]byte(keyJSON), password)
	if err != nil {
		return nil, err
	}

	xprv := key.XPrv
	if len(path) > 0 {
		xprv = key.XPrv.Derive(path)
	}
	return xprv.Sign(data[:]), nil
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
