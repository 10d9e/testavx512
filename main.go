package main

import (
	"fmt"
	"os"
	"strconv"
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
	fmt.Printf("Normal Elapsed time: %s\n", elapsed)

	server := sha256simd.NewAvx512Server()
	h512 := sha256simd.NewAvx512(server)

	start = time.Now()
	for i := 0; i < iterations; i++ {
		h512.Write([]byte("hello world"))
		h512.Sum([]byte{})
	}
	elapsed = time.Since(start)
	fmt.Printf("AVX512 Elapsed time: %s\n", elapsed)

}
