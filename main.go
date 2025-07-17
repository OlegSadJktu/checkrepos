package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	gogit "github.com/go-git/go-git/v5"
)

type repo struct {
	branch   string
	str      string
	path     string
	unstaged int
}

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

	var repos []repo
	// Check all subdirectories for .git
	for _, subdir := range subdirs {
		repoStruct := repo{}
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

		repoStruct.path = subdirPath

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
		repoStruct.unstaged = unstagedCount

		// Get current branch
		branch, err := repo.Head()
		if err != nil {
			fmt.Printf("Error getting current branch at %s: %v\n", subdirPath, err)
			continue
		}
		repoStruct.branch = branch.Name().Short()

		repoStruct.str = fmt.Sprintf("uncommitted files: %d", unstagedCount)
		repos = append(repos, repoStruct)
	}

	maxlen := 0
	for _, repo := range repos {
		if len(repo.path) > maxlen {
			maxlen = len(repo.path)
		}
	}

	for _, repo := range repos {
		if repo.unstaged > 0 {
			color.Red("%-*s    %s\n", maxlen, repo.path, repo.str)
		} else {
			if repo.branch == "develop" {
				color.Green("%-*s    %s\n", maxlen, repo.path, repo.str)
			} else {
				color.Yellow("%-*s    %s; branch: %s\n", maxlen, repo.path, repo.str, repo.branch)
			}
		}
	}

}
