package kecacc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	// "testing/quick"
)

type txTestCase struct {
	description string
	hash        string
	address     string
	expected    bool
	expectedErr bool
}

var testCases = []txTestCase{
	{
		description: "an empty string",
		hash:        "",
		address:     "",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "address associated to hash",
		hash:        "0x920af48a13855d090cea34915f7176a8db93aaa2e3762b42a1708764bfe1046e",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    true,
		expectedErr: false,
	},
	{
		description: "address not associated to hash",
		hash:        "0x00c71f8049c9c5588f17388a7651a5df1730142245fef33f34c97d380f23b30c",
		address:     "0xE216378C0ed702D66e09D4aDBE9548C52604eB6E",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "Not matching hash / address",
		hash:        "0x525d85bee80ca76a50acbb930f1cfa006e9436415d4f39eb28367ede4fbb3bc4",
		address:     "0xE216378C0ed702D66e09D4aDBE9548C52604eB6E",
		expected:    false,
		expectedErr: false,
	},
	{
		description: "Burn transaction",
		hash:        "0x0b15f7c6c13e17fd0dcb05116be24964fb0e2e8d96fd0eaa5d3e686476f75692",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    true,
		expectedErr: false,
	},
	{
		description: "Non matching chain id in tx",
		hash:        "",
		address:     "",
		expected:    false,
		expectedErr: true,
	},
}

func TestTransaction(t *testing.T) {
	for _, testCase := range testCases {
		res, err := TxIsFrom(
			common.HexToHash(testCase.hash),
			common.HexToAddress(testCase.address),
		)
		errCaseKo := err != nil && !testCase.expectedErr
		expectedKo := !testCase.expected == res

		if errCaseKo || expectedKo {
			t.Fatalf("FAIL: %s(%s,%s)\nExpected: %t\nActual: %t\nErr: %s",
				testCase.description, testCase.hash, testCase.address, testCase.expected, res, err)
		}
		t.Logf("PASS: %s", testCase.description)
	}
}

func BenchmarkTransaction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testCases {
			TxIsFrom(common.HexToHash(test.hash), common.HexToAddress(test.address))
		}
	}
}
