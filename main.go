package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shionfujie/fileconv/conv"

	"github.com/shionfujie/fileconv/io/ioutil"
)

func main() {
	fileNames := os.Args[1:]
	outputs := make(chan string, len(fileNames))

	spinner := ioutil.NewSpinner(100 * time.Millisecond)
	defer spinner.Stop()
	for _, file := range fileNames {
		go func(file string) {
			c := conv.NewConverter(file)
			outputs <- formatMassage(file, c.SourceType, c.TargetType, c.Convert())
		}(file)
	}

	length := len(fileNames)
	width := length/10 + 1
	for i := 1; i <= length; i++ {
		spinner.Printf("[%*d/%d] %s\n", width, i, length, <-outputs)
	}
}

func formatMassage(file string, sourceType string, targetType string, err error) string {
	if err != nil {
		return fmt.Sprintf("Fail\tconvert %s: %v", file, err)
	}
	return fmt.Sprintf("ok  \tconvert %s: %s file converted successfully to %s file", file, sourceType, targetType)
}
