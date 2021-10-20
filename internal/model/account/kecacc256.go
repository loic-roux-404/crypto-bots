package key

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// const privKey = "PRIVATE_KEY"

// GetKecacc256 method
func GetKecacc256(memonic string, existingKeyStore string) (accounts.Account, error) {
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	
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
