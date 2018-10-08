package pseudohsm

import (
	"github.com/bytom-community/wasm/crypto/ed25519/chainkd"
	"github.com/pborman/uuid"
)

const (
	version = 1
	keytype = "bytom_kd"
)

// XKey struct type for keystore file
type XKey struct {
	ID      uuid.UUID
	KeyType string
	Alias   string
	XPrv    chainkd.XPrv
	XPub    chainkd.XPub
}

type encryptedKeyJSON struct {
	Crypto  cryptoJSON `json:"crypto"`
	ID      string     `json:"id"`
	Type    string     `json:"type"`
	Version int        `json:"version"`
	Alias   string     `json:"alias"`
	XPub    string     `json:"xpub"`
}

type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}
