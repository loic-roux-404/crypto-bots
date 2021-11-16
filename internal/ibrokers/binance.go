package ibrokers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
)

// BinanceHandler Handler structure
type BinanceHandler struct {
	client  *binance.Client
	btcPair bool
}

// NewBinance function
// Create a binance broker
// return *BinanceHandler
func NewBinance() Broker {
	c := binance.NewClient(
		os.Getenv("BINANCE_API_KEY"),
		os.Getenv("BINANCE_API_SECRET"),
	)

	if c == nil {
		log.Panic("Binance client connection failed")
	}

	return &BinanceHandler{client: c, btcPair: true}
}

// Buy func
func (b *BinanceHandler) Buy(symbol string) (BrokerOperationResponse, error) {
	// TODO comply with an common interface and manage orderSymbol return
	// TODO get pair by config
	tokenPair := token.NewBtcPair(symbol)
	return orderSymbol(b.client, tokenPair, binance.SideTypeBuy, binance.OrderTypeMarket)
}

// Sell func
func (b *BinanceHandler) Sell(symbol string) (BrokerOperationResponse, error) {
	// TODO comply with a common interface and manage orderSymbol return
	// TODO get pair by config
	tokenPair := token.NewBtcPair(symbol)
	return orderSymbol(
		b.client,
		tokenPair,
		binance.SideTypeSell,
		binance.OrderTypeMarket,
	)
}

// GetRoi function
func (*BinanceHandler) GetRoi(order BrokerOperationResponse) (int, error) {
	biOrder := order.(*binance.CreateOrderResponse)

	p, err := strconv.Atoi(biOrder.Price)
	if err != nil {
		// handle error
		// fmt.Println(err)
		return -2, err
	}

	return p, err
}

func orderSymbol(
	client *binance.Client,
	token token.Pair,
	sideType binance.SideType,
	orderType binance.OrderType,
) (*binance.CreateOrderResponse, error) {
	order, err := client.NewCreateOrderService().
		Symbol(token.ToString()).
		Side(sideType).
		Type(orderType).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity("5"). // TODO
		Do(context.Background())

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	fmt.Println(order)

	return order, nil
}
