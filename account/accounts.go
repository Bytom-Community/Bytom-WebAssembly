// Package account stores and tracks accounts within a Bytom Core.
package account

import (
	"encoding/hex"

	"github.com/bytom-community/wasm/blockchain/signers"
	"github.com/bytom-community/wasm/common"
)

var (
	contractPrefix = []byte("Contract:")
)

// ContractKeyHexString hash hex string
func ContractKeyHexString(hash common.Hash) string {
	str := hex.EncodeToString(hash[:])
	return string(append(contractPrefix, []byte(str)...))
}

// Account is structure of Bytom account
type Account struct {
	*signers.Signer
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

//CtrlProgram is structure of account control program
type CtrlProgram struct {
	AccountID      string
	Address        string
	KeyIndex       uint64
	ControlProgram []byte
	Change         bool // Mark whether this control program is for UTXO change
}
