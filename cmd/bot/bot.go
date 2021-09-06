package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/loic-roux-404/pump-bot/internal/handlers"
	"github.com/loic-roux-404/pump-bot/internal/telegram"
	"github.com/loic-roux-404/pump-bot/internal/helpers"
)

func main() {
	b := telegram.GetTgBot()
	cex := handlers.GetBroker("binance")
	msgLoop(b, cex)
}

func msgLoop(b *tb.Bot, broker handlers.Broker) {
	b.Handle(tb.OnChannelPost, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
		log.Printf("m: %v\n", m)
		symbol := helpers.GrepSymbol(m.Payload)
		orderBuy := broker.Buy(symbol)
		// ROI in percent
		roi := 0
		for {
			roi = broker.GetRoi(orderBuy)
			if roi >= 100 {
				broker.Sell(symbol)
				break
			}
		}
	})

	b.Start()
}
