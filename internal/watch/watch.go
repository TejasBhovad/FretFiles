package watch

import (
	"encoding/json"
	organise "fret-files/internal/organise"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type FileInfo struct {
	Path         string    `json:"path"`
	Modification time.Time `json:"modification"`
}

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
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	jsonFilePath := filepath.Join(watchPath, "known_files.json")
	knownFiles, err := loadKnownFiles(jsonFilePath)
	if err != nil {
		log.Println("Could not load known files:", err)
		knownFiles = make(map[string]time.Time)
	}

	var wg sync.WaitGroup
	updateKnownFiles(watchPath, knownFiles)

	// Save known files immediately after updating
	if err := saveKnownFiles(jsonFilePath, knownFiles); err != nil {
		log.Println("Error saving known files:", err)
	} else {
		log.Println("Known files saved successfully.")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			changed := checkForNewFiles(watchPath, knownFiles, &wg)
			if changed {
				if err := saveKnownFiles(jsonFilePath, knownFiles); err != nil {
					log.Println("Error saving known files:", err)
				} else {
					log.Println("Known files updated and saved successfully.")
				}
			}
		case <-done:
			log.Println("Stopping file watcher")
			wg.Wait()
			if err := saveKnownFiles(jsonFilePath, knownFiles); err != nil {
				log.Println("Error saving known files:", err)
			} else {
				log.Println("Known files saved successfully.")
			}
			return
		}
	}
}

func loadKnownFiles(filePath string) (map[string]time.Time, error) {
	knownFiles := make(map[string]time.Time)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Known files JSON does not exist; starting fresh.")
			return knownFiles, nil // Return empty map if file doesn't exist
		}
		return knownFiles, err
	}

	var files []FileInfo
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, err
	}

	for _, file := range files {
		knownFiles[file.Path] = file.Modification
	}

	return knownFiles, nil
}

func saveKnownFiles(filePath string, knownFiles map[string]time.Time) error {
	var files []FileInfo
	for path, modTime := range knownFiles {
		files = append(files, FileInfo{Path: path, Modification: modTime})
	}

	data, err := json.Marshal(files)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Creating directory: %s\n", dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	log.Printf("Saving known files to: %s\n", filePath)
	return os.WriteFile(filePath, data, 0644)
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
			info, err := os.Stat(fullPath)
			if err != nil {
				log.Println("Error getting file info:", err)
				continue
			}
			knownFiles[fullPath] = info.ModTime()
		}
	}
}

func checkForNewFiles(watchPath string, knownFiles map[string]time.Time, wg *sync.WaitGroup) bool {
	files, err := os.ReadDir(watchPath)
	if err != nil {
		log.Println("Error reading directory:", err)
		return false
	}

	changed := false

	for _, file := range files {
		if !file.IsDir() && !isIgnored(file.Name()) {
			fullPath := filepath.Join(watchPath, file.Name())
			info, err := os.Stat(fullPath)
			if err != nil {
				log.Println("Error getting file info:", err)
				continue
			}

			if knownModTime, exists := knownFiles[fullPath]; !exists || info.ModTime().After(knownModTime) {
				log.Println("New or modified file detected:", file.Name())
				wg.Add(1)
				go func() {
					defer wg.Done()
					handleNewFile(fullPath)
				}()
				knownFiles[fullPath] = info.ModTime()
				changed = true
			}
		}
	}

	return changed
}

func isIgnored(fileName string) bool {
	for _, ignored := range ignoredFiles {
		if strings.HasPrefix(ignored, "._") && strings.HasPrefix(fileName, "._") {
			return true
		}
		if fileName == ignored {
			return true
		}
	}
	return false
}

func handleNewFile(filePath string) {
	organise.DetermineFolder(filePath)
}
