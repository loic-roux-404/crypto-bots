package kecacc_test

import (
	"os"
	"testing"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/kecacc"
	"github.com/loic-roux-404/crypto-bots/internal/model/wallet"
	"github.com/loic-roux-404/crypto-bots/tests"
)

const KeystoreFileNotFound = "FAIL: Keystore file not found"

// TODO:
// Test created keystore file presence
// Refacto import keystore
func TestKecacc(t *testing.T) {
	acc, err := kecacc.NewWallet(
		"admin",
		tests.DummyKs,
		tests.DummyAddress,
		[]wallet.ImportedKey{
			{Priv: tests.DummyPriv},
		},
	)

	if err != nil || acc == nil {
		t.Fatalf("FAIL: error creating kecacc wallet, %p, %s", acc, err)
		return
	}

	add := acc.Account().Address.Hex()
	if add != tests.DummyAddress {
		t.Fatalf("FAIL: Address switch failed, current: %s, expected: %s", add, tests.DummyAddress)
	}

	ok, err := helpers.Exists(tests.DummyKs)
	if !ok || err != nil {
		t.Fatal(KeystoreFileNotFound)
	}

	err = Clean()
	if err != nil {
		t.Fatal(KeystoreFileNotFound)
	}
}

func Clean() error {
	e := os.Remove(tests.DummyKs)
	if e != nil {
		return e
	}

	return nil
}
