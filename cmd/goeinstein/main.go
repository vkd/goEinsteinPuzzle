package main

import (
	"log"

	"github.com/vkd/goeinstein"
)

func main() {
	log.Printf("Starting...")
	err := goeinstein.Main()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
