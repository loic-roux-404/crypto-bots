package kecacc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestIsValidAddress(t *testing.T) {
	t.Parallel()
	validAddress := "0x323b5d4c32345ced77393b3530b1eed0f346429d"
	invalidAddress := "0xabc"
	invalidAddress2 := "323b5d4c32345ced77393b3530b1eed0f346429d"
	{
		got := ValidateAddress(validAddress)
		expected := true

		if got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}

	{
		got := ValidateAddress(invalidAddress)
		expected := false

		if got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}

	{
		got := ValidateAddress(invalidAddress2)
		expected := false

		if got != expected {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}
}

func TestIsZeroAddress(t *testing.T) {
	t.Parallel()
	validAddress := common.HexToAddress("0x323b5d4c32345ced77393b3530b1eed0f346429d")
	zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	{
		isZeroAddress := IsZeroAddress(validAddress)

		if isZeroAddress {
			t.Error("Expected to be false")
		}
	}

	{
		isZeroAddress := IsZeroAddress(zeroAddress)

		if !isZeroAddress {
			t.Error("Expected to be true")
		}
	}

	{
		isZeroAddress := IsZeroAddress(validAddress.Hex())

		if isZeroAddress {
			t.Error("Expected to be false")
		}
	}

	{
		isZeroAddress := IsZeroAddress(zeroAddress.Hex())

		if !isZeroAddress {
			t.Error("Expected to be true")
		}
	}
}

func TestSigRSV(t *testing.T) {
	t.Parallel()

	sig := "0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301"
	r, s, v := SigRSV(sig)
	expectedR := "789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c6"
	expectedS := "2621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde023"
	expectedV := uint8(28)
	if hexutil.Encode(r[:])[2:] != expectedR {
		t.FailNow()
	}
	if hexutil.Encode(s[:])[2:] != expectedS {
		t.FailNow()
	}
	if v != expectedV {
		t.FailNow()
	}
}
