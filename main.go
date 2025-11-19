package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define command-line flags
	name := flag.String("name", "World", "name to greet")
	verbose := flag.Bool("verbose", false, "enable verbose output")

	// Parse the flags
	flag.Parse()

	// Greet the user
	if *verbose {
		fmt.Printf("Running EEC289Q CLI application\n")
		fmt.Printf("Greeting: %s\n", *name)
	}

	fmt.Printf("Hello, %s!\n", *name)

	// Check for any additional positional arguments
	args := flag.Args()
	if len(args) > 0 {
		fmt.Printf("Additional arguments: %v\n", args)
	}

	os.Exit(0)
}
