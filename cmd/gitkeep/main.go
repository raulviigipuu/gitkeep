package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/raulviigipuu/gitkeep/internal/gitkeep"
	"github.com/raulviigipuu/gitkeep/internal/gitutils"
	"github.com/raulviigipuu/gitkeep/internal/logx"
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
		logx.Fatal("😕 Git is not installed. Please install Git and try again.")
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		logx.Fatal(fmt.Sprintf("🚫 Error resolving absolute path: %s\n", err))
	}

	// Check if the path is within a Git repository
	isInGitRepo, repoPath := gitutils.CheckIfGitRepo(absPath)
	if !isInGitRepo {
		logx.Fatal(fmt.Sprintf("🔍 The path %s is not within a Git repository.\n", absPath))
	}

	fmt.Printf("🚀 Managing .gitkeep files in repository: %s\n", repoPath)

	// Call the function to manage .gitkeep files (to be implemented in gitkeep package)
	err = gitkeep.ManageGitkeepFiles(repoPath)
	if err != nil {
		logx.Fatal(fmt.Sprintf("❌ Error managing .gitkeep files: %s\n", err))
	}

	logx.Info("✅ Operation completed successfully.")
}
