package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme https",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove scheme http",
			inputURL: "http://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "remove trailing slash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "handle capitals",
			inputURL: "https://WWW.BOOT.DEV/blog/path/",
			expected: "www.boot.dev/blog/path",
		},
		{
			name:     "http with capital domain and trailing slash",
			inputURL: "http://EXAMPLE.com/Page/One/",
			expected: "example.com/page/one",
		},
		{
			name:     "root domain with trailing slash",
			inputURL: "https://boot.dev/",
			expected: "boot.dev",
		},
		{
			name:     "root domain without trailing slash",
			inputURL: "https://boot.dev",
			expected: "boot.dev",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
