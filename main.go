package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	sha256simd "github.com/minio/sha256-simd"

	sha256native "crypto/sha256"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please provide an argument")
		return
	}
	iterations, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Failed to convert string to int:", err)
		return
	}

	start := time.Now()
	for i := 0; i < iterations; i++ {
		sha256native.Sum256([]byte("hello world"))
	}
	elapsed := time.Since(start)
	fmt.Printf("Native Elapsed time: %s\n", elapsed)

	server := sha256simd.NewAvx512Server()
	h512 := sha256simd.NewAvx512(server)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		h512.Write([]byte("hello world"))
		h512.Sum([]byte{})
	}
	elapsed = time.Since(start)
	fmt.Printf("AVX512 Elapsed time: %s\n", elapsed)

	// concurrent

	start = time.Now()
	// Set up a pool of 10 goroutines
	poolSize := 10
	jobs := make(chan []byte, poolSize)
	results := make(chan [32]byte, poolSize)

	var wg sync.WaitGroup
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range jobs {
				hash := sha256native.Sum256([]byte(msg))
				results <- hash
			}
		}()
	}

	// Enqueue the messages to hash
	//messages := []string{"hello", "world", "foo", "bar", "baz", "qux", "foofoo", "barbar", "bazbaz", "quxqux"}
	go func() {

		for i := 0; i < iterations; i++ {
			// sha256native.Sum256([]byte("hello world"))
			jobs <- []byte("hello world")
		}
		close(jobs)
	}()

	// Collect the results
	var hashes [][32]byte
	for i := 0; i < iterations; i++ {
		hashes = append(hashes, <-results)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	elapsed = time.Since(start)
	fmt.Printf("Native Elapsed time: %s\n", elapsed)

}
