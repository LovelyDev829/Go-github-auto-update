package nik

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Repository represents the structure of a GitHub repository
type Repository struct {
	DefaultBranch string `json:"default_branch"`
}

// Comparison represents the structure of a branch comparison
type Comparison struct {
	AheadBy  int `json:"ahead_by"`
	BehindBy int `json:"behind_by"`
}

// Commit represents the structure of a GitHub commit
type Commit struct {
	SHA string `json:"sha"`
}

func fetchCommits(owner, repo, userAgent string) ([]Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var commits []Commit
	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func Get(owner string, repo string) (string, error) {
	userAgent := "StackOverflow-29845346"

	// Fetch commits and store SHAs in an array
	commits, err := fetchCommits(owner, repo, userAgent)
	if err != nil {
		// fmt.Println("Error fetching commits:", err)
		return "", err
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found")
	}

	// Return the SHA of the latest commit
	return commits[0].SHA, nil
}
