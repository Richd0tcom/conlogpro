package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// isPalindrome checks if a number is palindromic
func isPalindrome(num int) bool {
	if num < 0 {
		return false
	}

	original := num
	reversed := 0

	for num > 0 {
		digit := num % 10
		reversed = reversed*10 + digit
		num /= 10
	}

	return original == reversed
}

// isPrime checks if a number is prime
func isPrime(num int) bool {
	if num < 2 || num%2 == 0 {
		return false
	}

	if num == 2 {
		return true
	}

	sqrt := int(math.Sqrt(float64(num)))
	for i := 3; i <= sqrt; i += 2 {
		if num%i == 0 {
			return false
		}
	}
	return true
}

// worker processes numbers from the input channel and sends prime palindromes to the result channel
func worker(input <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for num := range input {
		if isPalindrome(num) && isPrime(num) {
			results <- num
		}
	}
}

// findNPrimePalindromes finds the first N prime palindromic numbers
func findNPrimePalindromes(N int, numWorkers int) ([]int, int) {
	numbers := make(chan int)    
	priPalChan := make(chan int) 
	done := make(chan int)  

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(numbers, priPalChan, &wg)
	}
	defer close(numbers)
	// Start number generator
	go func() {
		num := 2
		for {
			select {
			case <-done:

				return
			default:
				numbers <- num
				num++
			}
		}
	}()

	// Collect results
	found := make([]int, 0, N)
	sum := 0

	// Start result collector
	go func() {
		wg.Wait()
		close(priPalChan)
	}()

	// Collect N results
	for num := range priPalChan {
		found = append(found, num)
		sum += num
		if len(found) == N {
			close(done)
			break
		}
	}

	return found, sum
}

func main() {

	argsWithoutProg := os.Args[1:]
	N, err := strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		log.Fatalf("Invalid Argument")
	}
	numWorkers := runtime.NumCPU()

	// fmt.Printf("Finding first %d prime palindromic numbers...\n", N)

	_, sum := findNPrimePalindromes(N, numWorkers)

	// fmt.Println("Prime Palindromic Numbers:", numbers)
	fmt.Println("Sum:", sum)
}
