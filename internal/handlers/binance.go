package handlers

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/loic-roux-404/pump-bot/internal/helpers"
)

type binanceH struct {
	client *binance.Client
}

func NewBinance() *binanceH {
	c := binance.NewClient(
		os.Getenv("BINANCE_API_KEY"),
		os.Getenv("BINANCE_API_SECRET"),
	)

	if c == nil {
		log.Fatal("Binance client connection failed")
		os.Exit(3)
	}

	return &binanceH{client: c}
}

func (b *binanceH) Buy(symbol string) interface{} {
	return orderSymbol(b.client, symbol, binance.SideTypeBuy, binance.OrderTypeMarket)
}

func (b *binanceH) Sell(symbol string) interface{} {
	return orderSymbol(
		b.client,
		symbol,
		binance.SideTypeSell,
		binance.OrderTypeMarket,
	)
}

func (*binanceH) GetRoi(order interface{}) int {
	biOrder := order.(*binance.CreateOrderResponse)

	p, err := strconv.Atoi(biOrder.Price)
    if err != nil {
        // handle error
        log.Println(err)
        return -2
    }

	return p
}

func orderSymbol(
	client *binance.Client,
	symbol string,
	sideType binance.SideType,
	orderType binance.OrderType,
) *binance.CreateOrderResponse {
	order, err := client.NewCreateOrderService().Symbol(helpers.BtcSymbol(symbol)).
		Side(sideType).
		Type(orderType).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity("5").
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return nil
	}

	log.Println(order)
	return order
}
