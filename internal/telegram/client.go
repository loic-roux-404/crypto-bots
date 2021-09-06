package telegram

import (
	"strings"
	"os"
	"time"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

func GetTgBot() *tb.Bot {
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