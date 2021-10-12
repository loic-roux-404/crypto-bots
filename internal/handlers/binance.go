package handlers

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/loic-roux-404/pump-bot/internal/helpers"
)

// BinanceH Handler structure
type BinanceH struct {
	client *binance.Client
}

// NewBinance function
// Create a binance broker
func NewBinance() (*BinanceH, error) {
	c := binance.NewClient(
		os.Getenv("BINANCE_API_KEY"),
		os.Getenv("BINANCE_API_SECRET"),
	)

	if c == nil {
		return nil, fmt.Errorf("Binance client connection failed")
	}

	return &BinanceH{client: c}, nil
}

// Buy func
func (b *BinanceH) Buy(symbol string) (interface{}, error) {
	return orderSymbol(b.client, symbol, binance.SideTypeBuy, binance.OrderTypeMarket)
}

// Sell func
func (b *BinanceH) Sell(symbol string) (interface{}, error) {
	return orderSymbol(
		b.client,
		symbol,
		binance.SideTypeSell,
		binance.OrderTypeMarket,
	)
}

// GetRoi function 
func (*BinanceH) GetRoi(order interface{}) int {
	biOrder := order.(*binance.CreateOrderResponse)

	p, err := strconv.Atoi(biOrder.Price)
    if err != nil {
        // handle error
        fmt.Println(err)
        return -2
    }

	return p
}

func orderSymbol(
	client *binance.Client,
	symbol string,
	sideType binance.SideType,
	orderType binance.OrderType,
) (*binance.CreateOrderResponse, error) {
	order, err := client.NewCreateOrderService().
		Symbol(helpers.BtcSymbol(symbol)).
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
