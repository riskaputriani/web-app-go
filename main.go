package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	maxImageBytes = 20 << 20 // 20MB
)

var (
	homeTemplate = template.Must(template.ParseFiles("templates/home.html"))
	viewTemplate = template.Must(template.New("view.html").Funcs(template.FuncMap{
		"humanBytes": humanBytes,
	}).ParseFiles("templates/view.html"))
	httpClient = &http.Client{
		Timeout: 15 * time.Second,
	}
)

type ImageInfo struct {
	Status          string
	FinalURL        string
	ContentType     string
	ContentLength   int64
	DownloadedBytes int64
	Truncated       bool
	Format          string
	Width           int
	Height          int
	HasDimensions   bool
	DecodeError     string
	FetchError      string
	Duration        string
}

type ViewData struct {
	Title       string
	EmbedURL    template.URL
	InputURL    string
	DisplayURL  string
	Error       string
	Info        *ImageInfo
	MaxBytesMB  int
	MaxBytesRaw int64
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/go", handleForm)

	addr := ":8080"
	if port := strings.TrimSpace(os.Getenv("PORT")); port != "" {
		addr = ":" + port
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawURL := normalizeURL(strings.TrimSpace(r.URL.Query().Get("url")))
	if rawURL == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	escaped := url.PathEscape(rawURL)
	http.Redirect(w, r, "/"+escaped, http.StatusSeeOther)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		renderHome(w)
		return
	}

	rawPath := strings.TrimPrefix(r.URL.Path, "/")
	rawURL, err := url.PathUnescape(rawPath)
	if err != nil {
		renderView(w, ViewData{
			Title:       "Invalid URL",
			Error:       "Could not decode the URL path.",
			InputURL:    rawPath,
			MaxBytesMB:  maxImageBytes / (1 << 20),
			MaxBytesRaw: maxImageBytes,
		})
		return
	}

	normalizedURL := normalizeURL(rawURL)
	parsed, err := url.Parse(normalizedURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		renderView(w, ViewData{
			Title:       "Invalid URL",
			Error:       "The path must be a full http or https URL.",
			InputURL:    normalizedURL,
			MaxBytesMB:  maxImageBytes / (1 << 20),
			MaxBytesRaw: maxImageBytes,
		})
		return
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		renderView(w, ViewData{
			Title:       "Unsupported URL",
			Error:       "Only http and https URLs are supported.",
			InputURL:    normalizedURL,
			MaxBytesMB:  maxImageBytes / (1 << 20),
			MaxBytesRaw: maxImageBytes,
		})
		return
	}

	data := ViewData{
		Title:       "Image Preview",
		EmbedURL:    template.URL(parsed.String()),
		InputURL:    normalizedURL,
		DisplayURL:  parsed.String(),
		MaxBytesMB:  maxImageBytes / (1 << 20),
		MaxBytesRaw: maxImageBytes,
	}

	info := fetchImageInfo(r.Context(), parsed.String())
	data.Info = info
	if info.FetchError != "" {
		data.Error = info.FetchError
	}

	renderView(w, data)
}

// normalizeURL fixes single-slash schemes produced by path-based input.
func normalizeURL(raw string) string {
	normalized := strings.TrimSpace(raw)
	if strings.HasPrefix(normalized, "http:/") && !strings.HasPrefix(normalized, "http://") {
		return "http://" + strings.TrimPrefix(normalized, "http:/")
	}
	if strings.HasPrefix(normalized, "https:/") && !strings.HasPrefix(normalized, "https://") {
		return "https://" + strings.TrimPrefix(normalized, "https:/")
	}
	return normalized
}

// fetchImageInfo downloads a limited amount of data to inspect metadata safely.
func fetchImageInfo(ctx context.Context, imageURL string) *ImageInfo {
	info := &ImageInfo{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		info.FetchError = fmt.Sprintf("request error: %v", err)
		return info
	}
	req.Header.Set("User-Agent", "image-embedder/1.0")

	start := time.Now()
	resp, err := httpClient.Do(req)
	info.Duration = time.Since(start).Round(time.Millisecond).String()
	if err != nil {
		info.FetchError = fmt.Sprintf("fetch error: %v", err)
		return info
	}
	defer resp.Body.Close()

	info.Status = resp.Status
	info.ContentType = resp.Header.Get("Content-Type")
	info.ContentLength = resp.ContentLength
	info.FinalURL = resp.Request.URL.String()

	body, readErr := readLimited(resp.Body, maxImageBytes)
	if readErr != nil {
		info.FetchError = fmt.Sprintf("read error: %v", readErr)
		return info
	}

	info.DownloadedBytes = int64(len(body.Data))
	info.Truncated = body.Truncated

	if len(body.Data) == 0 {
		return info
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(body.Data))
	if err != nil {
		info.DecodeError = err.Error()
		return info
	}
	info.Format = format
	info.Width = cfg.Width
	info.Height = cfg.Height
	info.HasDimensions = true

	return info
}

type limitedRead struct {
	Data      []byte
	Truncated bool
}

// readLimited reads up to maxBytes to avoid unbounded downloads.
func readLimited(r io.Reader, maxBytes int64) (limitedRead, error) {
	if maxBytes <= 0 {
		return limitedRead{}, nil
	}

	limited := io.LimitReader(r, maxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return limitedRead{}, err
	}

	if int64(len(data)) > maxBytes {
		return limitedRead{Data: data[:maxBytes], Truncated: true}, nil
	}

	return limitedRead{Data: data, Truncated: false}, nil
}

func renderHome(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := homeTemplate.Execute(w, nil); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func renderView(w http.ResponseWriter, data ViewData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := viewTemplate.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func humanBytes(size int64) string {
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
