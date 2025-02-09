package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	// "strings"
	"sync"
)

type Part struct {
	offset, size int64
}

func SplitFile(inputPath string, numParts int) ([]Part, error) {
	const maxLineLength = 4096 //initially 100

	f, err := os.Open(inputPath)
	
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()

	if size == 0 {
		log.Fatalf("file empty")
	}
	splitSize := size / int64(numParts)

	buf := make([]byte, maxLineLength)

	parts := make([]Part, 0, numParts)
	offset := int64(0)
	for offset < size {
		seekOffset := max(offset+splitSize-maxLineLength, 0)
		if seekOffset > size {
			break
		}

		_, err := f.Seek(seekOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, _ := io.ReadFull(f, buf)
		chunk := buf[:n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			return nil, fmt.Errorf("newline not found at offset %d", offset+splitSize-maxLineLength)
		}
		remaining := len(chunk) - newline - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, Part{offset, nextOffset - offset})

		offset = nextOffset
	}
	return parts, nil
}

func ProcessPart(filepath string, offset int64, size int64, wg *sync.WaitGroup, results chan []string) {
	defer wg.Done()

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}
	f := io.LimitedReader{R: file, N: size}

	scanner := bufio.NewScanner(&f)
	line_list := []string{}

	lineCount := 0
	// wordCount := make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		line_list = append(line_list, line)
		// results <- line_list

	}

	results <- line_list

}

func ProcessLogFile(filepath string, klist []string, batchChan chan []string, km map[string]int) {
	var wg sync.WaitGroup
	// batchChan:= make(chan []string)
	// done := make(chan int)

	for _, kw := range klist {
		km[kw] = 0
	}

	parts, err := SplitFile(filepath, 10)

	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("file does not exist")
		}
		panic(err)
	}

	for _, part := range parts {
		wg.Add(1)
		go ProcessPart(filepath, part.offset, part.size, &wg, batchChan)
	}

	go func() {
		wg.Wait()
		close(batchChan)
	}()

}
