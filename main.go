package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/adrium/goheif"
)

// TODO: add cli args?
// TODO: optimize with parallelism

// Skip Writer for exif writing
type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	// CONSTANTS
	INPUT_PATH := "./input"
	OUTPUT_PATH := "./output"
	// FILE_TYPES_INPUT := (".heic", ".png")
	FILE_TYPE_OUTPUT := "jpg"

	defer timer("main")()

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
			MoveFile(inputFileName, outputFileName)
		} else {
			outputFileName = strings.ReplaceAll(inputFileName, ".HEIC", fmt.Sprintf(".%s", FILE_TYPE_OUTPUT))
			conversion_err := convertHeicToJpg(inputFileName, outputFileName)

			if conversion_err != nil {
				log.Fatal(conversion_err)
			}
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
		return err
	}
	defer fileInput.Close()

	// Extract exif to add back in after conversion
	exif, err := goheif.ExtractExif(fileInput)
	if err != nil {
		return err
	}

	img, err := goheif.Decode(fileInput)
	if err != nil {
		return err
	}

	fileOutput, err := os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fileOutput.Close()

	// Write both convert file + exif data back
	w, _ := newWriterExif(fileOutput, exif)
	err = jpeg.Encode(w, img, nil)
	if err != nil {
		return err
	}

	return nil
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}
	return nil
}
