package stream

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ErrorBoi/occamfi-test-task/internal/models"
)

type Timeframe struct {
	from   time.Time
	to     time.Time
	mu     sync.RWMutex
	prices map[string]float64
}

func (t *Timeframe) Add(streamName string, price models.TickerPrice) {
	if price.Time.Before(t.from) || price.Time.After(t.to) {
		return
	}

	priceFloat, err := strconv.ParseFloat(price.Price, 64)
	if err != nil {
		return
	}

	t.mu.Lock()
	t.prices[streamName] = priceFloat
	t.mu.Unlock()
}

func (t *Timeframe) Avg() string {
	t.mu.RLock()

	var result float64

	for _, v := range t.prices {
		result += v
	}

	avg := result / float64(len(t.prices))

	t.mu.RUnlock()

	return fmt.Sprintf("%.2f", avg)
}

func (t *Timeframe) Timestamp() int64 {
	return t.to.Unix()
}

func NewTimeframe(from time.Time, timeframeDuration time.Duration) *Timeframe {
	return &Timeframe{
		from:   from,
		to:     from.Add(timeframeDuration),
		prices: make(map[string]float64),
	}
}
