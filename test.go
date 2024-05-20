package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func verifyCommitSignature(repoPath string, commitHash string) (bool, error) {
    // Prepare the git command
    cmd := exec.Command("git", "-C", repoPath, "verify-commit", commitHash)
    
    // Capture the output
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out

    // Run the command
    err := cmd.Run()
    if err != nil {
        fmt.Println("Command output:", out.String())
        return false, err
    }

    // Check output for verification success
    if bytes.Contains(out.Bytes(), []byte("Good signature")) {
        return true, nil
    }
    return false, fmt.Errorf("signature verification failed: %s", out.String())
}

func main() {
    repoPath := "../repo-pulled"
    commitHash := "7fee4f1cd7118429ace2dc4d9d871e67fceae388" // Replace with actual commit hash

    // Verify the commit signature
    valid, err := verifyCommitSignature(repoPath, commitHash)
    if err != nil {
        fmt.Println("Error verifying commit:", err)
        return
    }

    if valid {
        fmt.Println("Commit signature is valid.")
    } else {
        fmt.Println("Commit signature is invalid.")
    }
}
