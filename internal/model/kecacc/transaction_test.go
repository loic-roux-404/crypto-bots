package kecacc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	// "testing/quick"
)

type txTestCase struct {
	description string
	hash        string
	signature   string
	address     string
	expected    bool
	expectedErr bool
}

var testCases = []txTestCase{
	{
		description: "an empty string",
		hash:        "",
		signature:   "",
		address:     "",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "address associated with hash",
		hash:        "0x2bfe64764cd97a8994d22fccd3d5b2d302fe221f07219f78b051498ec96d1add",
		signature:   "0xf0b31f19f4c4a51aa98420dac2e590c3c40284e07ee2e6827b917ac58b46589b2ee587301a95757535d463e4c18f0717a0891ed7ec146bfcd739e7063c9dcdd501",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    true,
		expectedErr: false,
	},
	{
		description: "address not associated to hash",
		hash:        "0xebaab87b41ac1e5e2ff92ada0dcd8bca9b679498139ee1508efe2c87cce7f3c6",
		signature:   "0xf0b31f19f4c4a51aa98420dac2e590c3c40284e07ee2e6827b917ac58b46589b2ee587301a95757535d463e4c18f0717a0891ed7ec146bfcd739e7063c9dcdd501",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "tx non associated with signature",
		hash:        "0x2bfe64764cd97a8994d22fccd3d5b2d302fe221f07219f78b051498ec96d1add",
		signature:   "0x12596324a64e419e19db71b17d34dbefd73de4d7aadfb9beb3b8048e7f9402be3f45ba18f0cc6fb087a4e6718c5ec25343f5a40db43ecbe7b12fc34363c2439001",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    false,
		expectedErr: false,
	},
	{
		description: "Burn transaction",
		hash:        "0x86c040662114d7d1116c26978b90d6cf3752256be46b13b0eefd882ec4d88f4d",
		signature:   "0x884a0d10b9811d1812a80a010b8d6e0b895dc2f4595635d8e1a024115fb62d6d0275a1e1a67711ee32b13ca1fda180a979d3baa4b9e0e735c1333bc087e37fc300",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    true,
		expectedErr: false,
	},
	// {
	// 	description: "Non matching chain id in tx",
	// 	hash:        "0xa1032916b2be75d6d2a7d0b2ca2ae4b24b38b28bcd4648126d0c44e080d10014",
	// 	signature:   "",
	// 	address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
	// 	expected:    false,
	// 	expectedErr: true,
	// },
}

func TestTxIsFrom(t *testing.T) {
	for _, testCase := range testCases {
		res, err := TxIsFrom(
			common.HexToHash(testCase.hash),
			[]byte(testCase.signature),
			common.HexToAddress(testCase.address).Bytes()[3:],
		)
		errCaseKo := err != nil && !testCase.expectedErr
		expectedKo := !testCase.expected == res

		if errCaseKo || expectedKo {
			t.Fatalf("FAIL: %s (%s,%s)\nExpected: %t\nActual: %t\nErr: %s",
				testCase.description, testCase.hash, testCase.address, testCase.expected, res, err)
		}
		t.Logf("PASS: %s", testCase.description)
	}
}

func BenchmarkTxIsFrom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			TxIsFrom(
				common.HexToHash(testCase.hash),
				[]byte(testCase.signature),
				common.HexToAddress(testCase.address).Bytes()[3:],
			)
		}
	}
}
