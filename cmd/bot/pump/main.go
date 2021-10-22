package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/loic-roux-404/crypto-bots/pkg/brokers"
	"github.com/loic-roux-404/crypto-bots/internal/ibrokers"
	"github.com/loic-roux-404/crypto-bots/internal/telegram"
	"github.com/loic-roux-404/crypto-bots/internal/model/token"
)

func main() {
	b := telegram.GetTgBot()
	cex, err := brokers.Get("binance")

	if (err != nil) {
		log.Fatal(err)
	}

	msgLoop(b, cex)
}

func msgLoop(b *tb.Bot, broker ibrokers.Broker) {
	b.Handle(tb.OnChannelPost, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
		log.Printf("m: %v\n", m)
		symbol := token.Grep(m.Payload)
		orderBuy, err := broker.Buy(symbol)
		log.Fatal(err)
		// ROI in percent
		roi := 0

		for {
			roi, err = broker.GetRoi(orderBuy)
			log.Panic(err)
			if roi >= 100 {
				broker.Sell(symbol)
				break
			}
		}
	})

	b.Start()
}
