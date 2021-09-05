package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/adshao/go-binance/v2"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
		bi *binance.Client
		b *tb.Bot
	)

func main() {
	bi = getBinanceClient()
	b = getTgBot()
}

func getBinanceClient() *binance.Client {
	tmpClient := binance.NewClient(
		os.Getenv("BINANCE_API_KEY"), 
		os.Getenv("BINANCE_API_SECRET"),
	)

	if (tmpClient == nil) {
		log.Fatal("Binance client connection failed")
		return nil
	}

	return tmpClient
}

func getTgBot() *tb.Bot {
	b, err := tb.NewBot(tb.Settings{
		URL: strings.Join([]string{"http://", os.Getenv("TELEGRAM_BOT_ENDPOINT"), ":443"}, ""),
		Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return b
}

func msgLoop() {
	b.Handle(tb.OnChannelPost, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
		log.Printf("m: %v\n", m)
		orderBuy := orderSymbol(grepSymbol(m), binance.SideTypeBuy, binance.OrderTypeMarket)
		// ROI in percent
		roi := 0
		for {
			roi = getRoi(orderBuy.Price)
			if roi >= 100 {
				orderSymbol(grepSymbol(m), binance.SideTypeSell, binance.OrderTypeMarket)
				break
			}
		}
	})

	b.Start()
}

func grepSymbol(m *tb.Message) string {
	return "SDN"
}

func getBtcSymbol(symbol string) string {
	return strings.Join([]string{symbol, "BTC"}, "")
}

func orderSymbol(
	symbol string, 
	sideType binance.SideType,
	orderType binance.OrderType,
) *binance.CreateOrderResponse {
	order, err := bi.NewCreateOrderService().
		Symbol(getBtcSymbol(symbol)).
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

func getRoi(entryPrice string) int {
	return 0
}
