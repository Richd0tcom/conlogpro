package main

import (

	"fmt"
	"os"
	"sort"

	// "runtime"
	"strings"
	"sync"
	// "sync"
)

type KV struct {
	output map[string]int
	mu     sync.Mutex
}

func getLines(batch chan []string, done <-chan int, fn func(string, map[string]int) map[string]int, km map[string]int) chan map[string]int {
	stream := make(chan map[string]int)

	go func() {
		defer close(stream)
		for lines := range batch {
			for _, line := range lines {
				select {
				case <-done:
					return

				case stream <- fn(line, km):
				}
			}

		}
	}()

	return stream
}

// count each unique occurence of a keyword and returns an associative map
// this function returns early because we assume keywords will only appear once per line
// due to the format `2023-10-28 12:00:01 - INFO - User logged in`
func CountKeywords(line string, km map[string]int) map[string]int {
	m := make(map[string]int)

	words := strings.Fields(strings.ToLower(line))

	for _, word := range words {

		_, ok := km[word]
		if ok {
			m[word]++
			continue //finish the line early
		}

	}

	return m
}

//collects all the results of the different fanned out channel into one channel
func fanIn(done <-chan int, channels ...chan map[string]int) chan map[string]int {
	var wg sync.WaitGroup
	fannedInStream := make(chan map[string]int)

	transfer := func(c chan map[string]int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return

			case fannedInStream <- i:
			}
		}
	}

	for _, ch := range channels {
		wg.Add(1)
		go transfer(ch)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}





//prints key word frequency in decending order
func printMapDescending(m map[string]int) {
    // Create a slice of key-value pairs
    type keyval struct {
        Key   string
        Value int
    }
    pairs := make([]keyval, 0, len(m))
    
    // Convert map to slice of kv pairs
    for k, v := range m {
        pairs = append(pairs, keyval{k, v})
    }
    
    // Sort slice by value in descending order
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].Value > pairs[j].Value
    })
    

    // Print with (some) alignment
    fmt.Println("\nKeyword Frequency:")
    fmt.Println(strings.Repeat("-",  15))
    for _, pair := range pairs {
        fmt.Printf("%s : %d\n", pair.Key, pair.Value)
    }
    fmt.Println(strings.Repeat("-", 15))
}

func main() {

	var keywordMap map[string]int = make(map[string]int)

	var oput KV
	oput.output = make(map[string]int)

	argsWithoutProg := os.Args[1:]
	pathToFile := argsWithoutProg[0]

	batchChan := make(chan []string)

	ProcessLogFile(pathToFile, argsWithoutProg[1:], batchChan, keywordMap)

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

	// actualCount := verifyCount(pathToFile, "info")
	// fmt.Printf("\nVerification count for 'info': %d\n", actualCount)
	fmt.Printf("Our count: %d\n", oput.output["info"])

	printMapDescending(oput.output)

}
