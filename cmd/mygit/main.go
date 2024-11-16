package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dibyanshu47/building-git/pkg/gitrepository"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mygit <command> [<args>]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cmdInit(os.Args[2:])
	// ... other cases for other commands
	default:
		fmt.Printf("%q is not a valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

func cmdInit(args []string) {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)
	path := initCmd.String("path", ".", "Path where to create the repository")
	err := initCmd.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing init command: %v\n", err)
		initCmd.Usage()
		os.Exit(1)
	}

	repo, err := gitrepository.RepoCreate(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing repository: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Initialized empty Git repository in %s\n", repo.Gitdir)
}
