package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

func RemoteClone() error {
	fmt.Print("Enter the name of GitHub user or org: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading the input. Please try again", err)
		return err
	}
	input = strings.TrimSuffix(input, "\n")

	apiURL := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=created", input)

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("GitHub API request failed with error", err)
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	var repos []Repository
	err = json.Unmarshal(body, &repos)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	selectedRepos := selectNewRepos(repos)
	confirmAndCloneRepos(selectedRepos)

	return nil
}

func selectNewRepos(repos []Repository) []Repository {
	var options []string
	for _, repo := range repos {
		options = append(options, repo.Name)
	}

	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select repositories to clone and initialize:",
		Options: options,
	}
	survey.AskOne(prompt, &selected)

	var selectedRepos []Repository
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

func confirmAndCloneRepos(repos []Repository) {
	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Clone and initialize %d repositories?", len(repos)),
	}
	survey.AskOne(prompt, &confirm)

	if !confirm {
		fmt.Println("Radicalization cancelled.")
		return
	}

	fmt.Printf("Cloning and initializing %d repositories...\n", len(repos))

	for i, repo := range repos {
		fmt.Printf("Processing %s (%d/%d)...\n", repo.Name, i+1, len(repos))

		// Clone the repository
		err := cloneRepo(repo.CloneURL, repo.Name)
		if err != nil {
			color.Red("Error cloning %s: %v\n", repo.Name, err)
			continue
		}

		// Initialize with rad
		err = runRadInitRemote(repo.Name, repo.Name)
		if err != nil {
			color.Red("Error initializing %s: %v\n", repo.Name, err)
		} else {
			color.Green("Cloned and initialized %s\n", repo.Name)
		}
	}

	fmt.Println("Radicalization Complete")
}

func cloneRepo(url, name string) error {
	cmd := exec.Command("git", "clone", url, name)
	return cmd.Run()
}

func runRadInitRemote(path, name string) error {
	cmd := exec.Command("rad", "init", "--name", name, "--description", "", "--public", "--no-confirm")
	cmd.Dir = path
	return cmd.Run()
}
