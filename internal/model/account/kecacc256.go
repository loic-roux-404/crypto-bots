package account

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
)

var (
	dir = helpers.GetCurrDir()
)

// Kecacc256 type
type Kecacc256 struct {
	store *keystore.KeyStore
	currentAccount accounts.Account
}

// NewKecacc256 kecacc
func NewKecacc256(pass string, wantedAcc string) (*Kecacc256, error) {
	if len(pass) <= 0 {
		return nil, fmt.Errorf("Missing a password")
	}

	ks := keystore.NewKeyStore(
		dir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	acc, err := initAccount(ks, pass, wantedAcc)

	if err != nil {
		return nil, err
	}

	err = ks.Unlock(acc, pass)

	if err != nil {
		return nil, err
	}

	return &Kecacc256{store: ks, currentAccount: acc}, nil
}

func initAccount(
	ks *keystore.KeyStore,
	pass string,
	wantedAcc string,
) (accounts.Account, error) {
	if wantedAcc == "" {
		return ks.NewAccount(pass)
	}

	return importKeyStore(ks, wantedAcc, pass)
}

func importKeyStore(
	ks *keystore.KeyStore,
	file string,
	pass string,
) (accounts.Account, error) {
	jsonBytes, err := ioutil.ReadFile(file)

	if err == nil {
		acc, err := ks.Import(jsonBytes, pass, pass); if err != nil {
			log.Printf("Warning: %s", err)
		}

		return acc, nil
	}

	log.Printf("Warning: %s", err)
	log.Printf("Creating keystore: %s", file)

	acc, err := ks.NewAccount(pass); if err != nil {
		log.Panicf("Error creating account: %s", err)
	}

	os.Rename(acc.URL.Path, filepath.Join(".", file))

	return acc, nil
}

// Account initialized
func (k *Kecacc256) Account() accounts.Account {
	return k.currentAccount
}

// Store to use methods
func (k *Kecacc256) Store() *keystore.KeyStore {
	return k.store
}
