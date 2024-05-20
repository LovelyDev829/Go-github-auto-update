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


func fetchCommits(owner, repo, token, userAgent string) ([]Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "token "+token)

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

func Get(owner string, repo string, token string) string {
	userAgent := "StackOverflow-29845346"

	// Fetch commits and store SHAs in an array
	commits, err := fetchCommits(owner, repo, token, userAgent)
	if err != nil {
		// fmt.Println("Error fetching commits:", err)
		return err.Error()
	}
	// Create a slice to store commit SHAs
	var commitSHAs []string
	for _, commit := range commits {
		commitSHAs = append(commitSHAs, commit.SHA)
	}
	return commitSHAs[0]
}
