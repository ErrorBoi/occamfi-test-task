package models

type PriceStreamSubscriber interface {
	SubscribePriceStream(ticker Ticker) (chan TickerPrice, chan error)
}

type MockStream struct {
	Name      string
	InitPrice float64
}
