package main

import (
	"fmt"

	sha256 "github.com/minio/sha256-simd"
)

func main() {
	server := sha256.NewAvx512Server()
	h512 := sha256.NewAvx512(server)
	h512.Write([]byte("hello world"))
	digest := h512.Sum([]byte{})
	fmt.Println(digest)
}
