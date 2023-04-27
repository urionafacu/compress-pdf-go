package main

import (
	"compresor-file/utils"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./compress input_filepath [output_filepath]")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	fmt.Println("Input path:", inputPath)
	var outputPath string
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	} else {
		outputPath = inputPath + "_compressed"
	}

	inputBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(inputBytes)
	ext := filepath.Ext(inputPath)
	format := ext[1:] // Remove the leading "."

	fmt.Println("Format", format)

	switch format {
	case "pdf":
		encodedData = utils.CompressPdf(encodedData)
		decodedData, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			fmt.Println("Error decoding base64 data:", err)
			return
		}
		err = utils.SaveToFile(outputPath, decodedData)
	case "jpeg", "jpg", "png":
		compressedData, err := utils.CompressImage(encodedData)
		if err != nil {
			fmt.Println("Error compressing image:", err)
			return
		}
		err = utils.SaveToFile(outputPath, compressedData)
	default:
		fmt.Println("Unsupported file format:", format)
		return
	}
}
