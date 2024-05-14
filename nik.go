package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	repoPath = "../nik-hello-world/" // Modify this path to the location of your local repo
)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
)

func main() {
	ticker := time.NewTicker(2 * time.Second) // Check for updates once a day
	defer ticker.Stop()

	for range ticker.C {
		// Print current time
		// fmt.Println("Checking Updates on Github at:", time.Now().Format(time.RFC3339))
		fmt.Printf("%sChecking Updates on Github at: %s%s\n", Green, time.Now().Format(time.RFC3339), Reset)

		err := gitPull(repoPath)
		if err != nil {
			fmt.Println("Error updating repository:", err)
		} else {
			// fmt.Println("Repository updated successfully.")
		}
	}
}

func gitPull(repoDir string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
