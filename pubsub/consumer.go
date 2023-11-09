package pubsub

import (
	"context"
	"fmt"
)

type Consumer interface {
	Handle(Transaction) error
	GetChannel() chan Transaction
}

func Run(ctx context.Context, cons ...Consumer) {
	for _, c := range cons {
		ch := c.GetChannel()
		go func(c Consumer) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Consumer Done")
					return
				case m := <-ch:
					err := c.Handle(m)
					if err != nil {
						fmt.Println(err)
					}
				}

			}
		}(c)
	}
}

type MaxMinConsumer struct {
	maxes map[string]float64
	mins  map[string]float64
	ch    chan Transaction
}

func NewMaxMinConsumer() *MaxMinConsumer {
	return &MaxMinConsumer{
		maxes: make(map[string]float64),
		mins:  make(map[string]float64),
		ch:    make(chan Transaction),
	}
}

func (m *MaxMinConsumer) Handle(t Transaction) error {
	_, exists := m.maxes[t.Ticker]
	if !exists {
		m.maxes[t.Ticker] = t.Price
		m.mins[t.Ticker] = t.Price
	} else {
		if t.Price > m.maxes[t.Ticker] {
			m.maxes[t.Ticker] = t.Price
		}
		if t.Price < m.mins[t.Ticker] {
			m.mins[t.Ticker] = t.Price
		}
	}
	fmt.Println("MAX:", t.Ticker, m.maxes[t.Ticker])
	fmt.Println("MIN:", t.Ticker, m.mins[t.Ticker])
	return nil
}

func (m *MaxMinConsumer) GetChannel() chan Transaction {
	return m.ch
}

type LogConsumer struct {
	ch chan Transaction
}

func NewLogConsumer() *LogConsumer {
	return &LogConsumer{
		ch: make(chan Transaction),
	}
}

func (m *LogConsumer) Handle(t Transaction) error {
	fmt.Println("GOT TRANS", t)
	return nil
}

func (m *LogConsumer) GetChannel() chan Transaction {
	return m.ch
}
