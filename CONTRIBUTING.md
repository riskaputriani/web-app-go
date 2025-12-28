# Contributing to Image Metadata Viewer

Thank you for your interest in contributing to Image Metadata Viewer! This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce**
- **Expected vs actual behavior**
- **Screenshots** (if applicable)
- **Environment details** (OS, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description**
- **Use case and motivation**
- **Proposed solution** (if you have one)
- **Alternative solutions** considered

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the code style** (run `make fmt` before committing)
3. **Add tests** for new features
4. **Update documentation** as needed
5. **Ensure all tests pass** (`make test`)
6. **Run linters** (`make lint`)
7. **Write clear commit messages**

## Development Setup

### Prerequisites

- Go 1.23 or higher
- Git
- Make (optional, but recommended)

### Setup Steps

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/image-metadata-viewer.git
cd image-metadata-viewer

# Add upstream remote
git remote add upstream https://github.com/ORIGINAL_OWNER/image-metadata-viewer.git

# Install dependencies
make deps

# Run tests
make test

# Run the application
make run
```

### Project Structure

```
.
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data structures
│   ├── services/        # Business logic
│   └── utils/           # Helper functions
├── pkg/                 # Public libraries
│   └── metadata/        # Metadata extraction
├── web/                 # Frontend assets
│   └── templates/       # HTML templates
├── docs/                # Documentation
└── .github/             # GitHub configuration
```

## Coding Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `go vet` to catch common mistakes
- Use meaningful variable and function names
- Add comments for exported functions

### Example Code Style

```go
// ExtractMetadata extracts comprehensive metadata from image data.
// It returns an ImageMetadata struct containing all extracted information.
func ExtractMetadata(data []byte, contentType, fileName string) *models.ImageMetadata {
    meta := &models.ImageMetadata{
        FileName:          fileName,
        FileSize:          int64(len(data)),
        FileSizeHuman:     utils.HumanBytes(int64(len(data))),
        MIMEType:          contentType,
    }

    // Extract basic image information
    cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
    if err != nil {
        meta.DecodeError = err.Error()
        return meta
    }

    meta.Format = format
    meta.Width = cfg.Width
    meta.Height = cfg.Height

    return meta
}
```

### Testing

- Write unit tests for new features
- Aim for >80% code coverage
- Use table-driven tests where appropriate

```go
func TestHumanBytes(t *testing.T) {
    tests := []struct {
        name     string
        input    int64
        expected string
    }{
        {"zero", 0, "0 B"},
        {"bytes", 500, "500 B"},
        {"kilobytes", 1536, "1.5 KiB"},
        {"megabytes", 1048576, "1.0 MiB"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := HumanBytes(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

### Documentation

- Add godoc comments for all exported functions
- Update README.md for user-facing changes
- Update API.md for API changes
- Include examples in documentation

## Commit Messages

Use clear and descriptive commit messages:

```
feat: add support for WebP images
fix: correct aspect ratio calculation for portrait images
docs: update API documentation with new endpoints
refactor: simplify metadata extraction logic
test: add tests for image service
chore: update dependencies
```

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test -v ./internal/services/

# Run specific test
go test -v -run TestExtractMetadata ./pkg/metadata/
```

## Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run without building
make run

# Run with auto-reload (requires air)
make dev
```

## Questions?

Feel free to open an issue for:

- Questions about the codebase
- Clarification on contribution guidelines
- Discussion about new features

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions make this project better for everyone. Thank you for taking the time to contribute!
