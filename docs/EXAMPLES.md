# Example Outputs

This document shows example outputs from the Image Metadata Viewer API and web interface.

## Web Interface Examples

### Home Page

The landing page features three tabs:

1. **URL Tab**: Single image URL input
2. **Multiple URLs Tab**: Textarea for batch processing (one URL per line)
3. **Upload Tab**: Drag & drop file upload interface

### Result Page

After processing, you'll see:

- Image preview (small size for quick loading)
- Comprehensive metadata in organized sections
- Beautiful card-based layout
- Color-coded badges for important values

## API Response Examples

### Single Image URL (GET)

**Request:**

```bash
curl "http://localhost:8080/api/https://example.com/sample.jpg"
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "fileName": "sample.jpg",
      "fileSize": 164000,
      "fileSizeHuman": "160.2 KiB",
      "fileType": "JPEG",
      "fileTypeExtension": "jpg",
      "mimeType": "image/jpeg",
      "source": "remote",
      "uploadedAt": "0001-01-01T00:00:00Z",
      "width": 720,
      "height": 1030,
      "aspectRatio": "0.699",
      "megapixels": 0.742,
      "colorSpace": "sRGB",
      "colorMode": "RGB",
      "colorComponents": 3,
      "samplesPerPixel": 3,
      "orientation": "Horizontal (normal)",
      "xResolution": 72,
      "yResolution": 72,
      "resolutionUnit": "inches",
      "software": "Adobe Photoshop 26.0 (Windows)",
      "modifyDate": "2025:10:31 21:41:50+0700",
      "createDate": "2025:09:29 11:02:25+0900",
      "creatorTool": "Adobe Photoshop 26.0 (Windows)",
      "format": "jpeg",
      "status": "200 OK",
      "finalURL": "https://example.com/sample.jpg",
      "contentLength": 164000,
      "lastModified": "Sat, 28 Dec 2024 20:59:30 GMT",
      "downloadedBytes": 164000,
      "truncated": false,
      "duration": "234ms"
    }
  ]
}
```

### Multiple URLs (POST with JSON)

**Request:**

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://example.com/image1.jpg",
      "https://example.com/image2.png"
    ]
  }'
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "fileName": "image1.jpg",
      "fileSize": 245000,
      "fileSizeHuman": "239.3 KiB",
      "fileType": "JPEG",
      "fileTypeExtension": "jpg",
      "mimeType": "image/jpeg",
      "width": 1920,
      "height": 1080,
      "aspectRatio": "1.778",
      "megapixels": 2.0736,
      "colorSpace": "sRGB",
      "source": "remote",
      "status": "200 OK",
      "duration": "156ms"
    },
    {
      "fileName": "image2.png",
      "fileSize": 512000,
      "fileSizeHuman": "500.0 KiB",
      "fileType": "PNG",
      "fileTypeExtension": "png",
      "mimeType": "image/png",
      "width": 1280,
      "height": 720,
      "aspectRatio": "1.778",
      "megapixels": 0.9216,
      "colorSpace": "sRGB",
      "source": "remote",
      "status": "200 OK",
      "duration": "203ms"
    }
  ],
  "errors": []
}
```

### Multiple File Upload (POST with Multipart)

**Request:**

```bash
curl -X POST http://localhost:8080/api \
  -F "files=@photo1.jpg" \
  -F "files=@photo2.png" \
  -F "files=@photo3.jpg"
```

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "fileName": "photo1.jpg",
      "fileSize": 189440,
      "fileSizeHuman": "185.0 KiB",
      "fileType": "JPEG",
      "fileTypeExtension": "jpg",
      "mimeType": "image/jpeg",
      "source": "upload",
      "uploadedAt": "2025-12-29T10:30:45.123456Z",
      "width": 1024,
      "height": 768,
      "aspectRatio": "1.333",
      "megapixels": 0.786432,
      "colorSpace": "sRGB",
      "colorMode": "RGB",
      "colorComponents": 3,
      "samplesPerPixel": 3,
      "orientation": "Horizontal (normal)",
      "xResolution": 96,
      "yResolution": 96,
      "resolutionUnit": "inches",
      "format": "jpeg"
    },
    {
      "fileName": "photo2.png",
      "fileSize": 456789,
      "fileSizeHuman": "446.1 KiB",
      "fileType": "PNG",
      "fileTypeExtension": "png",
      "mimeType": "image/png",
      "source": "upload",
      "uploadedAt": "2025-12-29T10:30:45.234567Z",
      "width": 800,
      "height": 600,
      "aspectRatio": "1.333",
      "megapixels": 0.48,
      "colorComponents": 3,
      "format": "png"
    },
    {
      "fileName": "photo3.jpg",
      "fileSize": 234567,
      "fileSizeHuman": "229.1 KiB",
      "fileType": "JPEG",
      "fileTypeExtension": "jpg",
      "mimeType": "image/jpeg",
      "source": "upload",
      "uploadedAt": "2025-12-29T10:30:45.345678Z",
      "width": 1920,
      "height": 1440,
      "aspectRatio": "1.333",
      "megapixels": 2.7648,
      "colorSpace": "sRGB",
      "orientation": "Horizontal (normal)",
      "software": "GIMP 2.10.32",
      "format": "jpeg"
    }
  ],
  "errors": []
}
```

### Error Response Examples

#### Invalid URL

```json
{
  "success": false,
  "error": "Invalid URL format"
}
```

#### Failed to Fetch

```json
{
  "success": false,
  "error": "fetch error: Get \"https://invalid-domain.com/image.jpg\": dial tcp: lookup invalid-domain.com: no such host"
}
```

#### File Too Large

```json
{
  "success": false,
  "data": [],
  "errors": ["large-image.jpg: file exceeds size limit"]
}
```

#### Mixed Success/Failure

```json
{
  "success": true,
  "data": [
    {
      "fileName": "valid-image.jpg",
      "fileSize": 123456,
      "width": 800,
      "height": 600
      // ... more metadata
    }
  ],
  "errors": [
    "https://broken-link.com/image.jpg: fetch error: 404 Not Found",
    "invalid-file.txt: decode error: image: unknown format"
  ]
}
```

## Metadata Fields Explanation

### Always Present

- `fileName`: Original filename or extracted from URL
- `fileSize`: Size in bytes
- `fileSizeHuman`: Human-readable size (e.g., "160.2 KiB")
- `fileType`: Image format (JPEG, PNG, GIF, etc.)
- `fileTypeExtension`: File extension
- `mimeType`: MIME type
- `source`: "remote" or "upload"

### Image Dimensions (if decodable)

- `width`: Width in pixels
- `height`: Height in pixels
- `aspectRatio`: Width/height ratio (3 decimal places)
- `megapixels`: Total megapixels

### Color Information (if available)

- `colorSpace`: Color space (e.g., "sRGB", "Uncalibrated")
- `colorMode`: Color mode (e.g., "RGB")
- `colorComponents`: Number of color components
- `samplesPerPixel`: Samples per pixel

### EXIF Data (if available)

- `orientation`: Image orientation
- `xResolution`: Horizontal resolution
- `yResolution`: Vertical resolution
- `resolutionUnit`: Unit for resolution (inches, centimeters)
- `software`: Software used to create/edit
- `createDate`: Creation date
- `modifyDate`: Modification date
- `creatorTool`: Creator tool/software

### Remote Image Only

- `status`: HTTP status code
- `finalURL`: Final URL after redirects
- `contentLength`: Content-Length header value
- `lastModified`: Last-Modified header value
- `downloadedBytes`: Bytes downloaded
- `truncated`: Whether download was truncated
- `duration`: Download duration

### Upload Only

- `uploadedAt`: Timestamp of upload

### Error Fields (if errors occurred)

- `fetchError`: Error fetching remote image
- `decodeError`: Error decoding image data

## Size Formatting Examples

| Bytes    | Human Readable |
| -------- | -------------- |
| 500      | 500 B          |
| 1024     | 1.0 KiB        |
| 1536     | 1.5 KiB        |
| 1048576  | 1.0 MiB        |
| 5242880  | 5.0 MiB        |
| 20971520 | 20.0 MiB       |

## Aspect Ratio Examples

| Dimensions | Aspect Ratio | Common Name     |
| ---------- | ------------ | --------------- |
| 1920×1080  | 1.778        | 16:9 (HD)       |
| 1280×720   | 1.778        | 16:9 (HD)       |
| 1024×768   | 1.333        | 4:3             |
| 1080×1920  | 0.562        | 9:16 (Portrait) |
| 1000×1000  | 1.000        | 1:1 (Square)    |
| 3840×2160  | 1.778        | 16:9 (4K)       |

## Megapixel Calculation Examples

| Dimensions | Calculation | Megapixels |
| ---------- | ----------- | ---------- |
| 1920×1080  | 2,073,600   | 2.07 MP    |
| 3840×2160  | 8,294,400   | 8.29 MP    |
| 800×600    | 480,000     | 0.48 MP    |
| 1024×768   | 786,432     | 0.79 MP    |

## Common Image Formats

| Format | Extension | MIME Type  |
| ------ | --------- | ---------- |
| JPEG   | jpg, jpeg | image/jpeg |
| PNG    | png       | image/png  |
| GIF    | gif       | image/gif  |
| BMP    | bmp       | image/bmp  |
| TIFF   | tif, tiff | image/tiff |
| WebP   | webp      | image/webp |

## Using the Examples

### JavaScript/Fetch

```javascript
const response = await fetch(
  "http://localhost:8080/api/https://example.com/image.jpg"
);
const result = await response.json();

if (result.success) {
  const metadata = result.data[0];
  console.log(`Size: ${metadata.fileSizeHuman}`);
  console.log(`Dimensions: ${metadata.width}×${metadata.height}`);
  console.log(`Aspect Ratio: ${metadata.aspectRatio}`);
}
```

### Python/Requests

```python
import requests

response = requests.get('http://localhost:8080/api/https://example.com/image.jpg')
data = response.json()

if data['success']:
    metadata = data['data'][0]
    print(f"Size: {metadata['fileSizeHuman']}")
    print(f"Dimensions: {metadata['width']}×{metadata['height']}")
    print(f"Megapixels: {metadata['megapixels']}")
```

### curl with jq

```bash
curl -s "http://localhost:8080/api/https://example.com/image.jpg" | \
  jq '.data[0] | {
    fileName,
    size: .fileSizeHuman,
    dimensions: "\(.width)×\(.height)",
    aspectRatio,
    megapixels
  }'
```

Output:

```json
{
  "fileName": "image.jpg",
  "size": "160.2 KiB",
  "dimensions": "720×1030",
  "aspectRatio": "0.699",
  "megapixels": 0.742
}
```
