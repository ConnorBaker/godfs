package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

func ReadBytes(path string, size int64) ([]byte, error) {
	// Open the file for reading
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	} else {
		// Close the file when we're done
		defer file.Close()
	}

	// Read the file's contents into a buffer
	log.Println("Allocating data buffer of size:", size)
	data := make([]byte, size)
	log.Println("Reading file into buffer")
	if _, err := io.ReadFull(file, data); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}
