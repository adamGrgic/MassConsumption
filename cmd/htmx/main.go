package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type Manifest map[string]string

// Downloads a file and returns its raw bytes
func fetchBytes(rawURL string) ([]byte, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for %s: %s", rawURL, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %w", rawURL, err)
	}

	return body, nil
}

// Hashes the content and returns an 8-char hash
func getHash(content []byte) string {
	hash := md5.Sum(content)
	return hex.EncodeToString(hash[:])[:8]
}

// Extracts the base name from the URL path
func baseNameFromURL(rawURL string) string {
	u, _ := url.Parse(rawURL)
	return filepath.Base(u.Path)
}

// Downloads, hashes, saves the file, and returns the logical name + final filename
func downloadAndHash(scriptURL, buildDir string) (string, string, error) {
	content, err := fetchBytes(scriptURL)
	if err != nil {
		return "", "", err
	}

	hash := getHash(content)
	base := baseNameFromURL(scriptURL)
	logical := strings.TrimSuffix(base, filepath.Ext(base))
	hashedName := fmt.Sprintf("%s.%s%s", logical, hash, filepath.Ext(base))

	// Save file
	outputPath := filepath.Join(buildDir, hashedName)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		return "", "", err
	}

	fmt.Printf("✅ Saved %s → %s\n", scriptURL, outputPath)
	return logical, filepath.ToSlash(filepath.Join("public/scripts/htmx", hashedName)), nil
}

func main() {
	fmt.Println("📦 Downloading and hashing HTMX scripts...")

	godotenv.Load()

	scripts := []string{
		os.Getenv("HTMX"),
		os.Getenv("HTMX_TARGET"),
	}

	buildDir := os.Getenv("HTMX_BUILD_DIR")
	if buildDir == "" {
		log.Fatal("HTMX_BUILD_DIR is not set")
	}

	manifestPath := os.Getenv("MEDIA_MANIFEST_PATH")
	if manifestPath == "" {
		log.Fatal("Manifest not set while setting HTMX scripts")
	}

	manifest := Manifest{}
	if data, err := os.ReadFile(manifestPath); err == nil {
		if err := json.Unmarshal(data, &manifest); err != nil {
			log.Fatalf("❌ Failed to parse existing manifest: %v", err)
		}
	}

	for _, scriptURL := range scripts {
		if scriptURL == "" {
			continue
		}

		logical, hashedPath, err := downloadAndHash(scriptURL, buildDir)
		if err != nil {
			log.Panic("Something went wrong while downloading HTMX scripts: ", err)
		}
		manifest[logical] = hashedPath
	}

	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal manifest:", err)
	}
	if err := os.WriteFile(manifestPath, manifestBytes, 0644); err != nil {
		log.Fatal("Failed to write manifest:", err)
	}

	fmt.Printf("\n📝 Manifest updated: %s\n", manifestPath)
}
