package fileOperations

import (
	"convert-heic/helpers/timer"
	"fmt"

	// "convert-heic/helpers/timer"
	"io"
	"os"
	"sync"
)

func MoveFile(wg *sync.WaitGroup, sourcePath string, destPath string) error {
	defer timer.FuncTimer("MoveFile")()

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
