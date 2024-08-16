package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

type Repo struct {
	Path string
	Name string
}

var isPrivate bool

func main() {
	flag.BoolVar(&isPrivate, "private", false, "Initialize repositories as private")
	flag.Parse()

	repos := findGitRepos(".")
	selectedRepos := selectRepos(repos)
	confirmAndInitRepos(selectedRepos)
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

	// Clear the progress line and print final count
	fmt.Printf("\rScanned %d directories. Found %d Git repositories.\n", scannedDirs, len(repos))

	return repos
}

func selectRepos(repos []Repo) []Repo {
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

func confirmAndInitRepos(repos []Repo) {
	visibilityStr := "public"
	if isPrivate {
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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool)

	go func() {
		for i, repo := range repos {
			select {
			case <-interrupt:
				fmt.Println("\nInterrupted. Stopping initialization process.")
				done <- true
				return
			default:
				fmt.Printf("Initializing %s (%d/%d)...\n", repo.Name, i+1, len(repos))
				err := runRadInit(repo.Path, repo.Name)
				if err != nil {
					color.Red("Error initializing %s: %v\n", repo.Name, err)
				} else {
					color.Green("Initialized %s as %s\n", repo.Name, visibilityStr)
				}
			}
		}
		done <- true
	}()

	<-done
	fmt.Println("Radicalization Complete")
}

func runRadInit(path, name string) error {
	visibilityFlag := "--public"
	if isPrivate {
		visibilityFlag = "--private"
	}

	cmd := exec.Command("rad", "init", "--name", name, "--description", "", visibilityFlag, "--no-confirm")
	cmd.Dir = path
	return cmd.Run()
}
