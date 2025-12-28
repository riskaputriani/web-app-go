# API Documentation

## Overview

Image Metadata Viewer provides a RESTful API for extracting comprehensive metadata from images. The API supports both remote URL fetching and file uploads.

## Base URL

```
http://localhost:8080/api
```

## Authentication

Currently, the API does not require authentication. This may change in future versions.

## Endpoints

### GET /api/{url}

Extract metadata from a remote image URL.

**Parameters:**

- `url` (path parameter, required): The full URL of the image (not percent-encoded)

**Example:**

```bash
curl "http://localhost:8080/api/https://example.com/image.jpg"
```

**Response:**

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
      "modifyDate": "2025:10:31 21:41:50+0700",
      "source": "remote",
      "status": "200 OK",
      "duration": "234ms"
    }
  ]
}
```

### POST /api

Process multiple URLs or upload files to extract metadata.

#### With JSON (Multiple URLs)

**Headers:**

- `Content-Type: application/json`

**Body:**

```json
{
  "urls": ["https://example.com/image1.jpg", "https://example.com/image2.jpg"]
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.jpg"
    ]
  }'
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      /* metadata for image1 */
    },
    {
      /* metadata for image2 */
    }
  ],
  "errors": []
}
```

#### With Multipart Form (File Upload)

**Headers:**

- `Content-Type: multipart/form-data`

**Form Fields:**

- `files` (file, multiple allowed): Image files to process

**Example:**

```bash
curl -X POST http://localhost:8080/api \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg" \
  -F "files=@image3.png"
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "fileName": "image1.jpg",
      "fileSize": 123456,
      "fileSizeHuman": "120.6 KiB",
      "width": 1920,
      "height": 1080,
      "source": "upload"
      /* ... more metadata ... */
    },
    {
      "fileName": "image2.jpg"
      /* ... metadata ... */
    },
    {
      "fileName": "image3.png"
      /* ... metadata ... */
    }
  ],
  "errors": []
}
```

## Response Format

### Success Response

```json
{
  "success": true,
  "data": [
    /* array of ImageMetadata objects */
  ],
  "errors": [
    /* array of error messages, if any */
  ]
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error description"
}
```

## Metadata Fields

| Field               | Type    | Description                    |
| ------------------- | ------- | ------------------------------ |
| `fileName`          | string  | Original filename              |
| `fileSize`          | int64   | File size in bytes             |
| `fileSizeHuman`     | string  | Human-readable file size       |
| `fileType`          | string  | Image format (JPEG, PNG, etc.) |
| `fileTypeExtension` | string  | File extension                 |
| `mimeType`          | string  | MIME type                      |
| `width`             | int     | Image width in pixels          |
| `height`            | int     | Image height in pixels         |
| `aspectRatio`       | string  | Width/height ratio             |
| `megapixels`        | float64 | Total megapixels               |
| `colorSpace`        | string  | Color space (sRGB, etc.)       |
| `colorMode`         | string  | Color mode (RGB, etc.)         |
| `orientation`       | string  | EXIF orientation               |
| `xResolution`       | int     | Horizontal resolution          |
| `yResolution`       | int     | Vertical resolution            |
| `resolutionUnit`    | string  | Resolution unit (inches, cm)   |
| `software`          | string  | Software used to create/edit   |
| `createDate`        | string  | Creation date from EXIF        |
| `modifyDate`        | string  | Modification date from EXIF    |
| `source`            | string  | "remote" or "upload"           |
| `status`            | string  | HTTP status (for remote)       |
| `duration`          | string  | Download duration (for remote) |

## Rate Limits

Currently no rate limits are enforced. This may change in future versions.

## Error Codes

| HTTP Status | Description                                |
| ----------- | ------------------------------------------ |
| 200         | Success                                    |
| 400         | Bad Request (invalid parameters)           |
| 502         | Bad Gateway (failed to fetch remote image) |
| 500         | Internal Server Error                      |

## Size Limits

- Maximum file size: **20 MB** per image
- Maximum upload size: **21 MB** total per request

## Supported Formats

- JPEG / JPG
- PNG
- GIF
- BMP
- TIFF
- WebP

## Examples

### JavaScript/Fetch

```javascript
// Single URL
const response = await fetch(
  "http://localhost:8080/api/https://example.com/image.jpg"
);
const data = await response.json();
console.log(data.data[0]);

// Multiple URLs
const response = await fetch("http://localhost:8080/api", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    urls: ["https://example.com/image1.jpg", "https://example.com/image2.jpg"],
  }),
});
const data = await response.json();

// File upload
const formData = new FormData();
formData.append("files", fileInput.files[0]);
formData.append("files", fileInput.files[1]);

const response = await fetch("http://localhost:8080/api", {
  method: "POST",
  body: formData,
});
const data = await response.json();
```

### Python/Requests

```python
import requests

# Single URL
response = requests.get('http://localhost:8080/api/https://example.com/image.jpg')
data = response.json()
print(data['data'][0])

# Multiple URLs
response = requests.post('http://localhost:8080/api', json={
    'urls': [
        'https://example.com/image1.jpg',
        'https://example.com/image2.jpg'
    ]
})
data = response.json()

# File upload
files = {
    'files': [
        ('files', open('image1.jpg', 'rb')),
        ('files', open('image2.jpg', 'rb'))
    ]
}
response = requests.post('http://localhost:8080/api', files=files)
data = response.json()
```

### cURL

```bash
# Single URL
curl "http://localhost:8080/api/https://example.com/image.jpg"

# Multiple URLs
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{"urls":["https://example.com/image1.jpg","https://example.com/image2.jpg"]}'

# File upload
curl -X POST http://localhost:8080/api \
  -F "files=@image1.jpg" \
  -F "files=@image2.jpg"
```

## Notes

- Images are processed in memory and never stored on the server
- EXIF data availability depends on the image file
- Remote URLs must be publicly accessible
- Large images may be truncated at 20MB during processing
