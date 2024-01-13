package gitutils

import (
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// CheckIfGitRepo checks if the given path is within a Git repository
func CheckIfGitRepo(path string) (bool, string) {
	maxDepth := 100 // Maximum directory depth to prevent potential infinite loops

	originalPath := path
	for depth := 0; depth < maxDepth; depth++ {
		_, err := git.PlainOpen(path)
		if err == nil {
			return true, path // Found a Git repository, return its path
		}

		if path == "/" || path == "." {
			break // Reached the root without finding a Git repo
		}

		// Move to the parent directory
		path = filepath.Dir(path)
	}

	return false, originalPath // Not a Git repository, return original path
}

// IsGitInstalled checks if Git is installed and available in the system
func IsGitInstalled() bool {
	cmd := exec.Command("git", "--version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// IsPathIgnored checks if the given path is ignored by Git
func IsPathIgnored(path string) (bool, error) {
	// Check if the path is in a Git repository and find the repository root
	isInGitRepo, repoPath := CheckIfGitRepo(path)
	if !isInGitRepo {
		return false, nil
	}

	// Resolve the path relative to the repository root
	relPath, err := filepath.Rel(repoPath, path)
	if err != nil {
		return false, err
	}

	// Run the check-ignore command
	cmd := exec.Command("git", "check-ignore", "-q", relPath)
	cmd.Dir = repoPath // Set the working directory to the repository root
	err = cmd.Run()

	return err == nil, nil // if err is nil, the command executed successfully, indicating the path is ignored
}
