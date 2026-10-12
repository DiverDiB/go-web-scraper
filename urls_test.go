package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	baseURL, _ := url.Parse("https://crawler-test.com")

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "absolute and relative urls",
			input:    `<html><body><a href="https://crawler-test.com/path1">Link 1</a><a href="/path2">Link 2</a></body></html>`,
			expected: []string{"https://crawler-test.com/path1", "https://crawler-test.com/path2"},
		},
		{
			name:     "multiple absolute urls",
			input:    `<html><body><a href="https://example.com/a">A</a><a href="https://example.com/b">B</a></body></html>`,
			expected: []string{"https://example.com/a", "https://example.com/b"},
		},
		{
			name:     "no links present",
			input:    `<html><body><p>No links here</p></body></html>`,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.input, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	baseURL, _ := url.Parse("https://crawler-test.com")

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "relative image src",
			input:    `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected: []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:     "multiple images with absolute and relative",
			input:    `<html><body><img src="/img1.jpg"><img src="https://other.com/img2.jpg"></body></html>`,
			expected: []string{"https://crawler-test.com/img1.jpg", "https://other.com/img2.jpg"},
		},
		{
			name:     "no images present",
			input:    `<html><body><p>No images</p></body></html>`,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getImagesFromHTML(tc.input, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
