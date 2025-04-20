package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

func runMakeTarget(target string) error {
	cmd := exec.Command("make", target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	godotenv.Load()
	fmt.Println("running build...")

	mediaManifest := os.Getenv("MEDIA_MANIFEST_PATH")

	err := os.RemoveAll("./public")
	if err != nil {
		fmt.Println("public directory does not exist, creating now ...")
		if err := os.MkdirAll("./public", 0755); err != nil {
			log.Panicf("Could not create directory %s: %v", "./public", err)
		}
	}

	dirs := []string{
		"public",
		filepath.Join("public", "scripts"),
		filepath.Join("public", "scripts", "bun"),
		filepath.Join("public", "scripts", "htmx"),
		filepath.Join("public", "styles"),
	}

	// Create each directory
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Panicf("Could not create directory %s: %v", dir, err)
		}
	}

	// remove previous manifest
	// os.Stat(mediaManifest)
	// if err != nil {
	// 	os.Remove(mediaManifest)
	// }

	file, err := os.Create(mediaManifest)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("{}")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	// TODO: write minimal setup for manifest file
	// TODO: create optional flag for compiling after
	fmt.Println("Compiling bundler...")
	if err := runMakeTarget("bundler"); err != nil {
		fmt.Println("❌ bundler build failed:", err)
		return
	}

	fmt.Println("✅ Bundler complete")

	fmt.Println("Build: CSS")
	if err := runMakeTarget("css"); err != nil {
		fmt.Println("❌ CSS build failed:", err)
		return
	}

	fmt.Println("✅ CSS build complete. ")
	fmt.Println("Build: JS scripts")

	if err := runMakeTarget("js"); err != nil {
		fmt.Println("❌ JS build failed:", err)
		return
	}

	fmt.Println("✅ JS build complete.")

	fmt.Println("Build: HTMX scripts")

	if err := runMakeTarget("htmx"); err != nil {
		fmt.Println("❌ JS build failed:", err)
		return
	}

	fmt.Println("✅ HTMX build complete.")
}
