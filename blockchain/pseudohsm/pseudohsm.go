// Package pseudohsm provides a pseudo HSM for development environments.
package pseudohsm

import (
	"github.com/bytom-community/mobile/sdk/errors"
)

// pre-define errors for supporting bytom errorFormatter
var (
	ErrDuplicateKeyAlias    = errors.New("duplicate key alias")
	ErrInvalidAfter         = errors.New("invalid after")
	ErrLoadKey              = errors.New("key not found or wrong password ")
	ErrTooManyAliasesToList = errors.New("requested aliases exceeds limit")
	ErrDecrypt              = errors.New("could not decrypt key with given passphrase")
)
