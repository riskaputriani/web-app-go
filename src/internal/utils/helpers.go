package utils

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// HumanBytes converts bytes to human-readable format
func HumanBytes(size int64) string {
	if size < 0 {
		return "unknown"
	}
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit && exp < 3; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGT"[exp])
}

// NormalizeURL fixes single-slash schemes produced by path-based input
func NormalizeURL(raw string) string {
	normalized := strings.TrimSpace(raw)
	if strings.HasPrefix(normalized, "http:/") && !strings.HasPrefix(normalized, "http://") {
		return "http://" + strings.TrimPrefix(normalized, "http:/")
	}
	if strings.HasPrefix(normalized, "https:/") && !strings.HasPrefix(normalized, "https://") {
		return "https://" + strings.TrimPrefix(normalized, "https:/")
	}
	return normalized
}

// ExtensionFromName extracts file extension from filename
func ExtensionFromName(name string) string {
	ext := strings.TrimPrefix(strings.ToLower(path.Ext(name)), ".")
	return ext
}

// FileNameFromURL extracts filename from URL
func FileNameFromURL(parsed *url.URL) string {
	if parsed == nil {
		return ""
	}
	base := path.Base(parsed.Path)
	if base == "." || base == "/" {
		return ""
	}
	return base
}

// ContentTypeBase extracts base content type without parameters
func ContentTypeBase(contentType string) string {
	if contentType == "" {
		return ""
	}
	parts := strings.Split(contentType, ";")
	return strings.TrimSpace(parts[0])
}

// FormatToExtension converts image format to file extension
func FormatToExtension(format string) string {
	switch strings.ToLower(format) {
	case "jpeg":
		return "jpg"
	case "tiff":
		return "tif"
	default:
		return strings.ToLower(format)
	}
}

// CalculateAspectRatio calculates aspect ratio from dimensions
func CalculateAspectRatio(width, height int) string {
	if height == 0 {
		return "N/A"
	}
	ratio := float64(width) / float64(height)
	return fmt.Sprintf("%.3f", ratio)
}

// CalculateAspectRatioFraction returns a simplified width:height ratio.
func CalculateAspectRatioFraction(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	divisor := gcd(width, height)
	return fmt.Sprintf("%d:%d", width/divisor, height/divisor)
}

// CalculateMegapixels calculates megapixels from dimensions
func CalculateMegapixels(width, height int) float64 {
	return float64(width*height) / 1000000.0
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	if a == 0 {
		return 1
	}
	return a
}
