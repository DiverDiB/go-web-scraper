package main

import (
	//"fmt"
	"net/url"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {
	//fmt.Printf("inputURL = %s\n", inputURL)
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}
	//fmt.Printf("parsedURL = %s\n", parsedURL)

	// Combine host and path, then clean up any trailing slashes
	fullPath := parsedURL.Host + parsedURL.Path
	//fmt.Printf("fullPath1 = %s\n", fullPath)
	fullPath = strings.TrimSuffix(fullPath, "/")
	//fmt.Printf("fullPath2 = %s\n", fullPath)
	//fmt.Printf("lower case = %s\n", strings.ToLower(fullPath))

	// Convert to lowercase to ensure consistency
	return strings.ToLower(fullPath), nil
}
