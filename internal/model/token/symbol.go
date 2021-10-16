package token

import "strings"

// Pair type
type Pair struct{
	buy string
	sell string
}

// NewSymbol func
func NewSymbol(buy string, sell string) *Pair{
	return &Pair{buy: buy, sell: sell}
}

// NewBtcPair generate pair with BTC
func NewBtcPair(symbol string) *Pair {
	return NewSymbol(symbol, "BTC")
}

// Grep search a pair in char list
func Grep(m string) string {
	return "SDN"
}

// ToString method
func (s *Pair) ToString() string {
	return strings.Join([]string{s.buy, s.sell}, "/")
}
