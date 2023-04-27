package utils

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func ImagesToPdf(base64Images []string) (string, error) {
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
