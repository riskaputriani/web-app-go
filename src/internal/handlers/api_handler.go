package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/ahrdadan/image-metadata-viewer/src/internal/models"
	"github.com/ahrdadan/image-metadata-viewer/src/internal/services"
	"github.com/ahrdadan/image-metadata-viewer/src/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// APIHandler handles REST API requests
type APIHandler struct {
	imageService *services.ImageService
}

// NewAPIHandler creates a new APIHandler
func NewAPIHandler(imageService *services.ImageService) *APIHandler {
	return &APIHandler{
		imageService: imageService,
	}
}

// HandleGetMetadata handles GET /api/* for URL metadata retrieval
func (h *APIHandler) HandleGetMetadata(c *fiber.Ctx) error {
	rawPath := c.Params("*")
	if rawPath == "" {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "URL parameter is required",
		})
	}

	// Decode URL (it's not percent-encoded in this endpoint)
	imageURL := utils.NormalizeURL(rawPath)

	// Validate URL
	parsed, err := url.Parse(imageURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "Invalid URL format",
		})
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "Only http and https URLs are supported",
		})
	}

	// Process the URL
	meta := h.imageService.ProcessRemoteURL(c.Context(), parsed.String())

	if meta.FetchError != "" {
		return c.Status(http.StatusBadGateway).JSON(models.APIErrorResponse{
			Success: false,
			Error:   meta.FetchError,
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    []models.ImageMetadata{*meta},
	})
}

// HandlePostMetadata handles POST /api for multiple URLs or file uploads
func (h *APIHandler) HandlePostMetadata(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")

	// Check if it's a multipart form (file upload)
	if strings.Contains(contentType, "multipart/form-data") {
		return h.handleFileUpload(c)
	}

	// Check if it's JSON (URL list)
	if strings.Contains(contentType, "application/json") {
		return h.handleJSONURLs(c)
	}

	return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
		Success: false,
		Error:   "Content-Type must be application/json or multipart/form-data",
	})
}

// handleJSONURLs processes JSON payload with URLs
func (h *APIHandler) handleJSONURLs(c *fiber.Ctx) error {
	var payload struct {
		URLs []string `json:"urls"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "Invalid JSON payload",
		})
	}

	if len(payload.URLs) == 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "No URLs provided",
		})
	}

	results := make([]models.ImageMetadata, 0, len(payload.URLs))
	errors := make([]string, 0)

	for _, rawURL := range payload.URLs {
		normalized := utils.NormalizeURL(strings.TrimSpace(rawURL))

		parsed, err := url.Parse(normalized)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			errors = append(errors, fmt.Sprintf("Invalid URL: %s", rawURL))
			continue
		}

		meta := h.imageService.ProcessRemoteURL(c.Context(), parsed.String())

		if meta.FetchError != "" {
			errors = append(errors, fmt.Sprintf("%s: %s", rawURL, meta.FetchError))
		}

		results = append(results, *meta)
	}

	response := models.APIResponse{
		Success: len(results) > 0,
		Data:    results,
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	return c.JSON(response)
}

// handleFileUpload processes multipart file uploads
func (h *APIHandler) handleFileUpload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "Could not parse multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "No files provided",
		})
	}

	results := make([]models.ImageMetadata, 0, len(files))
	errors := make([]string, 0)

	for _, fileHeader := range files {
		meta, err := h.processAPIUpload(fileHeader)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %s", fileHeader.Filename, err.Error()))
			continue
		}
		results = append(results, *meta)
	}

	if len(results) == 0 {
		return c.Status(http.StatusBadRequest).JSON(models.APIErrorResponse{
			Success: false,
			Error:   "No valid images processed",
		})
	}

	response := models.APIResponse{
		Success: true,
		Data:    results,
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	return c.JSON(response)
}

// processAPIUpload processes a single uploaded file for API
func (h *APIHandler) processAPIUpload(fileHeader *multipart.FileHeader) (*models.ImageMetadata, error) {
	if fileHeader.Size > services.MaxUploadBytes {
		return nil, fmt.Errorf("file exceeds size limit")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	data := make([]byte, fileHeader.Size)
	_, err = file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	meta := h.imageService.ProcessUpload(data, contentType, fileHeader.Filename)

	if meta.DecodeError != "" {
		return nil, fmt.Errorf("decode error: %s", meta.DecodeError)
	}

	return meta, nil
}
