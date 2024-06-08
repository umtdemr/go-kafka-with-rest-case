package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/logger"
	"github.com/umtdemr/go-kafka-with-rest-case/pkg/server"
	"log"
)

func main() {
	logger.GetLogger()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("An error has occurred while creating a file watcher: %s\n", err)
	}
	defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("modified file", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error", err)
			}
		}
	}()
	err = watcher.Add("log.log")
	if err != nil {
		log.Fatal(err)
	}
	server.Run()
}
