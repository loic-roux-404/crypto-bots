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
)

var (
	dir = helpers.GetCurrDir()
)

// Kecacc256 type
type Kecacc256 struct {
	keystore *keystore.KeyStore
	currentAccount accounts.Account
}

var (
	errPassMissing = errors.New("Missing a password")
	errAccCreation = errors.New("Error creating account")
	errAccNotFound = errors.New("No account in keystore : ")
)

// NewKecacc256 kecacc
func NewKecacc256(pass string, importKs string, fromAcc string) (*Kecacc256, error) {

	if len(pass) <= 0 {
		return nil, errPassMissing
	}

	ks := keystore.NewKeyStore(
		dir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	kecacc := &Kecacc256{keystore: ks, currentAccount: accounts.Account{}}

	err := kecacc.initAccount(pass, importKs)

	if err != nil {
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

	return kecacc, nil
}

func (k *Kecacc256) initAccount(pass string, wantedKsFile string) error {
	if len(wantedKsFile) <= 0 {
		acc, err := k.keystore.NewAccount(pass)

		if err != nil {
			return  err
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

func (k *Kecacc256) addKs(
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

	acc, err := k.keystore.NewAccount(pass); if err != nil {
		log.Panic(errAccCreation.Error(), err)
	}

	os.Rename(acc.URL.Path, filepath.Join(".", file))

	return acc, nil
}

func (k *Kecacc256) changeCurrAcc(address string) error {
	// Create account definitions
	fromAccDef := accounts.Account{
		Address: common.HexToAddress(address),
	}

	if ValidateAddress(fromAccDef) {
		return fmt.Errorf("%s : %s", ErrInvalid, address)
	}

	// Find the signing account
	signAcc, err := k.keystore.Find(fromAccDef); if err == nil {
		k.currentAccount = signAcc
	} else {
		return fmt.Errorf("%s %s", errAccNotFound, address)
	}

	return nil
}

// Account initialized
func (k *Kecacc256) Account() accounts.Account {
	return k.currentAccount
}

// Ks to use methods
func (k *Kecacc256) Ks() *keystore.KeyStore {
	return k.keystore
}
