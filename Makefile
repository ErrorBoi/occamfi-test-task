run: # Run app
	go run -race cmd/price_tracker/main.go

build: # Build app
	go build -o bin/price_tracker cmd/price_tracker/main.go

lint: # Lint code using golangci-lint
	golangci-lint run -v ./...