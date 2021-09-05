package pumpbot

import (
	"os"

	_ "github.com/joho/godotenv/autoload" 

	"github.com/adshao/go-binance/v2"
)

func main() {
	var (
		apiKey = os.Getenv("BINANCE_API_KEY")
		secretKey = os.Getenv("BINANCE_API_SECRET")
	)
	client := binance.NewClient(apiKey, secretKey)
}

func tgWatch() string {
	
	return ""
}