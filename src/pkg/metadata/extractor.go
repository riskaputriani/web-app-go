package metadata

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"time"

	"github.com/ahrdadan/image-metadata-viewer/src/internal/models"
	"github.com/ahrdadan/image-metadata-viewer/src/internal/utils"
	"github.com/rwcarlsen/goexif/exif"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// ExtractMetadata extracts comprehensive metadata from image data
func ExtractMetadata(data []byte, contentType, fileName string) *models.ImageMetadata {
	meta := &models.ImageMetadata{
		FileName:          fileName,
		FileSize:          int64(len(data)),
		FileSizeHuman:     utils.HumanBytes(int64(len(data))),
		MIMEType:          contentType,
		FileTypeExtension: utils.ExtensionFromName(fileName),
		UploadedAt:        time.Now(),
	}

	// Decode image config for basic dimensions
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		meta.DecodeError = err.Error()
		return meta
	}

	meta.Format = format
	meta.Width = cfg.Width
	meta.Height = cfg.Height
	meta.AspectRatio = utils.CalculateAspectRatio(cfg.Width, cfg.Height)
	meta.AspectRatioFraction = utils.CalculateAspectRatioFraction(cfg.Width, cfg.Height)
	meta.Megapixels = utils.CalculateMegapixels(cfg.Width, cfg.Height)
	meta.FileType = strings.ToUpper(format)

	if meta.FileTypeExtension == "" {
		meta.FileTypeExtension = utils.FormatToExtension(format)
	}

	// Extract EXIF data
	extractEXIF(data, meta)

	// Set color space information
	if cfg.ColorModel != nil {
		meta.ColorComponents = 3 // Most images have RGB
		meta.SamplesPerPixel = 3
	}

	return meta
}

// extractEXIF extracts EXIF metadata from image data
func extractEXIF(data []byte, meta *models.ImageMetadata) {
	x, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		// EXIF not available or couldn't decode
		return
	}

	// Orientation
	if tag, err := x.Get(exif.Orientation); err == nil {
		if val, err := tag.Int(0); err == nil {
			meta.Orientation = orientationToString(val)
		}
	}

	// Software
	if tag, err := x.Get(exif.Software); err == nil {
		if val, err := tag.StringVal(); err == nil {
			meta.Software = val
			meta.CreatorTool = val
		}
	}

	// DateTime
	if tag, err := x.Get(exif.DateTime); err == nil {
		if val, err := tag.StringVal(); err == nil {
			meta.ModifyDate = val
		}
	}

	// DateTimeOriginal
	if tag, err := x.Get(exif.DateTimeOriginal); err == nil {
		if val, err := tag.StringVal(); err == nil {
			meta.CreateDate = val
		}
	}

	// XResolution
	if tag, err := x.Get(exif.XResolution); err == nil {
		if numer, denom, err := tag.Rat2(0); err == nil && denom != 0 {
			meta.XResolution = int(numer / denom)
		}
	}

	// YResolution
	if tag, err := x.Get(exif.YResolution); err == nil {
		if numer, denom, err := tag.Rat2(0); err == nil && denom != 0 {
			meta.YResolution = int(numer / denom)
		}
	}

	// ResolutionUnit
	if tag, err := x.Get(exif.ResolutionUnit); err == nil {
		if val, err := tag.Int(0); err == nil {
			meta.ResolutionUnit = resolutionUnitToString(val)
		}
	}

	// ColorSpace
	if tag, err := x.Get(exif.ColorSpace); err == nil {
		if val, err := tag.Int(0); err == nil {
			if val == 1 {
				meta.ColorSpace = "sRGB"
				meta.ColorMode = "RGB"
			} else {
				meta.ColorSpace = "Uncalibrated"
			}
		}
	}
}

// orientationToString converts EXIF orientation value to string
func orientationToString(orientation int) string {
	switch orientation {
	case 1:
		return "Horizontal (normal)"
	case 2:
		return "Mirror horizontal"
	case 3:
		return "Rotate 180"
	case 4:
		return "Mirror vertical"
	case 5:
		return "Mirror horizontal and rotate 270 CW"
	case 6:
		return "Rotate 90 CW"
	case 7:
		return "Mirror horizontal and rotate 90 CW"
	case 8:
		return "Rotate 270 CW"
	default:
		return fmt.Sprintf("Unknown (%d)", orientation)
	}
}

// resolutionUnitToString converts EXIF resolution unit to string
func resolutionUnitToString(unit int) string {
	switch unit {
	case 2:
		return "inches"
	case 3:
		return "centimeters"
	default:
		return "unknown"
	}
}
