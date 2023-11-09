package pubsub

import (
	"context"
	"math/rand"
	"time"
)

// A transaction provider publishing transactions over a channel
type TransactionProvider interface {
	ListenForTransactions(context.Context) <-chan Transaction
}

// Randomly generate transactions
type RandomTransactionProvider struct {
	config    *RandomTransactionProviderConfig
	dataQueue chan Transaction
}

// Config for random provider
type RandomTransactionProviderConfig struct {
	duration         time.Duration
	availableTickers []string
	queueLength      int
}
type RandomTransactionProviderOption func(*RandomTransactionProviderConfig)

func NewRandomTransactionProvider(options ...RandomTransactionProviderOption) *RandomTransactionProvider {
	config := &RandomTransactionProviderConfig{
		duration:         time.Second,
		availableTickers: []string{"GBP", "USD"},
		queueLength:      1000,
	}

	for _, option := range options {
		option(config)
	}

	return &RandomTransactionProvider{config: config, dataQueue: make(chan Transaction, config.queueLength)}
}

func WithQueueLength(length int) RandomTransactionProviderOption {
	return func(config *RandomTransactionProviderConfig) {
		config.queueLength = length
	}
}
func WithAvailableTickers(tickers []string) RandomTransactionProviderOption {
	return func(config *RandomTransactionProviderConfig) {
		config.availableTickers = tickers
	}
}
func WithDuration(duration time.Duration) RandomTransactionProviderOption {
	return func(config *RandomTransactionProviderConfig) {
		config.duration = duration
	}
}

func (p *RandomTransactionProvider) generateTransaction() Transaction {
	return Transaction{
		Format:    "16",
		Ticker:    p.config.availableTickers[rand.Intn(len(p.config.availableTickers))],
		Quantity:  rand.Int(),
		Price:     rand.Float64(),
		Timestamp: time.Now(),
	}
}

// asynchronously publish transactions in go routine. Stops when context cancelled
func (p *RandomTransactionProvider) ListenForTransactions(ctx context.Context) <-chan Transaction {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				p.dataQueue <- p.generateTransaction()
				time.Sleep(p.config.duration)
			}
		}
	}()
	return p.dataQueue
}
