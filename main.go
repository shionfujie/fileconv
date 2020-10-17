package main

import (
	"fmt"
	"log"
	"os/exec"
	"os"
)

func main() {
	 file := os.Args[1]
	 _, err := exec.Command("djvu2pdf", file).Output()
	 if err != nil {
		 log.Fatal(err)
	 }
	 fmt.Printf("\r%s is successfully converted to pdf", file)
}