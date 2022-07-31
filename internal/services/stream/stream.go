package stream

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ErrorBoi/occamfi-test-task/internal/app"
	"github.com/ErrorBoi/occamfi-test-task/internal/models"

	"go.uber.org/zap"
)

type Stream struct {
	Name      string
	initPrice float64
	log       *zap.SugaredLogger
}

func (s *Stream) SubscribePriceStream(
	ticker models.Ticker,
) (tickerPriceChan chan models.TickerPrice, errorChan chan error) {
	tickerPriceChan = make(chan models.TickerPrice)
	errorChan = make(chan error)

	go s.start(ticker, tickerPriceChan, errorChan)

	return tickerPriceChan, errorChan
}

func (s *Stream) start(ticker models.Ticker, tickerPriceChan chan models.TickerPrice, errorChan chan error) {
	rand.Seed(time.Now().Unix())

	for {
		//nolint:gosec // it's not necessary to generate random numbers with crypto/rand
		networkErrorOccurred := rand.Intn(101) <= 1 // 1% chance of network error occurrence
		if networkErrorOccurred {
			errorChan <- app.ErrNetwork

			s.log.Errorw("some network error occurred",
				app.ErrorTag, app.ErrNetwork.Error(),
				app.TickerTag, ticker,
				app.StreamNameTag, s.Name,
			)

			close(tickerPriceChan)
			close(errorChan)

			break
		}

		price := fmt.Sprintf("%f", s.changePrice(ticker))
		now := time.Now()

		//nolint:gosec // it's not necessary to generate random numbers with crypto/rand
		delay := 1 + rand.Intn(30) // random input delay [1;30] seconds

		<-time.After(time.Duration(delay) * time.Second)

		tickerPriceChan <- models.TickerPrice{
			Ticker: ticker,
			Time:   now,
			Price:  price,
		}

		s.log.Infow("sent ticker price mock entity",
			app.TickerTag, ticker,
			app.PriceTag, price,
			app.StreamNameTag, s.Name,
			app.TimestampTag, now.Unix(),
		)
	}
}

func (s *Stream) changePrice(ticker models.Ticker) float64 {
	oldPrice := s.initPrice

	//nolint:gosec // it's not necessary to generate random numbers with crypto/rand
	changeCoefficient := 1 + (rand.Float64()*2-1)/100 // coefficient in range [0.99; 1.01)

	s.initPrice *= changeCoefficient

	s.log.Infow("changed stream's price",
		app.TickerTag, ticker,
		app.StreamNameTag, s.Name,
		app.OldPriceTag, oldPrice,
		app.NewPriceTag, s.initPrice,
	)

	return s.initPrice
}

func NewStream(name string, initPrice float64, log *zap.SugaredLogger) *Stream {
	streamLog := log.Named("stream")

	return &Stream{
		Name:      name,
		initPrice: initPrice,
		log:       streamLog,
	}
}
