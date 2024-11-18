package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dibyanshu47/building-git/pkg/gitobject"
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
	case "cat-file":
		cmdCatFile(os.Args[2:])
	case "hash-object":
		cmdHashObject(os.Args[2:])
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

func cmdCatFile(args []string) {
	catFileCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)
	err := catFileCmd.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing cat-file command: %v\n", err)
		os.Exit(1)
	}

	objectType := catFileCmd.Args()[0]
	object := catFileCmd.Args()[1]

	repo, err := gitrepository.RepoFind(".", true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding repository: %v\n", err)
		os.Exit(1)
	}

	obj := gitobject.ObjectFind(repo, object, objectType, true)

	gitObject, err := gitobject.ObjectRead(repo, obj)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading object: %v\n", err)
		os.Exit(1)
	}

	serializedObject := gitObject.Serialize()
	os.Stdout.Write(serializedObject)
}

func cmdHashObject(args []string) {
	hashObjectCmd := flag.NewFlagSet("hash-object", flag.ExitOnError)
	write := hashObjectCmd.Bool("write", false, "Write object to the object database")
	objectType := hashObjectCmd.String("type", "", "Type of the object")
	err := hashObjectCmd.Parse(args)
	if err != nil {
		fmt.Printf("Error parsing hash-object command: %v\n", err)
		hashObjectCmd.Usage()
		os.Exit(1)
	}

	var repo *gitrepository.GitRepository
	if *write {
		repo, err = gitrepository.RepoFind(".", true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding repository: %v\n", err)
			os.Exit(1)
		}
	} else {
		repo = nil
	}

	path := hashObjectCmd.Args()[0]
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	sha, err := gitobject.ObjectHash(file, *objectType, repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error hashing object: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(sha)
}
