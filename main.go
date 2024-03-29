package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"regexp"

	"image-conversion/helpers/constants"
	"image-conversion/helpers/fileOperations"
	"image-conversion/helpers/helpers"
	"image-conversion/helpers/timer"
)

// Only for converting .heic files into .jpg

func main() {
	defer timer.FuncTimer("main")()

	if _, existErr := os.Stat(constants.OUTPUT_PATH); os.IsNotExist(existErr) {
		if mkdirErr := os.Mkdir(constants.OUTPUT_PATH, os.ModePerm); mkdirErr != nil {
			log.Fatal(mkdirErr)
		}
	}

	allEntries, err := os.ReadDir(constants.INPUT_PATH)

	var imageEntries []fs.DirEntry
	validImageExtensions := []string{".jpg", ".jpeg", ".png", ".heic", ".tiff", ".eps", ".bmp"}

	for _, e := range allEntries {
		fileExtension := strings.ToLower(filepath.Ext(e.Name()))

		// Filter out files that aren't images
		if helpers.StringInSlice(fileExtension, validImageExtensions) {
			imageEntries = append(imageEntries, e)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, e := range imageEntries {
		inputFileName := fmt.Sprintf("%s/%s", constants.INPUT_PATH, e.Name())
		outputFileName := fmt.Sprintf("%s/%s", constants.OUTPUT_PATH, e.Name())

		fileWithDesiredOutputType := strings.Contains(
			strings.ToLower(inputFileName),
			strings.ToLower(fmt.Sprintf(".%s", constants.OUTPUT_FILE_TYPE)),
		)

		if fileWithDesiredOutputType {
			outputFileName = strings.ReplaceAll(inputFileName, constants.INPUT_PATH, constants.OUTPUT_PATH)

			wg.Add(1)

			go func() {
				defer wg.Done()
				fileOperations.MoveFile(inputFileName, outputFileName)
			}()
		} else {
			re := regexp.MustCompile(`(?i)heic`)
			outputFileName = re.ReplaceAllString(outputFileName, constants.OUTPUT_FILE_TYPE)
			
			wg.Add(1)

			go func() {
				defer wg.Done()
				fileOperations.ImageConvertToTarget(inputFileName, outputFileName)
			}()
		}
	}

	log.Println("Waiting for goroutines to complete...")
	wg.Wait()
	log.Printf("Completed image conversion for %v files", len(imageEntries))
}
