package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var logger log.Logger
var once sync.Once

func GetLogger() *log.Logger {
	once.Do(func() {
		logFile, err := os.OpenFile("log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("nope...", err)
			return
		}

		writer := io.MultiWriter(logFile, os.Stdout)
		logger.SetOutput(writer)
	})
	return &logger
}
