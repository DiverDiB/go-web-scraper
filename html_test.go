package main

import (
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "h1 present",
			input:    "<html><body><h1>Hello World</h1></body></html>",
			expected: "Hello World",
		},
		{
			name:     "h1 missing, h2 fallback",
			input:    "<html><body><h2>Fallback Heading</h2></body></html>",
			expected: "Fallback Heading",
		},
		{
			name:     "both missing",
			input:    "<html><body><p>Just text</p></body></html>",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "main tag priority",
			input: `<html><body>
				<p>Outside paragraph.</p>
				<main><p>Main paragraph.</p></main>
			</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "fallback to first paragraph when main missing",
			input: `<html><body>
				<p>First general paragraph.</p>
				<p>Second general paragraph.</p>
			</body></html>`,
			expected: "First general paragraph.",
		},
		{
			name:     "no paragraph present",
			input:    "<html><body><h1>Only a title</h1></body></html>",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}
