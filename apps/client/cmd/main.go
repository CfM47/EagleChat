package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const logFilePath = "/data/responses.log"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./client <ip> <port>")
		os.Exit(1)
	}

	ip := os.Args[1]
	port := os.Args[2]
	url := fmt.Sprintf("http://%s:%s/status", ip, port)

	fmt.Printf("Starting client...\nPolling %s every second.\nLogging responses to %s\n", url, logFilePath)

	// Ensure log file directory exists.
	if err := os.MkdirAll(path.Dir(logFilePath), 0755); err != nil {
		log.Fatalf("Error creating log directory: %v", err)
	}

	// Open or create the log file.
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		timestamp := time.Now().Format(time.RFC3339)
		var logMessage string

		resp, err := http.Get(url)
		if err != nil {
			logMessage = fmt.Sprintf("%s - ERROR - Message: %s\n", timestamp, err.Error())
		} else {
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				logMessage = fmt.Sprintf("%s - ERROR - Failed to read response body: %s\n", timestamp, readErr.Error())
			} else {
				logMessage = fmt.Sprintf("%s - SUCCESS - Status: %s - Data: %s\n", timestamp, resp.Status, string(body))
			}
			resp.Body.Close()
		}

		if _, err := logFile.WriteString(logMessage); err != nil {
			log.Printf("Failed to write to log file: %v", err) // Log to stderr if file write fails
		}
	}
}
