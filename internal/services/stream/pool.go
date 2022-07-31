package stream

import (
	"sync"
	"time"

	"github.com/ErrorBoi/occamfi-test-task/internal/app"
	"github.com/ErrorBoi/occamfi-test-task/internal/models"
	"go.uber.org/zap"
)

type Pool struct {
	streams      map[string]models.PriceStreamSubscriber
	mu           sync.Mutex
	ticker       models.Ticker
	shutdownChan chan struct{}
	streamChan   chan models.TickerPrice
	log          *zap.SugaredLogger
	timeframe    *Timeframe
	timeframeMu  sync.Mutex
}

func (p *Pool) GetStream() <-chan models.TickerPrice {
	return p.streamChan
}

func (p *Pool) Add(stream *Stream) {
	p.streams[stream.Name] = stream
}

func (p *Pool) Run() error {
	p.mu.Lock()

	p.timeframe = NewTimeframe(time.Now().Truncate(1*time.Minute), 1*time.Minute)

	p.initStreams()

	p.log.Infow("starting pool",
		app.TickerTag, p.ticker,
	)

	go func() {
		for {
			select {
			case <-time.After(time.Minute + time.Millisecond): // minute passed, time to create new timeframe
				var currentTimeframe *Timeframe

				p.timeframeMu.Lock()

				currentTimeframe = p.timeframe

				// create new timeframe
				now := time.Now().Truncate(time.Minute)
				p.timeframe = NewTimeframe(now, time.Minute)

				p.timeframeMu.Unlock()

				// Send data to channel
				p.streamChan <- models.TickerPrice{
					Ticker: p.ticker,
					Time:   time.Unix(currentTimeframe.Timestamp(), 0),
					Price:  currentTimeframe.Avg(),
				}
			case <-p.shutdownChan:
				break
			}
		}
	}()

	return nil
}

func (p *Pool) Shutdown() {
	p.log.Infow("stream pool graceful shutdown")

	p.mu.Unlock()
	p.shutdownChan <- struct{}{}
}

func (p *Pool) initStreams() {
	for name, sub := range p.streams {
		tickerPriceChan, errorChan := sub.SubscribePriceStream(p.ticker)

		go func(name string, tickerPriceChan <-chan models.TickerPrice, errorChan <-chan error) {
			for {
				select {
				case price := <-tickerPriceChan:
					p.timeframeMu.Lock()
					p.timeframe.Add(name, price)
					p.timeframeMu.Unlock()
				case err := <-errorChan:
					var errorStr string
					if err != nil {
						errorStr = err.Error()
					}

					p.log.Errorw("retrieving data from stream error",
						app.StreamNameTag, name,
						app.TickerTag, p.ticker,
						app.ErrorTag, errorStr,
					)

					break
				case <-p.shutdownChan:
					break
				}
			}
		}(name, tickerPriceChan, errorChan)
	}
}

func NewPool(ticker models.Ticker, log *zap.SugaredLogger) *Pool {
	poolLog := log.Named("stream pool")

	return &Pool{
		streams:      make(map[string]models.PriceStreamSubscriber),
		ticker:       ticker,
		shutdownChan: make(chan struct{}),
		streamChan:   make(chan models.TickerPrice),
		log:          poolLog,
	}
}
