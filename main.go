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
			// outputFileName = strings.ReplaceAll(inputFileName, constants.INPUT_PATH, constants.OUTPUT_PATH)
			// go fileOperations.MoveFile(inputFileName, outputFileName)
		} else {
			outputFileName = strings.ReplaceAll(outputFileName, ".HEIC", fmt.Sprintf(".%s", constants.FILE_TYPE_OUTPUT))
			// go fileOperations.ConvertHeicToJpg(wg, inputFileName, outputFileName)

			wg.Add(1)

			go func(input string, output string) {
				defer wg.Done()
				fileOperations.ConvertHeicToJpg(inputFileName, outputFileName)
			}(inputFileName, outputFileName)

			// time.Sleep(1 * time.Second)
		}

		// log.Println("--------------")

		// log.Println(inputFileName)
		// log.Println(outputFileName)
	}

	log.Println("Conversion Passed")
	wg.Wait()
}
