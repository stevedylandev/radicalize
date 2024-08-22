package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func RemoteClone() error {
	fmt.Print("Enter the name of GitHub user or org")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading the input. Please try again", err)
		return err
	}
	input = strings.TrimSuffix(input, "\n")

	apiURL := fmt.Sprintf("https://api.github.com/users/%s/repos?sort=created", input)

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Github API request failed with error", err)
		return err
	}

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

	newRepos := selectNewRepos(repos)

	for _, repo := range newRepos {
		fmt.Printf("Repository: %s\n", repo.FullName)
		fmt.Printf("Description: %s\n", repo.Description)
		fmt.Printf("URL: %s\n", repo.HTMLURL)
		fmt.Printf("---\n")
	}

	return nil
}

func selectNewRepos(repos []Repository) []Repository {
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
