package main

import (
	"os"
	"path"
	"testing"
)

func TestMain(t *testing.T) {
	dir := "build"

	main()

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("Build directory does not exist: %v", err)
	}
	if files, err := os.ReadDir(dir); err != nil || len(files) == 0 {
		t.Fatalf("Build directory is empty or cannot be read: %v", err)
	}

	expectedFiles := []string{
		"index.html",
		"404.html",
		"static/icon.svg",
		"static/styles/style.css",
		"static/styles/reset.css",
	}

	for _, file := range expectedFiles {
		if _, err := os.Stat(path.Join(dir, file)); os.IsNotExist(err) {
			t.Fatalf("Expected file does not exist: %s", file)
		}
	}
}
