package helpers

import "strings"

func BtcSymbol(symbol string) string {
	return strings.Join([]string{symbol, "BTC"}, "")
}

func GrepSymbol(m string) string {
	return "SDN"
}
