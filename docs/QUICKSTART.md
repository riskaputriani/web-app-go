# Quick Start Guide

Get up and running with Image Metadata Viewer in 5 minutes!

## Prerequisites

- Go 1.23+ installed
- Git

## Installation

### Option 1: Quick Run (Recommended for Testing)

```bash
# Clone the repository
git clone https://github.com/yourusername/image-metadata-viewer.git
cd image-metadata-viewer

# Download dependencies
go mod download

# Run directly
go run ./cmd/server
```

Visit `http://localhost:8080` in your browser!

### Option 2: Build and Run

```bash
# Clone the repository
git clone https://github.com/yourusername/image-metadata-viewer.git
cd image-metadata-viewer

# Install dependencies
make deps

# Build the application
make build

# Run the built binary
./build/image-metadata-viewer
```

### Option 3: Using Docker

```bash
# Clone the repository
git clone https://github.com/yourusername/image-metadata-viewer.git
cd image-metadata-viewer

# Build Docker image
docker build -t image-metadata-viewer .

# Run container
docker run -p 8080:8080 image-metadata-viewer
```

Visit `http://localhost:8080`!

### Option 4: Docker without Clone

```bash
# Pull and run (once published)
docker pull username/image-metadata-viewer
docker run -p 8080:8080 username/image-metadata-viewer
```

## First Steps

### 1. Process a Single Image URL

1. Open http://localhost:8080
2. Click the "URL" tab
3. Enter an image URL: `https://picsum.photos/800/600`
4. Click "View Metadata"

### 2. Process Multiple URLs

1. Click the "Multiple URLs" tab
2. Enter multiple URLs, one per line:
   ```
   https://picsum.photos/800/600
   https://picsum.photos/1920/1080
   https://picsum.photos/400/300
   ```
3. Click "Process Batch"

### 3. Upload Files

1. Click the "Upload" tab
2. Drag and drop image files, or click to select
3. Click "Upload & View Metadata"

### 4. Use the REST API

```bash
# Get metadata from a URL
curl "http://localhost:8080/api/https://picsum.photos/800/600"

# Upload a file
curl -X POST http://localhost:8080/api -F "files=@your-image.jpg"

# Process multiple URLs
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://picsum.photos/800/600",
      "https://picsum.photos/1920/1080"
    ]
  }'
```

## Configuration

### Change Port

```bash
# Using environment variable
PORT=3000 go run ./cmd/server

# Or with Make
PORT=3000 make run
```

### Production Deployment

```bash
# Build optimized binary
make build

# Run in production
PORT=8080 ./build/image-metadata-viewer
```

## Common Issues

### Port Already in Use

```bash
# Change the port
PORT=3000 make run
```

### Module Errors

```bash
# Clean and reinstall dependencies
go clean -modcache
go mod download
go mod tidy
```

### Template Errors

Make sure the `web/templates/` directory exists with `home.html` and `view.html`.

## Next Steps

- Read the [full documentation](README.md)
- Check out the [API documentation](docs/API.md)
- Learn about [contributing](CONTRIBUTING.md)
- View the [project structure](docs/PROJECT_SUMMARY.md)

## Need Help?

- Open an issue on GitHub
- Check existing issues for solutions
- Read the documentation thoroughly

## Example Commands

```bash
# Development workflow
make deps          # Install dependencies
make fmt           # Format code
make test          # Run tests
make run           # Run application
make build         # Build binary

# Production workflow
make build-all     # Build for all platforms
make docker-build  # Build Docker image
make clean         # Clean build artifacts

# Quality checks
make lint          # Run linter
make vet           # Run go vet
make check         # Run all checks
```

## Testing the API

### Using curl

```bash
# Single URL
curl "http://localhost:8080/api/https://picsum.photos/800/600"

# Pretty print with jq (if installed)
curl "http://localhost:8080/api/https://picsum.photos/800/600" | jq

# Multiple files
curl -X POST http://localhost:8080/api \
  -F "files=@image1.jpg" \
  -F "files=@image2.png"
```

### Using httpie (if installed)

```bash
# Single URL
http GET localhost:8080/api/https://picsum.photos/800/600

# Multiple URLs
http POST localhost:8080/api urls:='["https://picsum.photos/800/600"]'

# File upload
http --form POST localhost:8080/api files@image.jpg
```

### Using Postman

1. Import the API collection (if provided)
2. Or manually create requests:
   - GET: `http://localhost:8080/api/YOUR_IMAGE_URL`
   - POST: `http://localhost:8080/api` with form-data or JSON body

## Performance Tips

1. **For remote URLs**: The first request might be slower due to download
2. **For large images**: They will be truncated at 20MB
3. **Batch processing**: Process multiple images in one request for efficiency
4. **Caching**: Currently not implemented, each request downloads fresh

## Security Notes

- Images are **never stored** on the server
- Processing happens in memory only
- 20MB size limit per image
- 15-second timeout for remote URLs
- No authentication required (can be added if needed)

## What You Can Extract

- **File Info**: Name, size, type, extension, MIME type
- **Dimensions**: Width, height, aspect ratio, megapixels
- **EXIF Data**: Orientation, resolution, software, dates
- **Color Info**: Color space, mode, components
- **HTTP Info** (remote): Status, headers, download time

Enjoy using Image Metadata Viewer! 🎉
