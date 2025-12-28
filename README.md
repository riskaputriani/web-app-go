# Image Metadata Viewer

A modern, fast web application for extracting comprehensive metadata from images. Built with Go and Fiber framework, supporting both remote URL fetching and local file uploads.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52.5-00ACD7?style=flat)](https://gofiber.io/)
[![CI](https://github.com/ahrdadan/image-metadata-viewer/actions/workflows/ci-cd.yml/badge.svg)](https://github.com/ahrdadan/image-metadata-viewer/actions/workflows/ci-cd.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## ✨ Features

- 🖼️ **Multiple Input Methods**

  - Single or batch URL processing
  - Drag & drop file upload
  - Multiple file upload support
  - Textarea input for batch URLs (newline separated)

- 📊 **Comprehensive Metadata Extraction**

  - Image dimensions (width, height, aspect ratio, megapixels)
  - File information (size, type, MIME type)
  - EXIF data (orientation, resolution, software, dates)
  - Color space information
  - XMP metadata support
  - HTTP headers for remote images

- 🚀 **REST API**

  - GET endpoint for single URL metadata
  - POST endpoint for batch processing
  - JSON response format
  - Support for both URLs and file uploads

- 🎨 **Modern UI**

  - Responsive design
  - Drag & drop interface
  - Real-time file preview
  - Beautiful gradient design
  - Mobile-friendly

- 🔒 **Privacy & Security**
  - No data storage
  - Images processed in memory only
  - 20MB size limit per image
  - Request timeout protection

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or higher
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/ahrdadan/image-metadata-viewer.git
cd image-metadata-viewer

# Download dependencies
make deps

# Run the application
make run
```

The server will start at `http://localhost:8080`

### Using Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

## 📖 Usage

### Web Interface

1. **Single URL**: Enter an image URL in the URL tab
2. **Multiple URLs**: Enter multiple URLs (one per line) in the Multiple URLs tab
3. **Upload Files**: Drag & drop files or click to select in the Upload tab

### REST API

#### Get Metadata from URL

```bash
# GET request
curl "http://localhost:8080/api/https://example.com/image.jpg"
```

#### Batch Process URLs

```bash
# POST request with JSON
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.jpg"
    ]
  }'
```

#### Upload Files

```bash
# POST request with files
curl -X POST http://localhost:8080/api \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg"
```

#### API Response Format

```json
{
  "success": true,
  "data": [
    {
      "fileName": "image.jpg",
      "fileSize": 164000,
      "fileSizeHuman": "160.2 KiB",
      "fileType": "JPEG",
      "fileTypeExtension": "jpg",
      "mimeType": "image/jpeg",
      "width": 720,
      "height": 1030,
      "aspectRatio": "0.699",
      "megapixels": 0.742,
      "colorSpace": "sRGB",
      "orientation": "Horizontal (normal)",
      "xResolution": 72,
      "yResolution": 72,
      "resolutionUnit": "inches",
      "software": "Adobe Photoshop 26.0 (Windows)",
      "createDate": "2025:09:29 11:02:25+0900",
      "modifyDate": "2025:10:31 21:41:50+0700"
    }
  ]
}
```

## 🏗️ Project Structure

```
.
├── cmd/
│   └── server/          # Main application entry point
│       └── main.go
├── internal/            # Private application code
│   ├── handlers/        # HTTP handlers
│   │   ├── web_handler.go
│   │   └── api_handler.go
│   ├── models/          # Data models
│   │   └── image.go
│   ├── services/        # Business logic
│   │   └── image_service.go
│   └── utils/           # Utility functions
│       └── helpers.go
├── pkg/                 # Public libraries
│   └── metadata/        # Metadata extraction
│       └── extractor.go
├── web/                 # Web assets
│   └── templates/       # HTML templates
│       ├── home.html
│       └── view.html
├── docs/                # Documentation
├── .github/             # GitHub configurations
│   └── workflows/       # CI/CD workflows
├── Dockerfile           # Docker configuration
├── Makefile            # Build automation
├── go.mod              # Go module definition
└── README.md           # This file
```

## 🛠️ Development

### Available Make Commands

```bash
make help           # Show all available commands
make run            # Run the application
make build          # Build binary
make build-all      # Build for all platforms
make test           # Run tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make deps           # Download dependencies
make fmt            # Format code
make lint           # Run linter
make vet            # Run go vet
make dev            # Run with auto-reload (requires air)
make check          # Run all checks
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run all checks
make check
```

## 📝 Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)

Example:

```bash
PORT=3000 make run
```

## 🐳 Docker Deployment

### Build and Run

```bash
# Build image
docker build -t image-metadata-viewer .

# Run container
docker run -p 8080:8080 image-metadata-viewer
```

### Docker Compose

```yaml
version: "3.8"
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Express-inspired web framework
- [goexif](https://github.com/rwcarlsen/goexif) - EXIF extraction library
- [golang.org/x/image](https://pkg.go.dev/golang.org/x/image) - Extended image support

## 📧 Contact

Project Link: [https://github.com/ahrdadan/image-metadata-viewer.git](https://github.com/ahrdadan/image-metadata-viewer.git)

## 🗺️ Roadmap

- [ ] Add support for more image formats (RAW, HEIC)
- [ ] Implement caching for remote URLs
- [ ] Add image comparison feature
- [ ] Export metadata to various formats (JSON, CSV, XML)
- [ ] Add batch download for remote images
- [ ] Implement rate limiting for API
- [ ] Add authentication for API endpoints
- [ ] Create comprehensive test suite
- [ ] Add internationalization support

---

Made with ❤️ using Go and Fiber
