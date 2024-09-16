package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	watch "fret-files/internal"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	watchPath := os.Getenv("WATCH_PATH")
	if watchPath == "" {
		log.Fatal("WATCH_PATH is not set in .env file")
	}

	done := make(chan bool) // done = stop
	go watch.WatchFolder(watchPath, done)

	// try ctrl C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// basically catch
	select {
	case <-quit:
		log.Println("\nShutting down gracefully...")
		done <- true
	case <-done:
	}
}
