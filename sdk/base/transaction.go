package base

import (
	"encoding/hex"
	"encoding/json"
	"syscall/js"

	"github.com/bytom-community/wasm/sdk/lib"

	"github.com/bytom-community/wasm/blockchain/pseudohsm"
	chainjson "github.com/bytom-community/wasm/encoding/json"
)

const getKeyByXPub = "getKeyByXPub"

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
	Transaction string     `json:"raw_transaction"`
	Signatures  [][]string `json:"signatures"`
}

//SignTransaction1 sign server transaction
func SignTransaction1(args []js.Value) {
	defer lib.EndFunc(args[1])
	transaction := args[0].Get("transaction").String()
	password := args[0].Get("password").String()
	keyJSON := args[0].Get("key").String()
	if lib.IsEmpty(transaction) || lib.IsEmpty(password) || lib.IsEmpty(keyJSON) {
		args[1].Set("error", "args empty")
		return
	}
	var tx Template
	err := json.Unmarshal([]byte(transaction), &tx)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	signRet := make([][]string, len(tx.SigningInstructions))
	for k, v := range tx.SigningInstructions {
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
			if signRet[k] == nil {
				signRet[k] = make([]string, 0, len(v.SignData))
			}
			signRet[k] = append(signRet[k], hex.EncodeToString(signData))
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
