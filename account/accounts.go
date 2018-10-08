// Package account stores and tracks accounts within a Bytom Core.
package account

import (
	"encoding/hex"
	"errors"

	"github.com/bytom-community/wasm/blockchain/signers"
	"github.com/bytom-community/wasm/common"
)

const (
	maxAccountCache = 1000
)

var (
	accountIndexKey     = []byte("AccountIndex")
	accountPrefix       = []byte("Account:")
	aliasPrefix         = []byte("AccountAlias:")
	contractIndexPrefix = []byte("ContractIndex")
	contractPrefix      = []byte("Contract:")
)

// pre-define errors for supporting bytom errorFormatter
var (
	ErrDuplicateAlias  = errors.New("duplicate account alias")
	ErrFindAccount     = errors.New("fail to find account")
	ErrMarshalAccount  = errors.New("failed marshal account")
	ErrInvalidAddress  = errors.New("invalid address")
	ErrFindCtrlProgram = errors.New("fail to find account control program")
)

// ContractKey account control promgram store prefix
func ContractKey(hash common.Hash) []byte {
	return append(contractPrefix, hash[:]...)
}

// ContractKeyHexString hash hex string
func ContractKeyHexString(hash common.Hash) string {
	str := hex.EncodeToString(hash[:])
	return string(append(contractPrefix, []byte(str)...))
}

// Key account store prefix
func Key(name string) []byte {
	return append(accountPrefix, []byte(name)...)
}

func aliasKey(name string) []byte {
	return append(aliasPrefix, []byte(name)...)
}

func contractIndexKey(accountID string) []byte {
	return append(contractIndexPrefix, []byte(accountID)...)
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
