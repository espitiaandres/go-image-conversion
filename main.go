package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"strings"

	"github.com/adrium/goheif"
)

// Skip Writer for exif writing
type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func main() {
	if err := os.Mkdir("output", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	entries, err := os.ReadDir("input")
	if err != nil {
		fmt.Println("the strokes")
		// log.Fatal(err)
	}

	fmt.Println(entries)

	for _, e := range entries {
		fmt.Println(e.Name())

		input_file_name := fmt.Sprintf("input/%s", e.Name())
		output_file_name := fmt.Sprintf("output/%s", e.Name())

		heicFile := strings.Contains(
			strings.ToLower(input_file_name),
			strings.ToLower(".heic"),
		)

		fmt.Println(heicFile)

		if !heicFile {
			continue
		}

		conversion_err := convertHeicToJpg(input_file_name, strings.ReplaceAll(output_file_name, ".HEIC", ".jpg"))

		if conversion_err != nil {
			log.Fatal(conversion_err)
		}
	}

	log.Println("Conversion Passed")
}

// convertHeicToJpg takes in an input file (of heic format) and converts
// it to a jpeg format, named as the output parameters.
func convertHeicToJpg(input, output string) error {

	fmt.Println("input")
	fmt.Println(input)
	fmt.Println("ouput")
	fmt.Println(output)

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
