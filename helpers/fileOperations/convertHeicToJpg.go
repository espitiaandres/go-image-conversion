package fileOperations

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/adrium/goheif"
)

// convertHeicToJpg takes in an input file (of heic format) and converts
// it to a jpeg format, named as the output parameters.
func ConvertHeicToJpg(input string, output string) error {

	log.Println("1 -----")

	fileInput, err := os.Open(input)
	if err != nil {
		log.Println("os.Open() failed")
	}

	log.Println("2 -----")
	defer fileInput.Close()

	// Extract exif to add back in after conversion
	exif, err := goheif.ExtractExif(fileInput)

	log.Println("3 -----")
	if err != nil {
		log.Println("goheif.ExtractExif() failed")
	}

	log.Println("4 -----")
	img, err := goheif.Decode(fileInput)
	if err != nil {
		log.Println("goheif.Decode() failed")
	}

	log.Println("5 -----")
	fileOutput, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("os.OpenFile() failed")
	}

	log.Println("6 -----")
	defer fileOutput.Close()

	// Write both convert file + exif data back
	w, _ := newWriterExif(fileOutput, exif)

	log.Println("7 -----")
	err = jpeg.Encode(w, img, nil)

	log.Println("8 -----")
	if err != nil {
		log.Println("jpeg.Encode() failed")
	}

	log.Println("9 -----")

	return nil
}
