package internal

import (
	"github.com/phuslu/log"
	"github.com/sandmannmax/go-trade/internal/money"
)

type order struct {
	volume float64
	isBid  bool
}

func (o order) IsFilled() bool {
	return o.volume == 0
}

func (o order) Type() string {
	if o.isBid {
		return "BID"
	}
	return "ASK"
}

func getMinMax(o1 *order, o2 *order) (*order, *order) {
	if o2.volume < o1.volume {
		return o2, o1
	}

	return o1, o2
}

type limitOrder struct {
	order
	price  money.Money
	orders []*order
}

type orderbook struct {
	askOrders map[money.Money]*limitOrder
	askVolume float64
	bidOrders map[money.Money]*limitOrder
	bidVolume float64
}

func NewOrderBook() *orderbook {
	return &orderbook{
		askOrders: make(map[money.Money]*limitOrder),
		bidOrders: make(map[money.Money]*limitOrder),
	}
}

type LimitOrderParams struct {
	price  money.Money
	volume float64
	isBid  bool
}

func (o *orderbook) CreateLimitOrder(p LimitOrderParams) {
	newOrder := &order{volume: p.volume, isBid: p.isBid}

	log.Info().Str("price", p.price.Display()).Float64("volume", newOrder.volume).Str("type", newOrder.Type()).Msg("new limit order")

	orders := o.askOrders
	volume := &o.askVolume
	if p.isBid {
		orders = o.bidOrders
		volume = &o.bidVolume
	}

	if val, ok := orders[p.price]; ok {
		val.orders = append(val.orders, newOrder)
		val.volume += newOrder.volume
		*volume += newOrder.volume
		return
	}

	orders[p.price] = &limitOrder{
		price: p.price,
		order: order{
			isBid:  p.isBid,
			volume: p.volume,
		},
		orders: []*order{newOrder},
	}
	*volume += newOrder.volume
}

func (o *orderbook) CreateLimitBid(price money.Money, volume float64) {
	o.CreateLimitOrder(LimitOrderParams{price: price, volume: volume, isBid: true})
}

func (o *orderbook) CreateLimitAsk(price money.Money, volume float64) {
	o.CreateLimitOrder(LimitOrderParams{price: price, volume: volume, isBid: false})
}

type MarketOrderParams struct {
	volume float64
	isBid  bool
}

func (o *orderbook) CreateMarketOrder(p MarketOrderParams) {
	newOrder := &order{volume: p.volume, isBid: p.isBid}
	log.Info().Float64("volume", newOrder.volume).Str("type", newOrder.Type()).Msg("new market order")

	correspondingOrders := o.bidOrders
	correspondingVolume := &o.bidVolume
	if p.isBid {
		correspondingOrders = o.askOrders
		correspondingVolume = &o.askVolume
	}

	if newOrder.volume > *correspondingVolume {
		log.Error().Float64("requested-volume", newOrder.volume).Float64("availabe-volume", *correspondingVolume).Msg("market order can not be filled - not enough volume")
		return
	}

	for price, limitOrder := range correspondingOrders {
		for _, order := range limitOrder.orders {
			minOrder, maxOrder := getMinMax(newOrder, order)
			filledVolume := minOrder.volume
			minOrder.volume = 0
			maxOrder.volume -= filledVolume
			limitOrder.volume -= filledVolume
			*correspondingVolume -= filledVolume

			if order.IsFilled() {
				log.Info().Str("price", price.Display()).Float64("volume", filledVolume).Str("type", limitOrder.Type()).Msg("limit order filled")
				if len(limitOrder.orders) == 1 {
					limitOrder.orders = nil
				} else {
					limitOrder.orders = limitOrder.orders[1:]
				}
			}

			if newOrder.IsFilled() {
				log.Info().Float64("volume", p.volume).Bool("isBid", p.isBid).Msg("market order filled")
				break
			}
		}

		if limitOrder.IsFilled() {
			delete(correspondingOrders, price)
		}

		if newOrder.IsFilled() {
			break
		}
	}

	if !newOrder.IsFilled() {
		panic("MARKET ORDER NOT FILLED - CANT HAPPEN")
	}
}

func (o *orderbook) CreateMarketBid(volume float64) {
	o.CreateMarketOrder(MarketOrderParams{volume: volume, isBid: true})
}

func (o *orderbook) CreateMarketAsk(volume float64) {
	o.CreateMarketOrder(MarketOrderParams{volume: volume, isBid: false})
}
