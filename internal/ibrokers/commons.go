package ibrokers

// BrokerOperationResponse to adapt response for each brokers
type BrokerOperationResponse interface {}

// Broker type
type Broker interface {
    GetRoi(order BrokerOperationResponse) (int, error)
    Buy(symbol string) (BrokerOperationResponse, error)
    Sell(symbol string) (BrokerOperationResponse, error)
}
