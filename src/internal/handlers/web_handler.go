package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/ahrdadan/image-metadata-viewer/src/internal/models"
	"github.com/ahrdadan/image-metadata-viewer/src/internal/services"
	"github.com/ahrdadan/image-metadata-viewer/src/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// WebHandler handles web interface requests
type WebHandler struct {
	imageService *services.ImageService
	blobStore    *services.BlobStore
	maxBytesMB   int
	maxBytesRaw  int64
}

// NewWebHandler creates a new WebHandler
func NewWebHandler(imageService *services.ImageService, blobStore *services.BlobStore) *WebHandler {
	return &WebHandler{
		imageService: imageService,
		blobStore:    blobStore,
		maxBytesMB:   services.MaxImageBytes / (1 << 20),
		maxBytesRaw:  services.MaxImageBytes,
	}
}

// HandleHome renders the home page
func (h *WebHandler) HandleHome(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{
		"BaseURL": h.getBaseURL(c),
	})
}

// HandleDocs renders the API documentation page
func (h *WebHandler) HandleDocs(c *fiber.Ctx) error {
	return c.Render("docs", fiber.Map{
		"Title":   "API Docs",
		"BaseURL": h.getBaseURL(c),
	})
}

// HandleForm processes the URL form submission
func (h *WebHandler) HandleForm(c *fiber.Ctx) error {
	rawURLs := c.Query("url")
	if rawURLs == "" {
		return c.Redirect("/", http.StatusSeeOther)
	}

	// Check if multiple URLs (newline separated)
	urls := strings.Split(rawURLs, "\n")
	cleanedURLs := make([]string, 0)

	for _, u := range urls {
		normalized := utils.NormalizeURL(strings.TrimSpace(u))
		if normalized != "" {
			cleanedURLs = append(cleanedURLs, normalized)
		}
	}

	if len(cleanedURLs) == 0 {
		return c.Redirect("/", http.StatusSeeOther)
	}

	// If single URL, use simple redirect
	if len(cleanedURLs) == 1 {
		escaped := url.PathEscape(cleanedURLs[0])
		return c.Redirect("/"+escaped, http.StatusSeeOther)
	}

	// Multiple URLs - process batch
	return h.processBatchURLs(c, cleanedURLs)
}

// HandleView displays image view for a single URL
func (h *WebHandler) HandleView(c *fiber.Ctx) error {
	rawPath := c.Params("*")
	if rawPath == "" {
		return h.HandleHome(c)
	}

	rawURL, err := url.PathUnescape(rawPath)
	if err != nil {
		return c.Render("view", h.buildErrorView("Invalid URL", "Could not decode the URL path.", rawPath, c))
	}

	normalizedURL := utils.NormalizeURL(rawURL)
	parsed, err := url.Parse(normalizedURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return c.Render("view", h.buildErrorView("Invalid URL", "The path must be a full http or https URL.", normalizedURL, c))
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return c.Render("view", h.buildErrorView("Unsupported URL", "Only http and https URLs are supported.", normalizedURL, c))
	}

	// Process the URL
	meta := h.imageService.ProcessRemoteURL(c.Context(), parsed.String())

	imageResult := models.ImageResult{
		InputURL:   normalizedURL,
		DisplayURL: parsed.String(),
		EmbedURL:   parsed.String(),
		Metadata:   meta,
	}

	if meta.FetchError != "" {
		imageResult.Error = meta.FetchError
	}

	return c.Render("view", fiber.Map{
		"Title":       "Image Preview",
		"BaseURL":     h.getBaseURL(c),
		"MaxBytesMB":  h.maxBytesMB,
		"MaxBytesRaw": h.maxBytesRaw,
		"Images":      []models.ImageResult{imageResult},
		"IsUpload":    false,
		"IsBatch":     false,
	})
}

// HandleUpload processes file uploads
func (h *WebHandler) HandleUpload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Render("view", h.buildErrorView("Upload error", "Could not parse upload form.", "", c))
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Render("view", h.buildErrorView("Upload error", "No files selected.", "", c))
	}

	results := make([]models.ImageResult, 0, len(files))

	for _, fileHeader := range files {
		result := h.processUploadedFile(fileHeader)
		results = append(results, result)
	}

	return c.Render("view", fiber.Map{
		"Title":       "Image Preview",
		"BaseURL":     h.getBaseURL(c),
		"MaxBytesMB":  h.maxBytesMB,
		"MaxBytesRaw": h.maxBytesRaw,
		"Images":      results,
		"IsUpload":    true,
		"IsBatch":     len(results) > 1,
	})
}

// HandleBlob serves temporary uploaded images.
func (h *WebHandler) HandleBlob(c *fiber.Ctx) error {
	if h.blobStore == nil {
		return c.SendStatus(http.StatusNotFound)
	}

	blobID := c.Params("id")
	if blobID == "" {
		return c.SendStatus(http.StatusNotFound)
	}

	data, contentType, ok := h.blobStore.Get(blobID)
	if !ok {
		return c.SendStatus(http.StatusNotFound)
	}

	if contentType != "" {
		c.Set("Content-Type", contentType)
	}
	c.Set("Cache-Control", "private, max-age=3600")
	return c.Send(data)
}

// processUploadedFile processes a single uploaded file
func (h *WebHandler) processUploadedFile(fileHeader *multipart.FileHeader) models.ImageResult {
	result := models.ImageResult{
		InputURL: fileHeader.Filename,
	}

	if fileHeader.Size > services.MaxUploadBytes {
		result.Error = fmt.Sprintf("File exceeds %d MB limit.", h.maxBytesMB)
		return result
	}

	file, err := fileHeader.Open()
	if err != nil {
		result.Error = fmt.Sprintf("Failed to open file: %v", err)
		return result
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read file: %v", err)
		return result
	}

	if len(data) == 0 {
		result.Error = "Uploaded file is empty."
		return result
	}

	// Detect content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	// Process image
	meta := h.imageService.ProcessUpload(data, contentType, fileHeader.Filename)
	result.Metadata = meta

	if h.blobStore != nil {
		blobID := h.blobStore.Put(data, contentType)
		result.EmbedURL = "/blob/" + blobID
		result.IsBlob = true
	}

	return result
}

// processBatchURLs processes multiple URLs
func (h *WebHandler) processBatchURLs(c *fiber.Ctx, urls []string) error {
	results := make([]models.ImageResult, 0, len(urls))

	for _, imageURL := range urls {
		parsed, err := url.Parse(imageURL)
		if err != nil {
			results = append(results, models.ImageResult{
				InputURL: imageURL,
				Error:    "Invalid URL",
			})
			continue
		}

		meta := h.imageService.ProcessRemoteURL(c.Context(), parsed.String())

		result := models.ImageResult{
			InputURL:   imageURL,
			DisplayURL: parsed.String(),
			EmbedURL:   parsed.String(),
			Metadata:   meta,
		}

		if meta.FetchError != "" {
			result.Error = meta.FetchError
		}

		results = append(results, result)
	}

	return c.Render("view", fiber.Map{
		"Title":       "Batch Image Preview",
		"BaseURL":     h.getBaseURL(c),
		"MaxBytesMB":  h.maxBytesMB,
		"MaxBytesRaw": h.maxBytesRaw,
		"Images":      results,
		"IsUpload":    false,
		"IsBatch":     true,
	})
}

// buildErrorView builds an error view data map
func (h *WebHandler) buildErrorView(title, errorMsg, inputURL string, c *fiber.Ctx) fiber.Map {
	return fiber.Map{
		"Title":       title,
		"Error":       errorMsg,
		"InputURL":    inputURL,
		"BaseURL":     h.getBaseURL(c),
		"MaxBytesMB":  h.maxBytesMB,
		"MaxBytesRaw": h.maxBytesRaw,
		"Images":      []models.ImageResult{},
		"IsUpload":    false,
		"IsBatch":     false,
	}
}

// getBaseURL derives base URL from request
func (h *WebHandler) getBaseURL(c *fiber.Ctx) string {
	proto := c.Protocol()
	host := c.Hostname()

	if forwardedProto := c.Get("X-Forwarded-Proto"); forwardedProto != "" {
		proto = forwardedProto
	}
	if forwardedHost := c.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	if host == "" {
		host = "localhost:8080"
	}

	return fmt.Sprintf("%s://%s", proto, host)
}
