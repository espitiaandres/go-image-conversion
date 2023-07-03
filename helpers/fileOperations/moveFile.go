package fileOperations

import (
	"convert-heic/helpers/timer"
	"io"
	"log"
	"os"
)

func MoveFile(sourcePath string, destPath string) {
	defer timer.FuncTimer("MoveFile")()

	inputFile, err := os.Open(sourcePath)

	if err != nil {
		log.Printf("os.Open() failed: %s", err)
	}

	outputFile, err := os.Create(destPath)

	if err != nil {
		inputFile.Close()
		log.Printf("os.Create() failed: %s", err)
	}

	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)

	inputFile.Close()

	if err != nil {
		log.Printf("os.Open() failed: %s", err)
	}
}
