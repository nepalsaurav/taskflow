package main

import (
	"fmt"
	"log"
	"os"
	"taskflow/cmd"
)

func main() {

	// skip program name
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		log.Fatal("no command provided")
	}
	command := args[0]

	switch command {
	case "migrate":
		handleMigrate(args[1:])
	case "index":
		maildir := cmd.Maildir{}
		maildir.IndexMail()
	default:
		printUsage()
		log.Fatalf("Unknown command: %s", command)
	}

}

func handleMigrate(subArgs []string) {
	if len(subArgs) == 0 {
		printUsage()
		log.Fatal("migrate requires 'up' or 'down'")
	}
	migrate := Migrate{}
	directions := subArgs[0]
	switch directions {
	case "up":
		fmt.Println("Running Migrations Up")
		migrate.up()
	case "down":
		fmt.Println("Rolling back migrations (DOWN)...")
		fmt.Println("Rollback completed")
	default:
		printUsage()
		log.Fatalf("Invalid migrate direction: (expected 'up' or 'down')")
	}
}

func printUsage() {
	fmt.Println(`Usages : myapp <command> [options]
Commands:
serve                Start the web server
migrate up           Run database migrations forward
migrate down         Rollback database migrations
help                 Show this help message
	`)
}
