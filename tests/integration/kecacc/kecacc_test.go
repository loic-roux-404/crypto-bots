package kecacc_test

import (
	"testing"

	"github.com/loic-roux-404/crypto-bots/internal/helpers"
	"github.com/loic-roux-404/crypto-bots/internal/kecacc"
	"github.com/loic-roux-404/crypto-bots/internal/model/wallet"
	"github.com/loic-roux-404/crypto-bots/tests"
)

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
		t.Logf("FAIL: error creating kecacc wallet, %p, %s", acc, err)
	}

	add := acc.Account().Address.Hex()
	if add != tests.DummyAddress {
		t.Logf("FAIL: Address switch failed, current: %s, expected: %s", add, tests.DummyAddress)
	}

	ok, err := helpers.Exists(tests.DummyKs)
	if !ok || err != nil {
		t.Log("FAIL: Keystore file not found")
	}
}
