package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/raulviigipuu/gitkeep/internal/gitkeep"
	"github.com/raulviigipuu/gitkeep/internal/gitutils"
	"github.com/raulviigipuu/gitkeep/internal/logx"
)

var Version = "dev" // Will be overridden at build time using -ldflags

func main() {

	logx.Init(nil)

	// Define flags
	versionFlag := flag.Bool("v", false, "Show version and exit")
	helpFlag := flag.Bool("h", false, "Show help and exit")
	flag.Parse()

	// Handle -h
	if *helpFlag {
		printHelp()
		return
	}

	// Handle -v flag
	if *versionFlag {
		fmt.Println("gitkeep", Version)
		return
	}

	// Remaining args (optional path argument)
	args := flag.Args()
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	// Check if Git is installed
	if !gitutils.IsGitInstalled() {
		logx.Fatal("😕 Git is not installed. Please install Git and try again.")
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		logx.Fatal(fmt.Sprintf("🚫 Error resolving absolute path: %s", err))
	}

	// Ensure the path is inside a Git repository
	isInGitRepo, repoPath := gitutils.CheckIfGitRepo(absPath)
	if !isInGitRepo {
		logx.Fatal(fmt.Sprintf("🔍 The path %s is not within a Git repository.", absPath))
	}

	logx.Info(fmt.Sprintf("🚀 Managing .gitkeep files in directory: %s", repoPath))

	// Run core logic
	if err := gitkeep.ManageGitkeepFiles(repoPath); err != nil {
		logx.Fatal(fmt.Sprintf("❌ Error managing .gitkeep files: %s", err))
	}

	logx.Info("✅ Operation completed successfully.")
}

func printHelp() {
	fmt.Println("gitkeep - Add and remove .gitkeep files in a directory (git repo)")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitkeep [options] [path]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -v       Show version and exit")
	fmt.Println("  -h       Show help and exit")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  path     Optional path to the directory (default: current directory)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitkeep")
	fmt.Println("  gitkeep -v")
	fmt.Println("  gitkeep ./path/to/dir")
	fmt.Println()
}
