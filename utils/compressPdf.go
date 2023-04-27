package utils

import (
	"encoding/base64"
	"fmt"
)

func CompressPdf(encodedData string) string {
	base64Images, err := PdfToImages(encodedData)

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

		compressedData, err := CompressImage(base64Image)

		if err != nil {
			fmt.Println(err)
			return ""
		}

		compressedImages = append(compressedImages, base64.StdEncoding.EncodeToString(compressedData))
	}

	compressedPdf, err := ImagesToPdf(compressedImages)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return compressedPdf
}
