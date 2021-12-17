package networking

import (
	"testing"
)

type urlTestCase struct {
	description string
	url         string
	shouldBeWs  bool
	shouldBeRPC bool
	// expectedErr error
}

// testCases for transaction
// ToRlp function is used to get raw tx
var testCases = []urlTestCase{
	{description: "empty str", url: "", shouldBeWs: false},
	{description: "non protocol", url: "ws", shouldBeWs: false},
	{description: "empty ws domain", url: "ws://", shouldBeWs: false},
	{description: "ws with ssl and route", url: "wss://eth-ropsten.alchemyapi.io/v2/ggBP3Q83g-AC_f2sy7tz_7bjkT1ISSma", shouldBeWs: true},
	{description: "simple ws", url: "wss://eth-ropsten.alchemyapi.io", shouldBeWs: true},
	{description: "rpc compatible", url: "http://eth-ropsten.alchemyapi.io", shouldBeWs: true},
	{description: "rpc direct protocol", url: "rpc://eth-ropsten.alchemyapi.io", shouldBeWs: true},
}

func TestIsWs(t *testing.T) {
	for _, testCase := range testCases {
		res := IsWs(testCase.url)

		if !res && testCase.shouldBeWs {
			t.Fatalf("FAIL: %s \nExpected: %t\nActual: %t\nErr: %s",
				testCase.description, testCase.shouldBeWs, res, "")
		}

		t.Logf("PASS: %s", testCase.description)
	}
}

func BenchmarkTxIsFrom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, testCase := range testCases {
			IsWs(testCase.url)
		}
	}
}
