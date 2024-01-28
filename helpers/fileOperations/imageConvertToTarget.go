package fileOperations

import (
	"image/jpeg"
	"log"
	"os"

	"github.com/adrium/goheif"
)

// ImageConvertToTarget takes in an input file and converts
// it to a jpeg format, named as the output parameters.
func ImageConvertToTarget(input string, output string) {

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

	// Decode the image frjom the file path into an image.Image type
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
	w, _ := newWriterExif(fileOutput, exif)

	// TODO: have encoders for all file types - need to research this
	// TODO: imports, write switch/case statement for all different file types
	// https://github.com/dawnlabs/photosorcery/blob/master/convert.go

	err = jpeg.Encode(w, img, nil)

	if err != nil {
		log.Println("jpeg.Encode() failed")
	}

	log.Printf("Successfully converted %s", output)
}
