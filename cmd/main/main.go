package main

import (
	"context"
	"github.com/spf13/viper"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/kafka"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/logger"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/server"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/store"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// read env
	viper.SetConfigFile(".env")
	errReadConfig := viper.ReadInConfig()

	if errReadConfig != nil {
		log.Fatal("could not read the env file")
	}

	dbStore, pgConnErr := store.NewStore()

	if pgConnErr != nil {
		log.Fatal("could not connect to db")
	}

	defer dbStore.Close()

	// init tables
	dbInitErr := dbStore.Init()

	if dbInitErr != nil {
		log.Fatal("error while initiating the tables")
	}

	// initialize the file and stdout logger
	logger.GetLogger()

	// block app with the stop chan
	stop := make(chan os.Signal, 1)
	consumerStarted := make(chan bool, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go kafka.StartConsumer(consumerStarted, dbStore)
	// wait kafka consumer
	<-consumerStarted

	// start kafka producer
	go kafka.StartProducer()

	log.Println("Starting HTTP server...")

	server.Run(dbStore)
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
