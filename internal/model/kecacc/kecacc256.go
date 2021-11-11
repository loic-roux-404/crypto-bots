package kecacc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

var (
	dir = helpers.GetCurrDir()
)

// ImportedKey config
type ImportedKey struct {
	priv string `mapstructure:"priv"`
	pass string `mapstructure:"pass"`
}

// KeccacWallet type
type KeccacWallet struct {
	keystore       *keystore.KeyStore
	currentAccount accounts.Account
	pass           string
}

var (
	errPassMissing = errors.New("Missing a password")
	errAccCreation = errors.New("Error creating account")
	errAccNotFound = errors.New("No account in keystore : ")
)

// NewWallet kecacc
func NewWallet(
	pass string,
	importKs string,
	fromAcc string,
	importKeys []ImportedKey,
) (kecacc *KeccacWallet, err error) {
	if len(pass) <= 0 {
		return nil, errPassMissing
	}

	ks := keystore.NewKeyStore(
		dir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	kecacc = &KeccacWallet{keystore: ks, currentAccount: accounts.Account{}, pass: pass}

	kecacc.AddPrivs(importKeys)

	err = kecacc.initAccount(pass, importKs)
	if err != nil {
		return nil, err
	}

	err = ks.Unlock(kecacc.currentAccount, pass)

	if err != nil {
		return nil, err
	}

	if len(fromAcc) > 0 {
		err = kecacc.changeCurrAcc(fromAcc)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Info: Logged address %s", kecacc.currentAccount.Address)

	return kecacc, nil
}

func (k *KeccacWallet) initAccount(pass string, wantedKsFile string) error {
	if len(wantedKsFile) <= 0 {
		acc, err := k.keystore.NewAccount(pass)

		if err != nil {
			return err
		}

		k.currentAccount = acc
	}

	acc, err := k.addKs(wantedKsFile, pass)

	if err != nil {
		return err
	}

	k.currentAccount = acc

	return nil
}

func (k *KeccacWallet) addKs(
	file string,
	pass string,
) (accounts.Account, error) {
	jsonBytes, err := ioutil.ReadFile(file)

	if err == nil {
		acc, err := k.keystore.Import(jsonBytes, pass, pass)

		if err != keystore.ErrAccountAlreadyExists {
			log.Printf("Warning: %s", err)
		}

		return acc, nil
	}

	log.Printf("Warning: %s", err)
	log.Printf("Creating keystore: %s", file)

	acc, err := k.keystore.NewAccount(pass)
	if err != nil {
		log.Panic(errAccCreation.Error(), err)
	}

	os.Rename(acc.URL.Path, filepath.Join(".", file))

	return acc, nil
}

// AddPrivs key to keystore
func (k *KeccacWallet) AddPrivs(importKeys []ImportedKey) {
	if importKeys == nil || len(importKeys) <= 0 {
		return
	}

	for _, imp := range importKeys {
		finalPriv, err := crypto.HexToECDSA(imp.priv)
		if err != nil {
			log.Printf("Warn: error importing a public key, skipping...")
			continue
		}
		k.keystore.ImportECDSA(finalPriv, imp.pass)
	}
}

func (k *KeccacWallet) changeCurrAcc(address string) error {
	// Create account definitions
	fromAccDef := accounts.Account{
		Address: common.HexToAddress(address),
	}

	if ValidateAddress(fromAccDef) {
		return fmt.Errorf("%s : %s", ErrAccInvalid, address)
	}

	// Find the signing account
	signAcc, err := k.keystore.Find(fromAccDef)
	if err == nil {
		k.currentAccount = signAcc
	} else {
		return fmt.Errorf("%s %s", errAccNotFound, address)
	}

	return nil
}

// IsTxFromCurrent account
func (k *KeccacWallet) IsTxFromCurrent(hash common.Hash) (bool, error) {
	signature, err := k.keystore.SignHash(k.currentAccount, hash.Bytes())

	if err != nil {
		return false, err
	}

	return TxIsFrom(hash, signature, k.currentAccount.Address[3:])
}

// Account initialized
func (k *KeccacWallet) Account() accounts.Account {
	return k.currentAccount
}

// Ks to use methods
func (k *KeccacWallet) Ks() *keystore.KeyStore {
	return k.keystore
}
