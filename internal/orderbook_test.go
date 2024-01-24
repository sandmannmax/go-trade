package internal

import (
	"testing"

	"github.com/phuslu/log"
	"github.com/sandmannmax/go-trade/internal/money"
)

func BenchmarkCreateLimitAsk(b *testing.B) {
	log.DefaultLogger = log.Logger{
		Level: log.ErrorLevel,
	}
	ob := NewOrderBook()
	for n := 0; n < b.N; n++ {
		ob.CreateLimitAsk(money.NewFromFloat(float64(n)), float64(n))

	}
}

func BenchmarkCreateLimitBid(b *testing.B) {
	log.DefaultLogger = log.Logger{
		Level: log.ErrorLevel,
	}
	ob := NewOrderBook()
	for n := 0; n < b.N; n++ {
		ob.CreateLimitBid(money.NewFromFloat(float64(n)), float64(n))

	}
}
