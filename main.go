package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
	"github.com/jung-kurt/gofpdf"
	"github.com/nfnt/resize"
)

func pdfToImages(pdfBase64 string) ([]string, error) {
	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return nil, err
	}

	tmpFile, err := ioutil.TempFile("", "temp_pdf_*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(pdfBytes); err != nil {
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	doc, err := fitz.New(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	var base64Images []string
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}

		base64Images = append(base64Images, base64.StdEncoding.EncodeToString(buf.Bytes()))
	}

	return base64Images, nil
}

func compressImage(base64Image string, format string) ([]byte, error) {
	imgBytes, err := base64.StdEncoding.DecodeString(base64Image)

	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))

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
		fmt.Println("bufLen", buf.Len())

		if buf.Len() <= maxSize || reduced {
			break
		}

		switch format {
		case "jpeg", "jpg":
			img = resize.Resize(0, 0, img, resize.Lanczos3)
			err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		case "png":
			img = resize.Resize(0, 0, img, resize.Lanczos3)
			err = png.Encode(&buf, img)
		default:
			return nil, fmt.Errorf("Unsupported image format: %s", format)
		}

		if err != nil {
			return nil, err
		}

		println(buf.Len(), quality)

		quality -= 10
		if quality < 0 {
			quality = 0
			reduced = true
		}
	}

	return buf.Bytes(), nil
}

func imagesToPdf(base64Images []string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 0)

	for _, base64Image := range base64Images {
		imgBytes, err := base64.StdEncoding.DecodeString(base64Image)
		if err != nil {
			return "", err
		}

		img, err := jpeg.Decode(bytes.NewReader(imgBytes))
		if err != nil {
			return "", err
		}

		imgWidth := img.Bounds().Dx()
		imgHeight := img.Bounds().Dy()

		pdfW, pdfH := pdf.GetPageSize()

		scale := pdfW / float64(imgWidth)
		scaledHeight := float64(imgHeight) * scale

		if scaledHeight > pdfH {
			scale = pdfH / float64(imgHeight)
			pdfW = float64(imgWidth) * scale
		} else {
			pdfH = scaledHeight
		}

		tmpFile, err := ioutil.TempFile("", "temp_image_*.jpg")
		if err != nil {
			return "", err
		}
		defer os.Remove(tmpFile.Name())

		if err := jpeg.Encode(tmpFile, img, nil); err != nil {
			return "", err
		}
		if err := tmpFile.Close(); err != nil {
			return "", err
		}

		width, height := pdf.GetPageSize()
		posX := (width - pdfW) / 2
		posY := (height - pdfH) / 2

		pdf.AddPage()
		pdf.Image(tmpFile.Name(), posX, posY, pdfW, pdfH, false, "", 0, "")
	}

	var pdfBuf bytes.Buffer
	err := pdf.Output(&pdfBuf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(pdfBuf.Bytes()), nil
}

func saveToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func compressPdf(encodedData string) string {
	base64Images, err := pdfToImages(encodedData)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	var compressedImages []string

	for _, base64Image := range base64Images {

		if err != nil {
			fmt.Println(err)
			return ""
		}

		compressedData, err := compressImage(base64Image, "jpg")

		if err != nil {
			fmt.Println(err)
			return ""
		}

		compressedImages = append(compressedImages, base64.StdEncoding.EncodeToString(compressedData))
	}

	compressedPdf, err := imagesToPdf(compressedImages)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return compressedPdf
}

func main() {
	format := "pdf"
	encodedData := "pdf_base64_data"

	filename := fmt.Sprintf("compressed_file.%s", format)
	downloadsFolder := filepath.Join(os.Getenv("HOME"), "Downloads")
	outputPath := filepath.Join(downloadsFolder, filename)

	if format == "pdf" {
		encodedData = compressPdf(encodedData)
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			fmt.Println("Error decoding base64 data:", err)
			return
		}
		err = saveToFile(outputPath, decodedData)
		if err != nil {
			fmt.Println("Error saving file:", err)
		} else {
			fmt.Println("File saved to:", outputPath)
		}
	} else {
		compressedData, _ := compressImage(encodedData, format)
		err := saveToFile(outputPath, compressedData)
		if err != nil {
			fmt.Println("Error saving file:", err)
		} else {
			fmt.Println("File saved to:", outputPath)
		}
	}
}
