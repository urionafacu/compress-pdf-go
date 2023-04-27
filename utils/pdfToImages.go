package utils

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/gen2brain/go-fitz"
)

func PdfToImages(pdfBase64 string) ([]string, error) {
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
