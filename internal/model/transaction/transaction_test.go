package transaction

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
		expectedErr: true,
	},
	{
		description: "address not associated to hash",
		hash:        "",
		address:     "",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "Not matching hash / address",
		hash:        "0x920af48a13855d090cea34915f7176a8db93aaa2e3762b42a1708764bfe1046e",
		address:     "0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17",
		expected:    false,
		expectedErr: false,
	},
	{
		description: "one is 0 address (burn)",
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
		if res != testCase.expected || err != nil && !testCase.expectedErr {
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
