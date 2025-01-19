package logger

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	// Set up the logger to output to a file or console
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string) {
	logger.Println("INFO: " + msg)
}

func Error(msg string) {
	logger.Println("ERROR: " + msg)
}

func Debug(msg string) {
	// Optionally log to a separate debug log
	logger.Println("DEBUG: " + msg)
}
