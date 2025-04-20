package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Created this while experimenting with css/js build processes.
// Decided to use bun for handling everything, but keeping this around for reference.
func main() {
	// Get the working directory where the command is being run
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}

	var file = flag.String("input-file", "", "the file to be hashed")

	flag.Parse()

	fmt.Println("file given", *file)
	// Define relative path to the CSS file (adjust based on your project structure)
	// cssRelativePath := "public/styles/main.css" // Adjust if needed

	absolutePath, err := filepath.Abs(filepath.Join(wd, *file))
	if err != nil {
		fmt.Println("Error resolving CSS path:", err)
		os.Exit(1)
	}

	// Check if the file exists before proceeding
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		fmt.Println("Error: CSS file not found at", absolutePath)
		os.Exit(1)
	}

	// Compute the hash
	hash, err := computeFileHash(absolutePath)
	if err != nil {
		fmt.Println("Error hashing CSS file:", err)
		os.Exit(1)
	}

	// Create new hashed filename (e.g., main.abc12345.css)
	newFileName := fmt.Sprintf("main.%s.css", hash)

	// Compute the absolute new file path
	cssDir := filepath.Dir(absolutePath)
	newFilePath := filepath.Join(cssDir, newFileName)

	// Rename the file
	if err := os.Rename(absolutePath, newFilePath); err != nil {
		fmt.Println("Error renaming file:", err)
		os.Exit(1)
	}

	// Save the hashed filename in a JSON manifest
	manifestPath := filepath.Join(cssDir, "manifest.json")
	manifestData := map[string]string{"css": newFileName}
	if err := writeManifest(manifestPath, manifestData); err != nil {
		fmt.Println("Error writing manifest:", err)
		os.Exit(1)
	}

	fmt.Printf("✅ CSS Cache Busting: Renamed to %s\n", newFileName)
}

// computeFileHash generates an 8-character MD5 hash of a file's contents.
func computeFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil))[:8], nil
}

// writeManifest saves the hashed filename to a JSON file (optional).
func writeManifest(path string, data map[string]string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, 0644)
}
