package main

import (
	"fmt"
	"os/exec"
  "log"
	"os"
	"path/filepath"
	"regexp"
  "strings"
)
func cleanAndCompressFiles(){
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Read all files in the current directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	// Define a regex to match and keep only alphanumeric characters and underscores
	re := regexp.MustCompile("[^a-zA-Z0-9_-]")

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}
    if !strings.HasSuffix(file.Name(), ".mp3"){
      continue
    }
		// Get the original file name
		oldName := file.Name()

		// Split the name and extension
		ext := filepath.Ext(oldName)
		name := oldName[:len(oldName)-len(ext)]

		// Create the new name by replacing special characters in the name part
		newName := re.ReplaceAllString(name, "_") + ext
		
		// If the names are different, rename the file
		if oldName != newName {
			oldPath := filepath.Join(dir, oldName)
			newPath := filepath.Join(dir, newName)

			fmt.Printf("Renaming '%s' to '%s'\n", oldName, newName)
			
			// Rename the file
			err := os.Rename(oldPath, newPath)
			if err != nil {
				log.Printf("Failed to rename '%s' to '%s': %v", oldName, newName, err)
			}
		}
	}
  newfiles, err := os.ReadDir(dir)
	  for _, file := range newfiles {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".mp3") {
			continue
		}
		err = compressAudio(file.Name(), "temp.mp3")
		if err != nil {
			fmt.Println(err)
		}
		err = os.Rename(filepath.Join(dir, "temp.mp3"), filepath.Join(dir, file.Name()))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func compressAudio(inputFile, outputFile string) error {
	fmt.Println("input " + inputFile)
	fmt.Println("output " + outputFile)
	// Define the ffmpeg command and arguments
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-b:a", "128k", outputFile)

	// Run the command and capture any errors
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run ffmpeg: %v", err)
	}
	err = os.Remove(inputFile)
	if err != nil {
		return fmt.Errorf("failed to delete: %v", err)
	}
	return nil
}
