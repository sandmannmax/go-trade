package sources

import (
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/phuslu/log"
)

type SourceData struct {
	Price        float64
	Volume       float64
	IsBid        bool
	ShouldDelete bool
}

type Streamer interface {
	Subscribe(r Receiver)
}

type Receiver interface {
	Receive(SourceData)
}

type bitfinexStreamer struct {
	receiverChan chan Receiver
	dataChan     chan SourceData
}

func NewBitfinexStreamer() *bitfinexStreamer {
	s := &bitfinexStreamer{
		receiverChan: make(chan Receiver),
		dataChan:     make(chan SourceData),
	}
	go s.run()
	go s.coordinate()
	return s
}

func (s *bitfinexStreamer) Subscribe(r Receiver) {
	s.receiverChan <- r
}

func (s bitfinexStreamer) run() {
	c, _, err := websocket.DefaultDialer.Dial("wss://api-pub.bitfinex.com/ws/2", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte("{\"event\":\"subscribe\",\"channel\":\"book\",\"symbol\":\"tBTCUSD\"}"))
	if err != nil {
		log.Error().Err(err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Err(err)
			continue
		}

		m := string(message)
		if m[0] == '{' {
			continue
		}
		splitMessage := strings.Split(m, ",")[1:]

		var (
			price  int
			count  int
			amount float64
		)
		for index, messagePart := range splitMessage {
			if index%3 == 0 {
				price, _ = strconv.Atoi(strings.ReplaceAll(messagePart, "[", ""))
			} else if index%3 == 1 {
				count, _ = strconv.Atoi(messagePart)
			} else {
				amount, _ = strconv.ParseFloat(strings.ReplaceAll(messagePart, "]", ""), 64)

				data := SourceData{
					Price:        float64(price),
					Volume:       float64(amount),
					IsBid:        amount > 0,
					ShouldDelete: count == 0,
				}

				if !data.IsBid {
					data.Volume = -data.Volume
				}

				s.dataChan <- data
			}
		}
	}
}

func (s bitfinexStreamer) coordinate() {
	receivers := []Receiver{}

	for {
		select {
		case r := <-s.receiverChan:
			receivers = append(receivers, r)
		case d := <-s.dataChan:
			for _, r := range receivers {
				r.Receive(d)
			}
		}
	}
}
