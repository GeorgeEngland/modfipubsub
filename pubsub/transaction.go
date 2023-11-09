package pubsub

import "time"

type Transaction struct {
	Format    string
	Ticker    string
	Quantity  int
	Price     float64
	Timestamp time.Time
}
