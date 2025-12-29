package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ahrdadan/image-metadata-viewer/src/internal/models"
	"github.com/ahrdadan/image-metadata-viewer/src/internal/utils"
	"github.com/ahrdadan/image-metadata-viewer/src/pkg/metadata"
)

const (
	// MaxImageBytes is the maximum size for image downloads
	MaxImageBytes = 20 << 20 // 20MB
	// MaxUploadBytes is the maximum size for uploads
	MaxUploadBytes = MaxImageBytes + (1 << 20)
)

// ImageService handles image processing operations
type ImageService struct {
	httpClient *http.Client
}

// NewImageService creates a new ImageService
func NewImageService() *ImageService {
	return &ImageService{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ProcessUpload processes an uploaded image file
func (s *ImageService) ProcessUpload(data []byte, contentType, fileName string) *models.ImageMetadata {
	meta := metadata.ExtractMetadata(data, contentType, fileName)
	meta.Source = "upload"
	return meta
}

// ProcessRemoteURL downloads and processes an image from a URL
func (s *ImageService) ProcessRemoteURL(ctx context.Context, imageURL string) *models.ImageMetadata {
	meta := &models.ImageMetadata{
		Source: "remote",
	}

	// Parse and validate URL
	parsed, err := url.Parse(imageURL)
	if err != nil {
		meta.FetchError = fmt.Sprintf("invalid URL: %v", err)
		return meta
	}

	meta.FinalURL = parsed.String()
	meta.FileName = utils.FileNameFromURL(parsed)
	meta.FileTypeExtension = utils.ExtensionFromName(meta.FileName)

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		meta.FetchError = fmt.Sprintf("request error: %v", err)
		return meta
	}
	req.Header.Set("User-Agent", "image-metadata-viewer/2.0")

	// Execute request
	start := time.Now()
	resp, err := s.httpClient.Do(req)
	meta.Duration = time.Since(start).Round(time.Millisecond).String()
	if err != nil {
		meta.FetchError = fmt.Sprintf("fetch error: %v", err)
		return meta
	}
	defer resp.Body.Close()

	// Check status
	if resp.StatusCode != http.StatusOK {
		meta.FetchError = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
		meta.Status = resp.Status
		return meta
	}

	meta.Status = resp.Status
	meta.MIMEType = resp.Header.Get("Content-Type")
	meta.ContentLength = resp.ContentLength
	meta.LastModified = resp.Header.Get("Last-Modified")

	// Download to temp file then read for metadata extraction
	tempPath, downloadedBytes, truncated, err := downloadToTempFile(resp.Body, MaxImageBytes)
	if err != nil {
		meta.FetchError = fmt.Sprintf("read error: %v", err)
		return meta
	}
	defer os.Remove(tempPath)

	body, err := os.ReadFile(tempPath)
	if err != nil {
		meta.FetchError = fmt.Sprintf("read error: %v", err)
		return meta
	}

	meta.DownloadedBytes = downloadedBytes
	meta.Truncated = truncated

	if len(body) == 0 {
		meta.FetchError = "empty response"
		return meta
	}

	// Extract metadata
	extracted := metadata.ExtractMetadata(body, meta.MIMEType, meta.FileName)

	// Merge data
	extracted.Source = meta.Source
	extracted.Status = meta.Status
	extracted.FinalURL = meta.FinalURL
	extracted.ContentLength = meta.ContentLength
	extracted.DownloadedBytes = meta.DownloadedBytes
	extracted.Truncated = meta.Truncated
	extracted.Duration = meta.Duration
	extracted.LastModified = meta.LastModified

	return extracted
}

// ProcessMultipleURLs processes multiple URLs concurrently
func (s *ImageService) ProcessMultipleURLs(ctx context.Context, urls []string) []*models.ImageMetadata {
	results := make([]*models.ImageMetadata, len(urls))

	for i, imageURL := range urls {
		results[i] = s.ProcessRemoteURL(ctx, imageURL)
	}

	return results
}

// downloadToTempFile streams data to a temp file with size limits and returns its path.
func downloadToTempFile(r io.Reader, maxBytes int64) (string, int64, bool, error) {
	if maxBytes <= 0 {
		return "", 0, false, fmt.Errorf("invalid maxBytes")
	}

	tmp, err := os.CreateTemp("", "image-metadata-*")
	if err != nil {
		return "", 0, false, err
	}
	defer tmp.Close()

	limited := io.LimitReader(r, maxBytes+1)
	written, err := io.Copy(tmp, limited)
	if err != nil {
		return tmp.Name(), written, false, err
	}

	truncated := false
	if written > maxBytes {
		truncated = true
		if err := tmp.Truncate(maxBytes); err != nil {
			return tmp.Name(), written, truncated, err
		}
		written = maxBytes
	}

	return tmp.Name(), written, truncated, nil
}
