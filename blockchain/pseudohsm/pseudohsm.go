// Package pseudohsm provides a pseudo HSM for development environments.
package pseudohsm

import (
	"github.com/bytom-community/mobile/sdk/errors"
)

// pre-define errors for supporting bytom errorFormatter
var (
	ErrDecrypt = errors.New("could not decrypt key with given passphrase")
)
