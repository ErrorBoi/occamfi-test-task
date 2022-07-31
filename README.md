# Occam.fi test task

## Run
- clone repo
- execute `go mod download`
- execute `make run`

## Algorithm
- Register input Streams 
  - I use set of predefined mocked Streams for testing
  - You can find them in `cmd/price_tracker/main.go`
  - Mock Streams mutate their asset price after random timeout
  - It makes mock inputs more dynamic and unpredictable
- Collect these input Streams in Stream Pool
- Read all messages from each Stream
- Map messages to timeframes depending on Time field
  - Timeframe is a bunch of prices from input messages, which were delivered during specific minute
- Compute average asset price during this timeframe
- Output average price to standard output

## Features
- Graceful Shutdown
- Race detector
- Linter (golangci-lint)

## Task description
[Click](task.md)