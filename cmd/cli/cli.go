package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define flags for customizing output
	ping := flag.Bool("ping", true, "Ping the server")

	// Parse command line flags
	flag.Parse()

	// Get remaining arguments after flags
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: heyecho [flags] [arguments to echo]")
		fmt.Println("Flags:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Respond to the command line arguments
	if *ping {
		fmt.Println("Pinging server...")
	}
}
