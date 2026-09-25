package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	outputFmt := flag.String("output", "text", "Output format (text, json)")
	flag.Parse()

	command := os.Args[1]
	subcommand := ""
	if len(os.Args) > 2 {
		subcommand = os.Args[2]
	}

	switch command {
	case "network":
		if subcommand == "status" {
			fmt.Printf("Network status: OK (Format: %s)\n", *outputFmt)
		} else {
			fmt.Println("Unknown network subcommand")
		}
	case "participant":
		if subcommand == "list" {
			fmt.Printf("Participant list (Format: %s)\n", *outputFmt)
		} else {
			fmt.Println("Unknown participant subcommand")
		}
	case "settlement":
		switch subcommand {
		case "create", "get", "accept", "reject", "reserve", "execute", "proof":
			fmt.Printf("Settlement %s (Format: %s)\n", subcommand, *outputFmt)
		default:
			fmt.Println("Unknown settlement subcommand")
		}
	case "demo":
		switch subcommand {
		case "bilateral":
			fmt.Println("Running Bilateral Demo...")
		case "three-party":
			fmt.Println("Running Three-Party Demo...")
		default:
			fmt.Println("Unknown demo subcommand")
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: afrorailctl [--output json] <command> <subcommand>")
	fmt.Println("\nCommands:")
	fmt.Println("  network status")
	fmt.Println("  participant list")
	fmt.Println("  settlement [create|get|accept|reject|reserve|execute|proof]")
	fmt.Println("  demo [bilateral|three-party]")
}
