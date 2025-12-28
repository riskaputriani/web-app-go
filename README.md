# Image Embed Tester (Go)

Simple web app to test embedding images by URL. You can:
- Submit the form at the root (`/`) to open the embed page.
- Access a URL directly in the path, e.g. `/https:/down-id.img.susercontent.com/file/id-11134004-82251-mit5169olsld91`.
- See image info: status, content-type, size, downloaded bytes, format, and dimensions.

## Features
- Auto-embed `<img src="URL">` from the path.
- URL normalization for `http:/` and `https:/` (no manual encoding needed).
- Image metadata from HTTP headers and `image.DecodeConfig`.

## Project Layout
- `main.go` Go server.
- `templates/home.html` form page.
- `templates/view.html` embed + image info page.
- `Dockerfile` ready for deployment (port 8080).

## Run Locally
```bash
go run .
```
Open:
- `http://localhost:8080/`
- `http://localhost:8080/https:/down-id.img.susercontent.com/file/id-11134004-82251-mit5169olsld91`

## Lint & Checks
```bash
go vet ./...
```
Optional (requires golangci-lint):
```bash
golangci-lint run
```

## Docker
Build and run:
```bash
docker build -t image-embed .
docker run --rm -p 8080:8080 image-embed
```

## Environment
- `PORT` (optional): server port. Default `8080`.

## Notes
- For URLs with query strings, use the form on the root page to ensure safe encoding.
