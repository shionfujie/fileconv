package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	file := os.Args[1]
	go spinner(100 * time.Millisecond)
	_, err := exec.Command("djvu2pdf", file).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\033[Fconvert %s: DJVU file converted successfully to PDF file\n", file)
}

func spinner(delay time.Duration) {
	fmt.Println()
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\033[F%c\n", r)
			time.Sleep(delay)
		}
	}
}
