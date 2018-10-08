package sdk

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"github.com/bytom-community/wasm/account"
	"github.com/bytom-community/wasm/blockchain/signers"
	"github.com/bytom-community/wasm/blockchain/txbuilder"
	"github.com/bytom-community/wasm/common"
	"github.com/bytom-community/wasm/consensus"
	"github.com/bytom-community/wasm/crypto"
	"github.com/bytom-community/wasm/crypto/ed25519/chainkd"
	"github.com/bytom-community/wasm/crypto/sha3pool"
	"github.com/bytom-community/wasm/protocol/vm/vmutil"
)

func createAccount(args []js.Value) {
	defer endFunc(args[1])
	var (
		alias     string
		quorum    int
		rootXPub  string
		nextIndex uint64 // account next index like 1 2 3 ...
	)

	alias = args[0].Get("alias").String()
	quorum = args[0].Get("quorum").Int()
	rootXPub = args[0].Get("rootXPub").String()
	nextIndex = uint64(args[0].Get("nextIndex").Int())

	var XPubs []chainkd.XPub
	xpub := new(chainkd.XPub)
	xpub.UnmarshalText([]byte(rootXPub))
	XPubs = append(XPubs, *xpub)

	normalizedAlias := strings.ToLower(strings.TrimSpace(alias))

	signer, err := signers.Create("account", XPubs, quorum, nextIndex)
	id := signers.IDGenerate()
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}

	acc := &account.Account{Signer: signer, ID: id, Alias: normalizedAlias}
	rawAccount, err := json.Marshal(acc)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	args[1].Set("data", string(rawAccount))
}

func createAccountReceiver(args []js.Value) {
	defer endFunc(args[1])
	var (
		acc       account.Account
		err       error
		nextIndex uint64
		cp        *account.CtrlProgram
	)
	err = json.Unmarshal([]byte(args[0].Get("account").String()), &acc)
	nextIndex = uint64(args[0].Get("nextIndex").Int())

	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	if len(acc.XPubs) == 1 {
		cp, err = createP2PKH(&acc, false, nextIndex)
	} else {
		cp, err = createP2SH(&acc, false, nextIndex)
	}
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}

	res, err := controlPrograms(cp)
	if err != nil {
		args[1].Set("error", err.Error())
		return
	}
	data, _ := json.Marshal(res)
	args[1].Set("db", string(data)) //insert web IndexedDB
	tx := txbuilder.Receiver{
		ControlProgram: cp.ControlProgram,
		Address:        cp.Address,
	}
	txj, _ := json.Marshal(tx)
	args[1].Set("data", string(txj))
}

func createP2PKH(acc *account.Account, change bool, nextIndex uint64) (*account.CtrlProgram, error) {
	path := signers.Path(acc.Signer, signers.AccountKeySpace, nextIndex)
	derivedXPubs := chainkd.DeriveXPubs(acc.XPubs, path)
	derivedPK := derivedXPubs[0].PublicKey()
	pubHash := crypto.Ripemd160(derivedPK)

	address, err := common.NewAddressWitnessPubKeyHash(pubHash, &consensus.ActiveNetParams)
	if err != nil {
		return nil, err
	}

	control, err := vmutil.P2WPKHProgram([]byte(pubHash))
	if err != nil {
		return nil, err
	}

	return &account.CtrlProgram{
		AccountID:      acc.ID,
		Address:        address.EncodeAddress(),
		KeyIndex:       nextIndex,
		ControlProgram: control,
		Change:         change,
	}, nil
}

func createP2SH(acc *account.Account, change bool, nextIndex uint64) (*account.CtrlProgram, error) {
	path := signers.Path(acc.Signer, signers.AccountKeySpace, nextIndex)
	derivedXPubs := chainkd.DeriveXPubs(acc.XPubs, path)
	derivedPKs := chainkd.XPubKeys(derivedXPubs)
	signScript, err := vmutil.P2SPMultiSigProgram(derivedPKs, acc.Quorum)
	if err != nil {
		return nil, err
	}
	scriptHash := crypto.Sha256(signScript)

	address, err := common.NewAddressWitnessScriptHash(scriptHash, &consensus.ActiveNetParams)
	if err != nil {
		return nil, err
	}

	control, err := vmutil.P2WSHProgram(scriptHash)
	if err != nil {
		return nil, err
	}

	return &account.CtrlProgram{
		AccountID:      acc.ID,
		Address:        address.EncodeAddress(),
		KeyIndex:       nextIndex,
		ControlProgram: control,
		Change:         change,
	}, nil
}

func controlPrograms(progs ...*account.CtrlProgram) (map[string]string, error) {
	var hash common.Hash
	res := make(map[string]string)
	for _, prog := range progs {
		accountCP, err := json.Marshal(prog)
		if err != nil {
			return nil, err
		}

		sha3pool.Sum256(hash[:], prog.ControlProgram)
		res[account.ContractKeyHexString(hash)] = string(accountCP)
	}
	return res, nil
}
