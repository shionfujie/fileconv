package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/shionfujie/fileconv/io/ioutil"
)

func main() {
	fileNames := os.Args[1:]
	outputs := make(chan string, len(fileNames))

	spinner := ioutil.NewSpinner(100 * time.Millisecond)
	defer spinner.Stop()
	for _, file := range fileNames {
		go func(file string) {
			c := NewConverter(file)
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

type Converter struct {
	SourceType string
	TargetType string
	file       string
	cmd        *exec.Cmd
	pathErr    *os.PathError
}

type cmd struct {
	name       string
	targetType string
}

var (
	djvu2pdf = cmd{"djvu2pdf", "PDF"}
)

var cmds = map[string]cmd{
	"DJVU": djvu2pdf,
}

func NewConverter(file string) (c *Converter) {
	c = &Converter{
		file: file,
	}
	c.SourceType = sourceType(file)
	if _, err := os.Stat(file); err != nil {
		c.pathErr = err.(*os.PathError)
		return
	}
	cmd := cmds[c.SourceType]
	c.TargetType = cmd.targetType
	c.cmd = exec.Command(cmd.name, file)
	return
}

func sourceType(file string) string {
	return strings.ToUpper(filepath.Ext(file)[1:])
}

func (c *Converter) Convert() error {
	if c.pathErr != nil {
		return c.pathErr.Err
	}
	if _, err := c.cmd.Output(); err != nil {
		return fmt.Errorf("failed to convert %s file to %s file: %v", c.SourceType, c.TargetType, err)
	}
	return nil
}

func convert(file string) error {
	if _, err := os.Stat(file); err != nil {
		return err.(*os.PathError).Err
	}
	if _, err := exec.Command("djvu2pdf", file).Output(); err != nil {
		return fmt.Errorf("failed to convert DJVU file to PDF file: %v", err)
	}
	return nil
}
