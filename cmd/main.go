package main

import (
	"os"

	"cmt/internal/app"
)

// main is the entry point for the application
func main() {
	exitCode := app.Run()

	os.Exit(exitCode)
}
