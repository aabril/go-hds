package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	// Define command line flags
	directory := flag.String("d", ".", "The directory to start searching from")
	pattern := flag.String("p", "", "The pattern to search for in file/directory names")
	workers := flag.Int("w", 8, "Number of worker goroutines")

	flag.Parse()

	if *pattern == "" {
		fmt.Println("Pattern is required")
		flag.Usage()
		os.Exit(1)
	}

	// Convert pattern to lowercase for case-insensitive matching
	patternLower := strings.ToLower(*pattern)

	// Create channels for communication
	filePaths := make(chan string)
	results := make(chan string)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(&wg, filePaths, results, patternLower)
	}

	// Start a goroutine to collect results
	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	// Walk the directory and send file paths to the worker pool
	go func() {
		err := filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Skip directories that cannot be opened
				if os.IsPermission(err) {
					// log.Printf("Skipping %s: %v\n", path, err)
					return filepath.SkipDir
				}
				return err
			}
			if !info.IsDir() {
				filePaths <- path
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking the path %v: %v\n", *directory, err)
		}
		close(filePaths)
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(results)
}

// Worker function to process file paths
func worker(wg *sync.WaitGroup, filePaths <-chan string, results chan<- string, pattern string) {
	defer wg.Done()
	for path := range filePaths {
		if strings.Contains(strings.ToLower(filepath.Base(path)), pattern) {
			results <- path
		}
	}
}
