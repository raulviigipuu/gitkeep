package gitkeep

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/raulviigipuu/gitkeep/internal/gitutils"
	"github.com/raulviigipuu/gitkeep/internal/logx"
)

// ManageGitkeepFiles manages .gitkeep files in the specified repository
func ManageGitkeepFiles(repoPath string) error {
	// Convert repoPath to absolute path to handle relative paths correctly
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return err
	}

	return filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // Propagate the error upwards
		}

		// Skip non-directory entries
		if !d.IsDir() {
			return nil
		}

		// Check if the path is inside the .git directory
		if isInsideGitDir(absRepoPath, path) {
			return filepath.SkipDir // Skip the .git directory and its subdirectories
		}

		// Check if the directory is ignored by Git
		isIgnored, err := gitutils.IsPathIgnored(path)
		if err != nil {
			return err
		}
		if isIgnored {
			return nil // Skip ignored directories
		}

		// Check if the directory is empty
		isEmpty, err := isDirEmpty(path)
		if err != nil {
			return err
		}

		gitkeepPath := filepath.Join(path, ".gitkeep")
		if isEmpty {
			// Add .gitkeep file if the directory is empty
			if err := os.WriteFile(gitkeepPath, nil, 0644); err != nil {
				return err
			}
			fmt.Printf("✨ Created .gitkeep in: %s\n", path)
		} else {
			// Remove .gitkeep file if the directory is not empty
			if _, err := os.Stat(gitkeepPath); err == nil {
				if err := os.Remove(gitkeepPath); err != nil {
					return err
				}
				logx.Info(fmt.Sprintf("🧹 Removed .gitkeep from: %s\n", path))
			}
		}

		return nil
	})
}

// isInsideGitDir checks if a given path is inside the .git directory
func isInsideGitDir(repoPath, path string) bool {
	relPath, err := filepath.Rel(repoPath, path)
	if err != nil {
		return false
	}
	return strings.HasPrefix(relPath, ".git") || strings.Contains(relPath, "/.git/")
}

// isDirEmpty checks if the given directory is empty
func isDirEmpty(path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	// Consider a directory empty if it only contains a .gitkeep file
	return len(files) == 0 || (len(files) == 1 && files[0].Name() == ".gitkeep"), nil
}
