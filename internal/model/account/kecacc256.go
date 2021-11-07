package account

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

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/model/transaction"
)

var (
	dir = helpers.GetCurrDir()
)

// KeccacWallet type
type KeccacWallet struct {
	keystore       *keystore.KeyStore
	currentAccount accounts.Account
}

var (
	errPassMissing = errors.New("Missing a password")
	errAccCreation = errors.New("Error creating account")
	errAccNotFound = errors.New("No account in keystore : ")
)

// NewErcWallet kecacc
func NewErcWallet(pass string, importKs string, fromAcc string) (*KeccacWallet, error) {
	if len(pass) <= 0 {
		return nil, errPassMissing
	}

	ks := keystore.NewKeyStore(
		dir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	kecacc := &KeccacWallet{keystore: ks, currentAccount: accounts.Account{}}

	err := kecacc.initAccount(pass, importKs)

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
	return transaction.TxIsFrom(hash, k.currentAccount.Address)
}

// Account initialized
func (k *KeccacWallet) Account() accounts.Account {
	return k.currentAccount
}

// Ks to use methods
func (k *KeccacWallet) Ks() *keystore.KeyStore {
	return k.keystore
}
