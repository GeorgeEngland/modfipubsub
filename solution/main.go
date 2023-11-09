package main

import (
	"context"
	"fmt"
	"modfi/pubsub"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Allow Exiting of Service by cancelling ctx with Ctrl C
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)

	go func() {
		// Wait for a Ctrl+C signal
		<-interrupt
		fmt.Println("Ctrl+C received. Cancelling context...")
		cancel()
	}()

	// Transactions provided by interface implemented by a random generator (but could be websocket etc.)
	transactionProvider := pubsub.NewRandomTransactionProvider(
		pubsub.WithAvailableTickers([]string{"GBP", "EUR"}),
		pubsub.WithDuration(time.Second),
		pubsub.WithQueueLength(100),
	)

	// Create the Consumers
	cs1 := pubsub.NewLogConsumer()
	cs2 := pubsub.NewMaxMinConsumer()

	// Create the producer with the random transaction provider
	// Create the producer with the given Consumers
	producer := pubsub.NewInMemoryProducer(
		pubsub.WithTransactionProvider(transactionProvider),
		pubsub.WithListeners(cs1, cs2))

	// Start the producer in the background
	go func() { producer.Run(ctx) }()

	// Start the consumers
	pubsub.Run(ctx, cs1, cs2)

	// Wait for Ctrl C
	<-ctx.Done()
	// The context was cancelled, exit the program
	fmt.Println("Context was cancelled.")

}
