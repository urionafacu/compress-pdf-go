package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

func CompressImage(base64Image string) ([]byte, error) {
	imgBytes, err := base64.StdEncoding.DecodeString(base64Image)

	if err != nil {
		return nil, err
	}

	img, format, err := image.Decode(bytes.NewReader(imgBytes))

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	maxSize := 512 * 1024 // 500 KB
	quality := 100
	reduced := false

	for {
		buf.Reset()
		fmt.Println("quality", quality)

		switch format {
		case "jpeg", "jpg":
			img = resize.Resize(0, 0, img, resize.Lanczos3)
			err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		case "png":
			img = resize.Resize(0, 0, img, resize.Lanczos3)
			encoder := png.Encoder{CompressionLevel: png.BestCompression}
			err = encoder.Encode(&buf, img)
		default:
			return nil, fmt.Errorf("Unsupported image format: %s", format)
		}

		if err != nil {
			return nil, err
		}

		fmt.Println("bufLen", buf.Len())

		if buf.Len() <= maxSize || reduced {
			break
		}

		quality -= 10
		if quality < 0 {
			quality = 0
			reduced = true
		}
	}

	return buf.Bytes(), nil
}
