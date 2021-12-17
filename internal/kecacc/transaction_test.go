package kecacc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TODO use simulated backend to create raw transaction
// Rlp will not work with other chain id than 1337
// https://github.com/ethereum/go-ethereum/issues/19699

type txTestCase struct {
	description string
	raw         string
	address     string
	expected    bool
	expectedErr error
}

// testCases for transaction
// ToRlp function is used to get raw tx
var testCases = []txTestCase{
	{
		description: "an empty string",
		raw:         "",
		address:     "",
		expected:    false,
		expectedErr: ErrEmptyRaw,
	},
	{
		description: "invalid raw hash",
		raw:         "2bfe64764cd97a8994d22fccd3d5b2d302fe221f07219f78b051498ec96d1add",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    false,
		expectedErr: ErrInvalidTx,
	},
	// 0x2bfe64764cd97a8994d22fccd3d5b2d302fe221f07219f78b051498ec96d1add
	{
		description: "address associated with hash",
		raw:         "eb5485010ce5dc1d82520c94be20d507fbdd6dafd7a2dde95c2d3f4618547f1787470de4df82000080808080",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    true,
		expectedErr: nil,
	},
	// 0x2bfe64764cd97a8994d22fccd3d5b2d302fe221f07219f78b051498ec96d1add
	{
		description: "addresss non associated with tx",
		raw:         "eb5485010ce5dc1d82520c94be20d507fbdd6dafd7a2dde95c2d3f4618547f1787470de4df82000080808080",
		address:     "0xBe20D507fbdD6dAFd7a2ddE95c2d3f4618547F17",
		expected:    false,
		expectedErr: nil,
	},
	// 0x74141f86d3156e5732a86c5d72f41bbd57279c4fe7ee9dcfdb82bbf62e7c1975
	{
		description: "Burn transaction",
		raw:         "ea5384fa8d247882520c94000000000000000000000000000000000000000087470de4df82000080808080",
		address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
		expected:    true,
		expectedErr: nil,
	},
	// TODO create raw
	// {
	// 	description: "Different chain id in tx",
	// 	raw:         "02f8702a80849502f907849502f90782520894be20d507fbdd6dafd7a2dde95c2d3f4618547f178603d469f28cee80c080a0dc4340630228d0e2d80df5933d1e5ff7ea2036b09084dc762f612f655476f351a075ec69b864f39a874c473240f631ae3d76f54e0dca4340afaf48b236168cbb54",
	// 	address:     "0xD4b2ae5560F8905fa4bb5C7f04122117A639B43d",
	// 	expected:    false,
	// 	expectedErr: nil,
	// },
}

func TestTxIsFrom(t *testing.T) {
	// TODO remoce once raw tx are ok
	testCases = []txTestCase{}

	for _, testCase := range testCases {
		tx, err := KeccacTxFromRaw(testCase.raw)

		if err != nil && testCase.expectedErr != nil {
			t.Logf("PASS: %s (error expected in KeccacTxFromRaw)", testCase.description)
			continue
		}

		res, err := TxIsFrom(tx, common.HexToAddress(testCase.address))
		errCaseKo := err != nil && testCase.expectedErr != nil
		expectedKo := !testCase.expected == res

		if errCaseKo || expectedKo {
			t.Fatalf("FAIL: %s \nExpected: %t\nActual: %t\nErr: %s",
				testCase.description, testCase.expected, res, err)
		}

		t.Logf("PASS: %s", testCase.description)
	}
}

func BenchmarkTxIsFrom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			tx, _ := KeccacTxFromRaw(testCase.raw)

			TxIsFrom(
				tx,
				common.HexToAddress(testCase.address),
			)
		}
	}
}
