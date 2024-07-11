package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define command line flags
	directory := flag.String("d", ".", "The directory to start searching from")
	pattern := flag.String("p", "", "The pattern to search for in file/directory names")

	flag.Parse()

	if *pattern == "" {
		fmt.Println("Pattern is required")
		flag.Usage()
		os.Exit(1)
	}

	// Convert pattern to lowercase for case-insensitive matching
	patternLower := strings.ToLower(*pattern)

	// Walk the directory
	err := filepath.Walk(*directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Skip directories that cannot be opened
			if os.IsPermission(err) {
				log.Printf("Skipping %s: %v\n", path, err)
				return filepath.SkipDir
			}
			return err
		}
		if strings.Contains(strings.ToLower(info.Name()), patternLower) {
			fmt.Println(path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %v: %v\n", *directory, err)
	}
}
