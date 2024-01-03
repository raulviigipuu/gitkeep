package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const gitKeeperFile = ".gitkeep"

func processDirectory(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Check if directory is empty or contains only .gitkeep
	isOnlyGitKeeper := len(files) == 1 && files[0].Name() == gitKeeperFile
	if len(files) == 0 || isOnlyGitKeeper {
		// Add .gitkeep if it's not there
		if !isOnlyGitKeeper {
			gitKeeperPath := filepath.Join(path, gitKeeperFile)
			if err := os.WriteFile(gitKeeperPath, []byte{}, 0644); err != nil {
				return err
			}
			fmt.Println("Added", gitKeeperPath)
		}
	} else {
		// Remove .gitkeep if it's there
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

func traverseDirectories(path string) error {
	return filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return processDirectory(currentPath)
		}
		return nil
	})
}

func main() {
	flag.Parse()
	rootDir := flag.Arg(0)
	if rootDir == "" {
		fmt.Println("Usage: gitkeeper <root_directory>")
		return
	}

	if err := traverseDirectories(rootDir); err != nil {
		fmt.Println("Error:", err)
	}
}
