package main

import (
	"fmt"
	"os"
	"os/exec"
)

var issues = false

func checkCommand(module string, level string) {
	_, err := exec.LookPath(module)
	if err != nil {
		fmt.Printf("❌ [%s] %s is missing!\n", level, module)
		issues = true
	}
}

func main() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		fmt.Println("❌ [warning] APP_ENV is not set!")
		issues = true
	}
	// checkCommand(appEnv, "warning")
	checkCommand("docker", "warning")
	checkCommand("bun", "fatal")

	if !issues {
		fmt.Printf("✅ All checks are passing\n")
	}
}
