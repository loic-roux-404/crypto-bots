package account

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// const privKey = "PRIVATE_KEY"

var (
		dir, _ = os.Getwd()
		ks = keystore.NewKeyStore(
		dir,
		keystore.StandardScryptN, 
		keystore.StandardScryptP,
	)
)

// Kecacc256 type
type Kecacc256 struct {
	store *keystore.KeyStore
	currentAccount accounts.Account
}

// NewKecacc256 kecacc
func NewKecacc256(memonic string, existingKeyStore string) (*Kecacc256, error) {
	acc, err := initAccount(memonic, existingKeyStore)
	
	if err != nil {
		return nil, err
	}

	return &Kecacc256{
		store: ks,
		currentAccount: acc,
	}, nil
}

func initAccount(memonic string, existingKeyStore string) (accounts.Account, error) {
	if existingKeyStore == "" {
		return ks.NewAccount(memonic)
	}

	exStore, err := ioutil.ReadFile(existingKeyStore)

	if (err != nil || len(exStore) > 0) {
		return importKeyStore(ks, existingKeyStore, memonic)
	}

	return ks.NewAccount(memonic)
}

func importKeyStore(
	ks *keystore.KeyStore,
	file string,
	memonic string,
) (accounts.Account, error) {
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}

	return ks.Import(jsonBytes, memonic, memonic)
}

// Account initialized
func (k *Kecacc256) Account() accounts.Account {
	return k.currentAccount
}

// Store to use methods
func (k *Kecacc256) Store() *keystore.KeyStore {
	return k.store
}
