package main

import (
	"github.com/phuslu/log"
	"github.com/sandmannmax/go-trade/internal"
	"github.com/sandmannmax/go-trade/internal/money"
)

func main() {
	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
	log.Info().Msg("Moin")
	ob := internal.NewOrderBook()
	ob.CreateLimitAsk(money.NewFromFloat(100), 100)
	ob.CreateLimitAsk(money.NewFromFloat(100), 100)
	ob.CreateLimitAsk(money.NewFromFloat(110), 100)
	ob.CreateLimitAsk(money.NewFromFloat(105), 100)
	ob.CreateLimitAsk(money.NewFromFloat(107), 100)

	ob.CreateLimitBid(money.NewFromFloat(90), 50)
	ob.CreateLimitBid(money.NewFromFloat(92), 50)
	ob.CreateLimitBid(money.NewFromFloat(85), 50)
	ob.CreateLimitBid(money.NewFromFloat(87), 50)

	ob.CreateMarketBid(480)
	ob.CreateMarketBid(120)
	ob.CreateMarketBid(120)

	ob.CreateMarketAsk(120)
}
