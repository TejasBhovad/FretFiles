package organise

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var folderExtensions = map[string][]string{
	"images": {".png", ".jpeg", ".webp", ".jpg", ".gif", ".tiff", ".raw"},
	"pdfs":   {".pdf"},
}

func DetermineFolder(filePath string) {
	loadEnvVars()

	fmt.Println("Path from organise:", filePath)

	_, fileName := splitFromLast(filePath, "/")
	extension := strings.ToLower(getFileExtension(fileName))

	folderName := getFolderName(extension)
	if folderName == "" {
		fmt.Println("This file type is not recognized.")
		return
	}

	moveFileToFolder(filePath, folderName)
}

func splitFromLast(s, sep string) (string, string) {
	lastIndex := strings.LastIndex(s, sep)

	if lastIndex == -1 {
		return s, ""
	}

	part1 := s[:lastIndex]
	part2 := s[lastIndex+len(sep):]

	return part1, part2
}

func getFileExtension(fileName string) string {
	lastDotIndex := strings.LastIndex(fileName, ".")

	if lastDotIndex == -1 {
		return ""
	}

	return fileName[lastDotIndex:]
}

func getFolderName(ext string) string {
	for folderName, extensions := range folderExtensions {
		for _, e := range extensions {
			if ext == e {
				return folderName
			}
		}
	}

	return ""
}

func moveFileToFolder(filePath, folderName string) {
	watchPath := os.Getenv("WATCH_PATH")
	targetDir := filepath.Join(watchPath, folderName)

	// Create the target directory if it doesn't exist
	createDirIfNotExists(targetDir)

	targetFilePath := filepath.Join(targetDir, filepath.Base(filePath))

	// Move the file to the new directory
	if err := os.Rename(filePath, targetFilePath); err != nil {
		log.Fatalf("Error moving file from %s to %s: %v", filePath, targetFilePath, err)
	}

	fmt.Printf("Moved %s to %s\n", filePath, targetFilePath)
}

func loadEnvVars() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	watchPath := os.Getenv("WATCH_PATH")
	if watchPath == "" {
		log.Fatal("WATCH_PATH is not set in .env file")
	}
}

func createDirIfNotExists(dir string) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory %s: %v", dir, err)
	}
}
