package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	repoPath    = "../nik-hello-world/" // Modify this path to the location of your local repo
	execCommand = repoPath + "hello-world.exe" // Command to run after git pull
)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
)

var failedCommit string // Store the hash of the last failed commit

func main() {
	ticker := time.NewTicker(2 * time.Second) // Check for updates once every 2 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Print current time
		fmt.Printf("%s----------------------------------------------------------------%s\n", Green, Reset)
		fmt.Printf("%sChecking Updates on Github at: %s%s\n", Green, time.Now().Format(time.RFC3339), Reset)

		// Check for changes in the repository
		changed, statusOutput, err := hasChanges(repoPath)
		if err != nil {
			fmt.Printf("%sError checking repository for changes: %s%s\n", Red, err, Reset)
			continue
		}

		if changed {
			fmt.Printf("%sChanges in the repository:%s\n", Yellow, Reset)
			fmt.Println(statusOutput)

			err := gitPull(repoPath)
			if err != nil {
				fmt.Printf("%sError updating repository: %s%s\n", Red, err, Reset)
			} else {
				fmt.Printf("%sRepository updated successfully: %s\n", Yellow, Reset)

				// Run the command after git pull
				err := runCommand(execCommand)
				if err != nil {
					fmt.Printf("%sError running command: %s%s\n", Red, err, Reset)

					// Remember the failed commit
					failedCommit = getCurrentCommit(repoPath)
				} else {
					fmt.Printf("%sCommand executed successfully. %s\n", Yellow, Reset)
				}
			}
		} else {
			fmt.Println("No changes in the repository or the branch is up to date.")

			// If the failed commit is detected as a change, skip pulling
			if failedCommit != "" {
				currentCommit := getCurrentCommit(repoPath)
				if currentCommit == failedCommit {
					fmt.Printf("%sSkipping pull as the failed commit is detected as a change.%s\n", Yellow, Reset)
					continue
				}
			}
		}
	}
}

func hasChanges(repoDir string) (bool, string, error) {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = repoDir
	err := cmd.Run()
	if err != nil {
		return false, "", err
	}

	// Check if the local branch is ahead of the remote branch
	cmd = exec.Command("git", "status", "-uno", "-sb")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return false, "", err
	}
	statusOutput := string(output)
	return strings.Contains(statusOutput, "behind") || strings.Contains(statusOutput, "diverged"), statusOutput, nil
}

func gitPull(repoDir string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runCommand(command string) error {
	cmd := exec.Command(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func getCurrentCommit(repoDir string) string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
