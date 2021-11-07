package telegram

import (
	"log"
	"os"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func GetTgBot() *tb.Bot {
	// ping
	PORT := "443"
	URL := strings.Join([]string{"http://", os.Getenv("TELEGRAM_BOT_ENDPOINT")}, "")
	finalUrl := strings.Join([]string{URL, ":", PORT}, "")

	b, err := tb.NewBot(tb.Settings{
		URL:    finalUrl,
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal("Can't connect to telegram endpoint ", err)
		return nil
	}

	return b
}
