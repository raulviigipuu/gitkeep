package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

const gitKeeperFile = ".gitkeep"

// Function to process a directory
func processDirectory(path string, matcher gitignore.Matcher) error {
	if matcher.Match([]string{}, filepath.Base(path)) {
		return nil // Skip ignored paths
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	isOnlyGitKeeper := len(files) == 1 && files[0].Name() == gitKeeperFile
	if len(files) == 0 || isOnlyGitKeeper {
		if !isOnlyGitKeeper {
			gitKeeperPath := filepath.Join(path, gitKeeperFile)
			if err := os.WriteFile(gitKeeperPath, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("Added", gitKeeperPath)
		}
	} else {
		gitKeeperPath := filepath.Join(path, gitKeeperFile)
		if _, err := os.Stat(gitKeeperPath); err == nil {
			if err := os.Remove(gitKeeperPath); err != nil {
				return err
			}
			fmt.Println("Removed", gitKeeperPath)
		}
	}
	return nil
}

// Function to create a matcher from .gitignore file
func createGitignoreMatcher(rootDir string) (gitignore.Matcher, error) {
	gitignorePath := filepath.Join(rootDir, ".gitignore")
	patterns, err := gitignore.ReadPatterns(nil, []string{gitignorePath})
	if err != nil {
		if os.IsNotExist(err) {
			return gitignore.NewMatcher(nil), nil // .gitignore not found, return default
		}
		return nil, err
	}
	return gitignore.NewMatcher(patterns), nil
}

// Recursive function to traverse directories
func traverseDirectories(path string, matcher gitignore.Matcher) error {
	return filepath.WalkDir(path, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return processDirectory(currentPath, matcher)
		}
		return nil
	})
}

func main() {
	flag.Parse()
	rootDir := flag.Arg(0)
	if rootDir == "" {
		fmt.Println("Usage: gitkeep <root_directory>")
		return
	}

	matcher, err := createGitignoreMatcher(rootDir)
	if err != nil {
		fmt.Println("Error creating .gitignore matcher:", err)
		return
	}

	if err := traverseDirectories(rootDir, matcher); err != nil {
		fmt.Println("Error:", err)
	}
}
