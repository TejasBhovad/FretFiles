package main

import (
	watch "fret-files/internal"

	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	WatchPath := os.Getenv("WATCH_PATH")
	fmt.Println("Hello, World!")
	watch.Add(WatchPath)

}
