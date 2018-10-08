package query

import (
	"github.com/bytom-community/wasm/crypto/ed25519/chainkd"
)

//AnnotatedAccount means an annotated account.
type AnnotatedAccount struct {
	ID       string         `json:"id"`
	Alias    string         `json:"alias,omitempty"`
	XPubs    []chainkd.XPub `json:"xpubs"`
	Quorum   int            `json:"quorum"`
	KeyIndex uint64         `json:"key_index"`
}
