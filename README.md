# Image Conversion with Go

Convert all images in a specified directory to a specified file type.

- To specify the directory that has all the images you want to convert, change the `INPUT_PATH` variable in `./helpers/constants/constants.go`. This can be anything, as long as it is a valid path 😄.
- To specify the directory that will contain all the **converted** images, change the `OUTPUT_PATH` variables in `./helpers/constants/constants.go`. Again this can be anything, as long as it is a valid path 😄.
- To specify the output type of all the **converted** images, change the `OUTPUT_FILE_TYPE` variable in `./helpers/constants/constants.go`.

## Installation

- Specify the necessary parameters in `./helpers/constants/constants.go` as detailed above.
- Run the command below. This will start your image conversion!
  - `go run main.go`

## Notes

- This script uses Goroutines and Waitgroups to see a bit of a performance boost, rather than employing synchronous code.
- If you have any suggestions on how to make this script better, please add an issue here:
  - https://github.com/espitiaandres/go-image-conversion/issues
