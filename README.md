# Image Conversion with Go

Convert .heic images in a specified directory to a jpg. Heic images are usually produced by apple products when using the ultra-wide camera mode. This script makes it easy to convert these images into the jpg format which is universally accepted.

- To specify the directory that has all the images you want to convert, change the `INPUT_PATH` variable in `./helpers/constants/constants.go`. This can be anything, as long as it is a valid path 😄.
- To specify the directory that will contain all the **converted** images, change the `OUTPUT_PATH` variables in `./helpers/constants/constants.go`. Again this can be anything, as long as it is a valid path 😄.
- Currently, only `.jpg` images are the only file type that is supported for the **converted** images.

## Installation

- Specify the necessary parameters in `./helpers/constants/constants.go` as detailed above.
- Run the command below. This will start your image conversion!
  - `go run main.go`

## Notes

- This script uses Goroutines and Waitgroups to see a bit of a performance boost, rather than employing synchronous code. Since this script reads images off a disk, this is I/O bound so implementing Goroutines and Waitgroups help with performance.
- If you have any suggestions on how to make this script better, please add an issue here:
  - https://github.com/espitiaandres/go-image-conversion/issues
