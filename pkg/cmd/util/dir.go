package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// EnsureDir ensures that the given directory exists.
// If an error occurs, the program fails.
func EnsureDir(dirName string) {
	if err := os.MkdirAll(dirName, 0755); err != nil {
		log.WithError(err).Fatalf("Cannot ensure directory %s exists", dirName)
	}
}

// WorkingDir returns the current working directory.
// If an error occurs, the program fails.
func WorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal("Cannot get current working directory")
	}
	return dir
}
