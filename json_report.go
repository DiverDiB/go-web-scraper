package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	// sort the map keys
	keys := make([]string, 0, len(pages))
	for key := range pages {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build a[]PageData slice in sorted order
	sorted_pages := make([]PageData, 0, len(pages))
	for _, key := range keys {
		sorted_pages = append(sorted_pages, pages[key])
	}

	// marshal
	data, err := json.MarshalIndent(sorted_pages, "", "  ")
	if err != nil {
		fmt.Printf("error upon Marshal: %v\n", err)
		os.Exit(1)
	}

	// write
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("error writing file: %v\n", err)
		os.Exit(1)
	}
	return nil
}
