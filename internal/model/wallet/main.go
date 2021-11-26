package wallet

// ImportedKey config
type ImportedKey struct {
	Priv string `mapstructure:"priv"`
	// Pass string `mapstructure:"pass"`
}

// Infos on wallet
type Infos struct {
	balanceAtFn func() (interface{}, error)
}

// BalanceAt adapter fn
func (i *Infos) BalanceAt() (interface{}, error) {
	return i.balanceAtFn()
}
