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
func NewKecacc256(memonic string, wantedAcc string) (*Kecacc256, error) {
	if len(memonic) <= 0 {
		return nil, fmt.Errorf("Missing a memonic")
	}

	ks := keystore.NewKeyStore(
		dir,
		keystore.StandardScryptN,
		keystore.StandardScryptP,
	)

	acc, err := initAccount(ks, memonic, wantedAcc)

	if err != nil {
		return nil, err
	}

	err = ks.Unlock(acc, memonic)

	if err != nil {
		return nil, err
	}

	return &Kecacc256{
		store: ks,
		currentAccount: acc,
	}, nil
}

func initAccount(
	ks *keystore.KeyStore,
	memonic string,
	wantedAcc string,
) (accounts.Account, error) {
	if wantedAcc == "" {
		return ks.NewAccount(memonic)
	}

	return importKeyStore(ks, wantedAcc, memonic)
}

func importKeyStore(
	ks *keystore.KeyStore,
	file string,
	memonic string,
) (accounts.Account, error) {
	jsonBytes, err := ioutil.ReadFile(file)

	if err == nil {
		acc, err := ks.Import(jsonBytes, memonic, memonic); if err != nil {
			log.Printf("Warn: %s", err)
		}

		return acc, nil
	}

	log.Printf("Warn: %s", err)
	log.Printf("Creating keystore: %s", file)

	acc, err := ks.NewAccount(memonic); if err != nil {
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
