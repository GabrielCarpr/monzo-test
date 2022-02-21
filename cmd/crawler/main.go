package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Seed URL not provided")
		log.Fatal()
	}

	target := os.Args[1]
	log.Printf("%v", target)
}
