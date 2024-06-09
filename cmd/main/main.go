package main

import (
	"context"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/kafka"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/logger"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.GetLogger()
	stop := make(chan os.Signal, 1)
	consumerStarted := make(chan bool, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go kafka.StartConsumer(consumerStarted)
	<-consumerStarted
	go kafka.StartProducer()

	log.Println("Starting HTTP server...")

	server.Run()
	// Wait for termination signal
	<-stop
	log.Println("Shutting down gracefully...")

	// Create a context with timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serverStop := make(chan struct{})
	// Shut down the server
	server.Shutdown(ctx, serverStop)

	// Wait for the server to complete shutdown
	<-serverStop
	log.Println("HTTP server stopped.")

	log.Println("All services stopped.")
}
