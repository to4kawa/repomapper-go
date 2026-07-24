package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: repomapper <path>")
		os.Exit(1)
	}

	repoPath := os.Args[1]
	absPath, err := filepath.Abs(repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	// Gitリポジトリとして開けるか確認
	repo, err := git.PlainOpen(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Not a git repository: %v\n", err)
		os.Exit(1)
	}

	// HEADのツリーを取得
	ref, err := repo.Head()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get HEAD: %v\n", err)
		os.Exit(1)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get commit: %v\n", err)
		os.Exit(1)
	}

	tree, err := commit.Tree()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get tree: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Repository: %s\n", absPath)
	fmt.Println("Files:")

	err = tree.Files().ForEach(func(f *object.File) error {
		fmt.Println(" -", f.Name)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing files: %v\n", err)
		os.Exit(1)
	}
}
