package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

type Repo struct {
	Path string
	Name string
}

func LocalClone(private bool) error {
	repos := findGitRepos(".")
	selectedRepos := selectLocalRepos(repos)
	confirmAndInitRepos(selectedRepos, private)
	return nil
}

func findGitRepos(root string) []Repo {
	var repos []Repo
	var scannedDirs int

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Error reading directory %v: %v\n", root, err)
		return repos
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		scannedDirs++
		fmt.Printf("\rScanned %d directories...", scannedDirs)

		dirPath := filepath.Join(root, entry.Name())
		gitDir := filepath.Join(dirPath, ".git")

		if _, err := os.Stat(gitDir); err == nil {
			repos = append(repos, Repo{
				Path: dirPath,
				Name: entry.Name(),
			})
		}
	}

	fmt.Printf("\rScanned %d directories. Found %d Git repositories.\n", scannedDirs, len(repos))

	return repos
}

func selectLocalRepos(repos []Repo) []Repo {
	var options []string
	for _, repo := range repos {
		options = append(options, repo.Name)
	}

	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select repositories to initialize:",
		Options: options,
	}
	survey.AskOne(prompt, &selected)

	var selectedRepos []Repo
	for _, name := range selected {
		for _, repo := range repos {
			if repo.Name == name {
				selectedRepos = append(selectedRepos, repo)
				break
			}
		}
	}

	return selectedRepos
}

func confirmAndInitRepos(repos []Repo, private bool) {
	visibilityStr := "public"
	if private {
		visibilityStr = "private"
	}

	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Initialize %d repositories as %s?", len(repos), visibilityStr),
	}
	survey.AskOne(prompt, &confirm)

	if !confirm {
		fmt.Println("Radicalization cancelled.")
		return
	}

	fmt.Printf("Initializing %d repositories as %s...\n", len(repos), visibilityStr)

	for i, repo := range repos {
		fmt.Printf("Initializing %s (%d/%d)...\n", repo.Name, i+1, len(repos))
		err := runRadInit(repo.Path, repo.Name, private)
		if err != nil {
			color.Red("Error initializing %s: %v\n", repo.Name, err)
		} else {
			color.Green("Initialized %s as %s\n", repo.Name, visibilityStr)
		}
	}

	fmt.Println("Radicalization Complete")
}

func runRadInit(path, name string, private bool) error {
	visibilityFlag := "--public"
	if private {
		visibilityFlag = "--private"
	}

	cmd := exec.Command("rad", "init", "--name", name, "--description", "", visibilityFlag, "--no-confirm")
	cmd.Dir = path
	return cmd.Run()
}
