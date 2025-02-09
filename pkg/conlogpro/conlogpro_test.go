package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)


var (
    testFile string
    keywords []string
)

func init() {
    // Define flags
    flag.StringVar(&testFile, "file", "", "path to log file")
}


// this is a
func verifyCount(filepath string, keyword string) int {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
			count++
		}
	}

	return count
}

func TestMain(t *testing.T) {

	// Parse test flags first
    flag.Parse()
    
    // Get the arguments after -args
    testArgs := flag.Args()
    
    // Find the -file argument
    var filepath string = testFile
    // for i := 0; i < len(testArgs); i++ {
    //     if testArgs[i] == "-file" && i+1 < len(testArgs) {
    //         filepath = testArgs[i+1]
	// 		fmt.Print("****", filepath)
    //         // Remove the -file and its value from the args slice
    //         testArgs = append(testArgs[:i], testArgs[i+2:]...)
    //         break
    //     }
    // }
    
    if filepath == "" {
        t.Fatalf("Error: -file flag is required")
    }

	// Remaining args are keywords
    keywords = testArgs

	if len(keywords) == 0 {
        t.Fatalf("Error: at least one keyword argument is required")
    }

	var keywordMap map[string]int = make(map[string]int)

	var oput KV
	oput.output = make(map[string]int)





	batchChan := make(chan []string)

	ProcessLogFile(filepath, keywords, batchChan, keywordMap)

	done := make(chan int)
	defer close(done)

	CPUCount := 12 // or runtime.NumCPU()

	getLnChans := make([]chan map[string]int, CPUCount)

	//fan out
	for i := 0; i < CPUCount; i++ {
		getLnChans[i] = getLines(batchChan, done, CountKeywords, keywordMap)
	}

	//fan in
	fstream := fanIn(done, getLnChans...)

	for kv := range fstream {
		oput.mu.Lock()
		for w, c := range kv {
			oput.output[w] += c
		}
		oput.mu.Unlock()
	}

	actualCount := verifyCount(filepath, "info")
	fmt.Printf("\nVerification count for 'info': %d\n", actualCount)
	fmt.Printf("Our count: %d\n", oput.output["info"])
	if actualCount != oput.output["info"] {
		t.Errorf("Result was incorrect, got: %x, want: %x.", oput.output["info"], actualCount)
	}
}
