package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"convert-heic/helpers/constants"
	"convert-heic/helpers/fileOperations"
	"convert-heic/helpers/timer"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	defer timer.FuncTimer("main")()

	if _, existErr := os.Stat(constants.OUTPUT_PATH); os.IsNotExist(existErr) {
		if mkdirErr := os.Mkdir(constants.OUTPUT_PATH, os.ModePerm); mkdirErr != nil {
			log.Fatal(mkdirErr)
		}
	}

	allEntries, err := os.ReadDir(constants.INPUT_PATH)

	var imageEntries []fs.DirEntry
	validImageExtensions := []string{".jpg", ".png", ".heic", ".tif", ".eps"}

	for _, e := range allEntries {
		fileExtension := strings.ToLower(filepath.Ext(e.Name()))

		if stringInSlice(fileExtension, validImageExtensions) {

			fmt.Println("helppp")
			fmt.Println(e.Name())
			imageEntries = append(imageEntries, e)
		}
	}

	// Filter out files that aren't images

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, e := range imageEntries {
		inputFileName := fmt.Sprintf("%s/%s", constants.INPUT_PATH, e.Name())
		outputFileName := fmt.Sprintf("%s/%s", constants.OUTPUT_PATH, e.Name())

		heicFile := strings.Contains(
			strings.ToLower(inputFileName),
			strings.ToLower(".heic"),
		)

		if !heicFile {
			outputFileName = strings.ReplaceAll(inputFileName, constants.INPUT_PATH, constants.OUTPUT_PATH)

			wg.Add(1)

			go func() {
				defer wg.Done()
				fileOperations.MoveFile(inputFileName, outputFileName)
			}()
		} else {
			outputFileName = strings.ReplaceAll(outputFileName, ".HEIC", fmt.Sprintf(".%s", constants.FILE_TYPE_OUTPUT))

			wg.Add(1)

			go func() {
				defer wg.Done()
				fileOperations.ConvertHeicToJpg(inputFileName, outputFileName)
			}()
		}
	}

	log.Println("Waiting for goroutines to complete...")
	wg.Wait()
	log.Printf("Completed image conversion for %v files", len(imageEntries))
}
