package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"convert-heic/helpers/constants"
	"convert-heic/helpers/fileOperations"
	"convert-heic/helpers/timer"
)

func main() {
	defer timer.FuncTimer("main")()

	if _, existErr := os.Stat(constants.OUTPUT_PATH); os.IsNotExist(existErr) {
		if mkdirErr := os.Mkdir(constants.OUTPUT_PATH, os.ModePerm); mkdirErr != nil {
			log.Fatal(mkdirErr)
		}
	}

	entries, err := os.ReadDir(constants.INPUT_PATH)

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, e := range entries {
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
	log.Printf("Completed image conversion for %v files", len(entries))
}
