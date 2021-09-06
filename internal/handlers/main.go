package handlers

type Broker interface {
    GetRoi(order interface{}) int
    Buy(symbol string) interface{}
    Sell(symbol string) interface{}
}

var brokers = map[string]Broker{
    "binance": NewBinance(),
}

func GetBroker(name string) Broker {
    return brokers[name]
}