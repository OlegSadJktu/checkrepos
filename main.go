package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	gogit "github.com/go-git/go-git/v5"
)

func main() {
	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Get all subdirectories
	subdirs, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	strs := []string{}
	paths := []string{}
	unstaged := []int{}
	// Check all subdirectories for .git
	for _, subdir := range subdirs {
		if !subdir.IsDir() {
			continue
		}
		subdirPath := filepath.Join(dir, subdir.Name())
		gitPath := filepath.Join(subdirPath, ".git")
		_, err := os.Stat(gitPath)
		if err != nil {
			fmt.Printf("Not a git repository: %s\n", subdirPath)
			continue
		}

		paths = append(paths, subdirPath)

		// Open the repository
		repo, err := gogit.PlainOpen(subdirPath)
		if err != nil {
			fmt.Printf("Error opening repository at %s: %v\n", subdirPath, err)
			continue
		}

		// Get working tree
		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Printf("Error getting worktree at %s: %v\n", subdirPath, err)
			continue
		}

		// Get status
		status, err := worktree.Status()
		if err != nil {
			fmt.Printf("Error getting status at %s: %v\n", subdirPath, err)
			continue
		}

		unstagedCount := 0
		for _, s := range status {
			if s.Worktree != gogit.UpdatedButUnmerged {
				unstagedCount++
			}
		}
		unstaged = append(unstaged, unstagedCount)

		strs = append(strs, fmt.Sprintf("uncommitted files: %d", unstagedCount))
	}

	maxlen := 0
	for _, str := range paths {
		if len(str) > maxlen {
			maxlen = len(str)
		}
	}

	for i, str := range strs {
		if unstaged[i] > 0 {
			color.Red("%-*s    %s\n", maxlen, paths[i], str)
		} else {
			color.Green("%-*s    %s\n", maxlen, paths[i], str)
		}
	}

}
