package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		panic(fmt.Errorf("failed to read the message.txt file, err: %s", err.Error()))
	}

	defer file.Close()
	for line := range getLinesChannel(file) {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		buf := make([]byte, 8)
		var currentLine string
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(fmt.Errorf("failed to read the message.txt file, err: %s", err.Error()))
			}

			currentLine += string(buf[:n])
			var parts []string
			if strings.Contains(currentLine, "\n") {
				parts = strings.Split(currentLine, "\n")
				lastLineIdx := len(parts) - 1
				for i := 0; i < lastLineIdx; i++ {
					lines <- parts[i]
				}

				currentLine = parts[lastLineIdx]
			}
		}

		if len(currentLine) > 0 {
			lines <- currentLine
		}
	}()
	return lines
}
