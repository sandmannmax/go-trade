package main

import (
	"github.com/phuslu/log"
	"github.com/sandmannmax/go-trade/internal/sources"
)

type mainReceiver struct {
	Id int
}

func (r mainReceiver) Receive(data sources.SourceData) {
	log.Debug().Int("id", r.Id).Float64("price", data.Price).Float64("volume", data.Volume).Bool("isBid", data.IsBid).Bool("shouldDelete", data.ShouldDelete).Msg("received")
}

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

	// ob := internal.NewOrderBook()
	// ob.CreateLimitAsk(money.NewFromFloat(100), 100)
	// ob.CreateLimitAsk(money.NewFromFloat(100), 100)
	// ob.CreateLimitAsk(money.NewFromFloat(110), 100)
	// ob.CreateLimitAsk(money.NewFromFloat(105), 100)
	// ob.CreateLimitAsk(money.NewFromFloat(107), 100)
	//
	// ob.CreateLimitBid(money.NewFromFloat(90), 50)
	// ob.CreateLimitBid(money.NewFromFloat(92), 50)
	// ob.CreateLimitBid(money.NewFromFloat(85), 50)
	// ob.CreateLimitBid(money.NewFromFloat(87), 50)
	//
	// ob.CreateMarketBid(480)
	// ob.CreateMarketBid(120)
	// ob.CreateMarketBid(120)
	//
	// ob.CreateMarketAsk(120)

	s := sources.NewBitfinexStreamer()
	s.Subscribe(mainReceiver{Id: 1})
	// time.Sleep(1 * time.Second)
	// s.Subscribe(mainReceiver{Id: 2})
	// time.Sleep(1 * time.Second)
	// s.Subscribe(mainReceiver{Id: 3})

	select {}
}
