package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/raulviigipuu/gitkeep/internal/gitkeep"
	"github.com/raulviigipuu/gitkeep/internal/gitutils"
)

func main() {
	flag.Parse()
	args := flag.Args()

	// Default path is the current directory if no argument is provided
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Validate if Git is installed
	if !gitutils.IsGitInstalled() {
		fmt.Println("😕 Git is not installed. Please install Git and try again.")
		os.Exit(1)
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("🚫 Error resolving absolute path: %s\n", err)
		os.Exit(1)
	}

	// Check if the path is within a Git repository
	isInGitRepo, repoPath := gitutils.CheckIfGitRepo(absPath)
	if !isInGitRepo {
		fmt.Printf("🔍 The path %s is not within a Git repository.\n", absPath)
		os.Exit(1)
	}

	fmt.Printf("🚀 Managing .gitkeep files in repository: %s\n", repoPath)

	// Call the function to manage .gitkeep files (to be implemented in gitkeep package)
	err = gitkeep.ManageGitkeepFiles(repoPath)
	if err != nil {
		fmt.Printf("❌ Error managing .gitkeep files: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Operation completed successfully.")
}
