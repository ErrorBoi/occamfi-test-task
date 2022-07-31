package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ErrorBoi/occamfi-test-task/internal/app"
	"github.com/ErrorBoi/occamfi-test-task/internal/models"
	"github.com/ErrorBoi/occamfi-test-task/internal/services/stream"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	sugar := logger.Sugar()

	streams := []models.MockStream{
		{Name: "Binance", InitPrice: 23800.00},
		{Name: "OKX", InitPrice: 23790.00},
		{Name: "FTX", InitPrice: 23825.00},
		{Name: "Mercuryo", InitPrice: 23290.00},
		{Name: "KuCoin", InitPrice: 23800.00},
	}

	streamPool := stream.NewPool(models.BTCUSDTicker, sugar)

	for _, s := range streams {
		source := stream.NewStream(s.Name, s.InitPrice, sugar)

		streamPool.Add(source)
	}

	go func() {
		if err := streamPool.Run(); err != nil {
			sugar.Fatalw("error running stream pool")
		}
	}()

	// Receive streamed data from pool
	go func() {
		for value := range streamPool.GetStream() {
			fmt.Printf("Timestamp: %d | Index Price: %s\n", value.Time.Unix(), value.Price)
			sugar.Infow("new avg value",
				app.PriceTag, value.Price,
				app.TickerTag, value.Ticker,
				app.TimestampTag, value.Time.Unix())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	streamPool.Shutdown()
}
