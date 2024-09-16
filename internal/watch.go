package watch

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ignoredFiles = []string{
	".DS_Store",
	".Trashes",
	".Spotlight-V100",
	".fseventsd",
	".AppleDouble",
	"._*",
	".DocumentRevisions-V100",
	".PKG",
}

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
		if !file.IsDir() && !isIgnored(file.Name()) {
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
		if !file.IsDir() && !isIgnored(file.Name()) {
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
func isIgnored(fileName string) bool {
	for _, ignored := range ignoredFiles {
		// Check for exact match or prefix match (for patterns like "._*")
		if strings.HasPrefix(ignored, "._") && strings.HasPrefix(fileName, "._") {
			return true
		}
		if fileName == ignored {
			return true
		}
	}
	return false
}
