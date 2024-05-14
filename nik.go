package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	repoPath = "../2_Cloning/Three.js-my-room/" // Modify this path to the location of your local repo
)

func main() {
	ticker := time.NewTicker(2 * time.Second) // Check for updates once a day
	defer ticker.Stop()

	for range ticker.C {
		err := gitPull(repoPath)
		if err != nil {
			fmt.Println("Error updating repository:", err)
		} else {
			fmt.Println("Repository updated successfully.")
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
