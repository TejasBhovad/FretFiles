package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	watch "fret-files/internal/watch"

	"github.com/joho/godotenv"
	"github.com/sevlyar/go-daemon"
)

func main() {
	cntxt := &daemon.Context{
		PidFileName: "daemon.pid",
		PidFilePerm: 0644,
		LogFileName: "daemon.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[daemon]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("Daemon started")

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

	log.Println("Daemon terminated")
}
