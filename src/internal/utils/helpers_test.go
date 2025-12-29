package utils

import (
	"testing"
)

func TestHumanBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"negative", -1, "unknown"},
		{"zero", 0, "0 B"},
		{"bytes", 500, "500 B"},
		{"exact_kb", 1024, "1.0 KiB"},
		{"kilobytes", 1536, "1.5 KiB"},
		{"megabytes", 1048576, "1.0 MiB"},
		{"large", 5242880, "5.0 MiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HumanBytes(tt.input)
			if result != tt.expected {
				t.Errorf("HumanBytes(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"already_normal", "https://example.com/image.jpg", "https://example.com/image.jpg"},
		{"single_slash_http", "http:/example.com/image.jpg", "http://example.com/image.jpg"},
		{"single_slash_https", "https:/example.com/image.jpg", "https://example.com/image.jpg"},
		{"with_spaces", "  https://example.com/image.jpg  ", "https://example.com/image.jpg"},
		{"no_change_needed", "http://example.com", "http://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeURL(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtensionFromName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"jpg", "image.jpg", "jpg"},
		{"png", "photo.png", "png"},
		{"uppercase", "PHOTO.JPG", "jpg"},
		{"no_extension", "filename", ""},
		{"multiple_dots", "file.name.jpeg", "jpeg"},
		{"hidden_file", ".gitignore", "gitignore"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtensionFromName(tt.input)
			if result != tt.expected {
				t.Errorf("ExtensionFromName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalculateAspectRatio(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected string
	}{
		{"landscape", 1920, 1080, "1.778"},
		{"portrait", 1080, 1920, "0.562"},
		{"square", 1000, 1000, "1.000"},
		{"zero_height", 1920, 0, "N/A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAspectRatio(tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("CalculateAspectRatio(%d, %d) = %s, want %s", tt.width, tt.height, result, tt.expected)
			}
		})
	}
}

func TestCalculateAspectRatioFraction(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected string
	}{
		{"landscape", 1920, 1080, "16:9"},
		{"portrait", 1080, 1920, "9:16"},
		{"square", 1000, 1000, "1:1"},
		{"classic", 1024, 768, "4:3"},
		{"zero_height", 1920, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAspectRatioFraction(tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("CalculateAspectRatioFraction(%d, %d) = %s, want %s", tt.width, tt.height, result, tt.expected)
			}
		})
	}
}

func TestCalculateMegapixels(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected float64
	}{
		{"hd", 1920, 1080, 2.0736},
		{"4k", 3840, 2160, 8.2944},
		{"small", 800, 600, 0.48},
		{"zero", 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMegapixels(tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("CalculateMegapixels(%d, %d) = %f, want %f", tt.width, tt.height, result, tt.expected)
			}
		})
	}
}

func TestFormatToExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"jpeg", "jpeg", "jpg"},
		{"JPEG", "JPEG", "jpg"},
		{"png", "png", "png"},
		{"tiff", "tiff", "tif"},
		{"gif", "gif", "gif"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatToExtension(tt.input)
			if result != tt.expected {
				t.Errorf("FormatToExtension(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContentTypeBase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple", "image/jpeg", "image/jpeg"},
		{"with_charset", "image/jpeg; charset=utf-8", "image/jpeg"},
		{"empty", "", ""},
		{"with_spaces", "image/png ; charset=utf-8", "image/png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContentTypeBase(tt.input)
			if result != tt.expected {
				t.Errorf("ContentTypeBase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
