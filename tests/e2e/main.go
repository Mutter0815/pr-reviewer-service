package main

import (
	"fmt"
	"os"
)

func main() {
	client := newHTTPClient()

	if err := runScenario(client); err != nil {
		fmt.Fprintf(os.Stderr, "e2e failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("e2e scenario passed")
}
