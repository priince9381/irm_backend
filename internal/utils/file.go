package utils

import (
	"os"
	"time"
)

// ParseTime parses a time string in RFC3339 format
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

// DeleteFile deletes a file at the given path
func DeleteFile(path string) error {
	return os.Remove(path)
}

// EnsureDir ensures that a directory exists, creating it if necessary
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
