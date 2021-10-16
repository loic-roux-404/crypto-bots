package handlers

// BrokerOperationResponse to adapt response for each brokers
type BrokerOperationResponse interface {}

// Broker type
type Broker interface {
    GetRoi(order BrokerOperationResponse) (int, error)
    Buy(symbol string) (BrokerOperationResponse, error)
    Sell(symbol string) (BrokerOperationResponse, error)
}

var binanceBroker, _ = NewBinance()

var brokers = map[string]Broker{
    "binance": binanceBroker,
}


// GetBroker in map
func GetBroker(name string) Broker {
    return brokers[name]
}
