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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/wallet"
)

var (
	dir = helpers.GetCurrDir()
)

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
	errPrivKey 	   = errors.New("Warn: error importing key starting with %s, [skipping]")
)

// NewWallet kecacc
func NewWallet(
	pass string,
	importKs string,
	fromAcc string,
	importKeys []wallet.ImportedKey,
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

	kecacc.AddPrivs(importKeys);

	err = kecacc.initAccount(pass, importKs); if err != nil {
		return nil, err
	}

	err = ks.Unlock(kecacc.currentAccount, pass)

	if err != nil {
		return nil, err
	}

	if len(fromAcc) > 0 {
		err = kecacc.changeCurrAcc(fromAcc); if err != nil {
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

	acc, err := k.addKs(wantedKsFile, pass); if err != nil {
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

		if err == keystore.ErrAccountAlreadyExists {
			log.Printf("Warning: %s", err.Error())
		}

		if err != nil && err != keystore.ErrAccountAlreadyExists  {
			println(err.Error())
			return accounts.Account{}, err
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
func (k *KeccacWallet) AddPrivs(importKeys []wallet.ImportedKey) {
	if importKeys == nil || len(importKeys) <= 0 {
		return
	}

	for _, imp := range importKeys {
		finalPriv, err := crypto.HexToECDSA(imp.Priv)
		if err != nil {
			log.Printf(errPrivKey.Error(), imp.Priv[:3])
			continue
		}

		_, err = k.keystore.ImportECDSA(finalPriv, k.pass)

		if err != nil {
			log.Printf(errPrivKey.Error(),  imp.Priv[:3])
			continue
		}
	}
}

func (k *KeccacWallet) changeCurrAcc(address string) error {
	// Create account definitions
	fromAccDef := accounts.Account{
		Address: common.HexToAddress(address),
	}

	if !ValidateAccAddress(fromAccDef) {
		return fmt.Errorf("%s : %s", ErrAccInvalid, address)
	}

	// Find the signing account
	signAcc, err := k.keystore.Find(fromAccDef); if err == nil {
		k.currentAccount = signAcc

		return nil
	}

	return fmt.Errorf("%s %s", errAccNotFound, address)
}

// IsTxFromCurrent account
func (k *KeccacWallet) IsTxFromCurrent(tx *types.Transaction) (bool, error) {
	return TxIsFrom(tx, k.currentAccount.Address)
}

// Account initialized
func (k *KeccacWallet) Account() accounts.Account {
	return k.currentAccount
}

// Ks to use methods
func (k *KeccacWallet) Ks() *keystore.KeyStore {
	return k.keystore
}
