package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strings"

	// "time"
	"convert-heic/helpers/fileOperations"
	"convert-heic/helpers/timer"

	"github.com/adrium/goheif"
)

// CONSTANTS
const INPUT_PATH string = "./input"
const OUTPUT_PATH string = "./output"
const FILE_TYPE_OUTPUT string = "jpg"

func main() {
	defer timer.FuncTimer("main")()

	if _, existErr := os.Stat(OUTPUT_PATH); os.IsNotExist(existErr) {
		if mkdirErr := os.Mkdir(OUTPUT_PATH, os.ModePerm); mkdirErr != nil {
			log.Fatal(mkdirErr)
		}
	}

	entries, err := os.ReadDir(INPUT_PATH)

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		inputFileName := fmt.Sprintf("%s/%s", INPUT_PATH, e.Name())
		outputFileName := ""

		heicFile := strings.Contains(
			strings.ToLower(inputFileName),
			strings.ToLower(".heic"),
		)

		if !heicFile {
			outputFileName = strings.ReplaceAll(inputFileName, INPUT_PATH, OUTPUT_PATH)
			go fileOperations.MoveFile(inputFileName, outputFileName)
		} else {
			outputFileName = strings.ReplaceAll(inputFileName, ".HEIC", fmt.Sprintf(".%s", FILE_TYPE_OUTPUT))
			// conversion_err := convertHeicToJpg(inputFileName, outputFileName)

			go convertHeicToJpg(inputFileName, outputFileName)

			// if conversion_err != nil {
			// 	log.Fatal(conversion_err)
			// }
		}

		fmt.Println("--------------")
		fmt.Println(inputFileName)
		fmt.Println(outputFileName)
	}

	log.Println("Conversion Passed")
}

// convertHeicToJpg takes in an input file (of heic format) and converts
// it to a jpeg format, named as the output parameters.
func convertHeicToJpg(input, output string) error {
	fileInput, err := os.Open(input)
	if err != nil {
		log.Println("os.Open() failed")
	}
	defer fileInput.Close()

	// Extract exif to add back in after conversion
	exif, err := goheif.ExtractExif(fileInput)
	if err != nil {
		log.Println("goheif.ExtractExif() failed")
	}

	img, err := goheif.Decode(fileInput)
	if err != nil {
		log.Println("goheif.Decode() failed")
	}

	fileOutput, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("os.OpenFile() failed")
	}
	defer fileOutput.Close()

	// Write both convert file + exif data back
	w, _ := fileOperations.NewWriterExif(fileOutput, exif)
	err = jpeg.Encode(w, img, nil)
	if err != nil {
		log.Println("jpeg.Encode() failed")
	}

	return nil
}
