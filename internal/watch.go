package watch

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

func WatchFolder(watchPath string, done chan bool) {
	knownFiles := make(map[string]time.Time)

	// first time store al files
	updateKnownFiles(watchPath, knownFiles)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkForNewFiles(watchPath, knownFiles)
		case <-done:
			log.Println("Stopping file watcher")
			return
		}
	}
}

func updateKnownFiles(watchPath string, knownFiles map[string]time.Time) {
	files, err := os.ReadDir(watchPath)
	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(watchPath, file.Name())
			info, err := file.Info()
			if err != nil {
				log.Println("Error getting file info:", err)
				continue
			}
			knownFiles[fullPath] = info.ModTime()
		}
	}
}

func checkForNewFiles(watchPath string, knownFiles map[string]time.Time) {
	files, err := os.ReadDir(watchPath)
	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(watchPath, file.Name())
			info, err := file.Info()
			if err != nil {
				log.Println("Error getting file info:", err)
				continue
			}

			knownModTime, exists := knownFiles[fullPath]
			if !exists || info.ModTime().After(knownModTime) {
				log.Println("New or modified file detected:", file.Name())
				knownFiles[fullPath] = info.ModTime()
			}
		}
	}
}
