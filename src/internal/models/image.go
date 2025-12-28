package models

import "time"

// ImageMetadata contains comprehensive metadata information about an image
type ImageMetadata struct {
	// Basic file information
	FileName          string    `json:"fileName"`
	FileSize          int64     `json:"fileSize"`
	FileSizeHuman     string    `json:"fileSizeHuman"`
	FileType          string    `json:"fileType"`
	FileTypeExtension string    `json:"fileTypeExtension"`
	MIMEType          string    `json:"mimeType"`
	Source            string    `json:"source"` // "upload" or "remote"
	UploadedAt        time.Time `json:"uploadedAt,omitempty"`

	// Image dimensions
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	AspectRatio string  `json:"aspectRatio"`
	Megapixels  float64 `json:"megapixels"`

	// Color information
	ColorSpace      string `json:"colorSpace,omitempty"`
	ColorMode       string `json:"colorMode,omitempty"`
	BitsPerSample   string `json:"bitsPerSample,omitempty"`
	ColorComponents int    `json:"colorComponents,omitempty"`
	SamplesPerPixel int    `json:"samplesPerPixel,omitempty"`

	// JPEG specific
	EncodingProcess  string `json:"encodingProcess,omitempty"`
	PhotoshopQuality int    `json:"photoshopQuality,omitempty"`

	// EXIF data
	Orientation    string `json:"orientation,omitempty"`
	XResolution    int    `json:"xResolution,omitempty"`
	YResolution    int    `json:"yResolution,omitempty"`
	ResolutionUnit string `json:"resolutionUnit,omitempty"`
	Software       string `json:"software,omitempty"`
	ModifyDate     string `json:"modifyDate,omitempty"`
	CreateDate     string `json:"createDate,omitempty"`

	// XMP metadata
	CreatorTool  string `json:"creatorTool,omitempty"`
	MetadataDate string `json:"metadataDate,omitempty"`
	Format       string `json:"format"`

	// HTTP metadata (for remote images)
	Status          string `json:"status,omitempty"`
	FinalURL        string `json:"finalURL,omitempty"`
	ContentLength   int64  `json:"contentLength,omitempty"`
	LastModified    string `json:"lastModified,omitempty"`
	DownloadedBytes int64  `json:"downloadedBytes,omitempty"`
	Truncated       bool   `json:"truncated,omitempty"`
	Duration        string `json:"duration,omitempty"`

	// Error information
	FetchError  string `json:"fetchError,omitempty"`
	DecodeError string `json:"decodeError,omitempty"`
}

// ViewData represents the data passed to view templates
type ViewData struct {
	Title       string
	BaseURL     string
	MaxBytesMB  int
	MaxBytesRaw int64
	Error       string
	Images      []ImageResult
	IsUpload    bool
	IsBatch     bool
}

// ImageResult represents a single image processing result
type ImageResult struct {
	InputURL   string
	DisplayURL string
	EmbedURL   string
	Metadata   *ImageMetadata
	Error      string
	IsBlob     bool
}

// HomeData represents the data passed to home template
type HomeData struct {
	BaseURL string
}

// APIResponse represents the JSON response for API endpoints
type APIResponse struct {
	Success bool            `json:"success"`
	Data    []ImageMetadata `json:"data,omitempty"`
	Errors  []string        `json:"errors,omitempty"`
	Message string          `json:"message,omitempty"`
}

// APIErrorResponse represents an error response
type APIErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
