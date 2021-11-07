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
		hash:       "",
		address:       "",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "address associated to hash",
		hash:       "",
		address:       "",
		expected:   true,
		expectedErr: true,
	},
	{
		description: "address not associated to hash",
		hash:       "",
		address:       "",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "wrong address, correct hash",
		hash:       "",
		address:       "",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "wrong hash, correct address",
		hash:       "",
		address:       "",
		expected:    false,
		expectedErr: true,
	},
	{
		description: "one is 0 address (burn)",
		hash:       "",
		address:       "",
		expected: false,
		expectedErr: true,
	},
}

func TestReverse(t *testing.T) {
	for _, testCase := range testCases {
		res, err := TxIsFrom(
			common.HexToHash(testCase.hash),
			common.HexToAddress(testCase.address),
		);
		if res != testCase.expected || err != nil && !testCase.expectedErr {
			t.Fatalf("FAIL: %s(%s,%s)\nExpected: %t\nActual: %t\nErr: %s",
				testCase.description, testCase.hash, testCase.address, testCase.expected, res, err)
		}
		t.Logf("PASS: %s", testCase.description)
	}
}


func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testCases {
			TxIsFrom(common.HexToHash(test.hash), common.HexToAddress(test.address))
		}
	}
}

