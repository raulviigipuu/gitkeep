package gitkeep

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/raulviigipuu/gitkeep/internal/gitutils"
	"github.com/raulviigipuu/gitkeep/internal/logx"
)

func ManageGitkeepFiles(repoPath string) error {
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return err
	}

	return checkAndManageDirectory(absRepoPath, absRepoPath)
}

func checkAndManageDirectory(currentPath, repoRoot string) error {
	// Normalize and skip .git folder (Windows-safe)
	relPath, err := filepath.Rel(repoRoot, currentPath)
	if err != nil {
		return err
	}
	relPath = filepath.ToSlash(relPath)
	lowerPath := strings.ToLower(relPath)
	if lowerPath == ".git" || strings.HasPrefix(lowerPath, ".git/") {
		return nil // skip .git dir
	}

	// Check if ignored by Git
	isIgnored, err := gitutils.IsPathIgnored(currentPath)
	if err != nil {
		return err
	}
	if isIgnored {
		return nil // skip ignored dirs
	}

	// Read dir entries
	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return err
	}

	isEmpty := true
	for _, entry := range entries {
		if entry.Name() == ".gitkeep" || strings.ToLower(entry.Name()) == "thumbs.db" {
			continue
		}
		isEmpty = false
		break
	}

	gitkeepPath := filepath.Join(currentPath, ".gitkeep")
	if isEmpty {
		// Add .gitkeep if not exists
		if _, err := os.Stat(gitkeepPath); os.IsNotExist(err) {
			err := os.WriteFile(gitkeepPath, nil, 0644)
			if err != nil {
				return err
			}
			logx.Info(fmt.Sprintf("✨ Created .gitkeep in: %s", currentPath))
		}
	} else {
		// Remove if exists
		if _, err := os.Stat(gitkeepPath); err == nil {
			if err := os.Remove(gitkeepPath); err != nil {
				return err
			}
			logx.Info(fmt.Sprintf("🧹 Removed .gitkeep from: %s", currentPath))
		}
	}

	// Recurse into subdirectories
	for _, entry := range entries {
		if entry.IsDir() {
			if err := checkAndManageDirectory(filepath.Join(currentPath, entry.Name()), repoRoot); err != nil {
				return err
			}
		}
	}

	return nil
}
