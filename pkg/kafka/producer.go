package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("An error has occurred while creating a file watcher: %s\n", err)
	}
	defer watcher.Close()

	producer, err := sarama.NewSyncProducer([]string{brokers}, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	defer producer.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					lastLog := readLastLog(event.Name)
					message := &sarama.ProducerMessage{
						Topic: topics,
						Value: sarama.StringEncoder(lastLog),
					}
					partition, offset, err := producer.SendMessage(message)
					if err != nil {
						log.Printf("Failed to send message: %v", err)
					} else {
						log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
					}
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

	<-stop
	log.Println("Producer stopped.")
}

func readLastLog(fileName string) string {
	fileHandle, fileOpenErr := os.Open(fileName)
	if fileOpenErr != nil {
		panic("error while opening the file...")
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size() - 5
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			break
		}

		if char[0] != 10 {
			line = fmt.Sprintf("%s%s", string(char), line)
		}

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}
	return line
}
