package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	// "time"
	"convert-heic/helpers/fileOperations"
	"convert-heic/helpers/timer"
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

			go fileOperations.ConvertHeicToJpg(inputFileName, outputFileName)

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
