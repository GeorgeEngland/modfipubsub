package pubsub

import (
	"context"
	"fmt"
)

type Producer interface {
	Run(context.Context)
	publishTransaction(m Transaction)
}

type InMemoryProducer struct {
	config *InMemoryProducerConfig
}
type InMemoryProducerConfig struct {
	consumers           []Consumer
	transactionProvider TransactionProvider
}
type Option func(*InMemoryProducerConfig)

func NewInMemoryProducer(options ...Option) *InMemoryProducer {

	config := &InMemoryProducerConfig{
		consumers:           []Consumer{},
		transactionProvider: NewRandomTransactionProvider(),
	}
	for _, option := range options {
		option(config)
	}

	return &InMemoryProducer{config: config}
}

func WithListeners(consumers ...Consumer) Option {
	return func(config *InMemoryProducerConfig) {
		config.consumers = append(config.consumers, consumers...)
	}
}

func WithTransactionProvider(provider TransactionProvider) Option {
	return func(config *InMemoryProducerConfig) {
		config.transactionProvider = provider
	}
}

func (p *InMemoryProducer) Run(ctx context.Context) error {
	ctxTx, cancel := context.WithCancel(ctx)
	defer cancel()
	c := p.config.transactionProvider.ListenForTransactions(ctxTx)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Producer Done")
			return nil
		case m := <-c:
			p.publishTransaction(m)
		}
	}
}

func (p *InMemoryProducer) publishTransaction(m Transaction) {
	for _, consumer := range p.config.consumers {
		consumer.GetChannel() <- m
	}
}
