package main

import (
	"encoding/json"
	"fmt"
	"log"
	nethttp "net/http"
	getonlinecommit "nik/package"
	globalvar "nik/package"
	"os"
	"os/exec"
	"sync"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	owner = "LovelyDev829"
	repo = "nik-hello-world"
	token = "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"
	interval = 5 * time.Second    // Check interval
	localPath = "../repo-pulled/"
	execCommand = localPath + "hello-world.exe"
	successHashFile = "LastSuccessCommitHash.txt"
	failedHashFile = "LastFailedCommitHash.txt"
)

var (
    supperFlag bool
    mutex      sync.Mutex
)

const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
)

func run_circle() {
	ticker := time.NewTicker(interval) // Check for updates once every 2 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Print current time
		fmt.Printf("%s----------------------------------------------------------------%s\n", Green, Reset)
		if !supperFlag {
			fmt.Printf("%sSupperFlag is turned off. you can turn it on at /update-supper-flag.%s\n", Red, Reset)
			continue
		}
		fmt.Printf("%sChecking Updates on Github at: %s%s\n", Green, time.Now().Format(time.RFC3339), Reset)

		// Check for changes in the repository
		changed, lastSuccessCommitHash, lastCommitOnline, err := hasChanges()
		if err != nil {
			fmt.Printf("%sError checking repository for changes%s%s\n", Red, err, Reset)
			continue
		}

		if changed {
			fmt.Printf("%sChanges Dectected.%s\n", Yellow, Reset)

			err := pull_repo()
			if err != nil {
				if err.Error() == "repository does not exist" {
					clone_repo()
				} else if err.Error() != "already up-to-date" {
					fmt.Printf("%sError updating repository: %s%s\n", Red, err.Error(), Reset)
					continue
				} else {
					continue
				}
				
			} else {
				fmt.Printf("%sRepository updated successfully%s\n", Yellow, Reset)
			}

			// Run the command after git pull or clone
			no_err := runCommand(execCommand)
			if no_err != nil {
				fmt.Printf("%sError running command: %s%s\n", Red, no_err, Reset)

				// Remember the failed commit
				err := globalvar.SetVariableToFile(failedHashFile, lastCommitOnline)
				if err != nil {
					fmt.Printf("%sError updating failedHashFile%s%s\n", Red, err, Reset)
				} else {
					fmt.Printf("%sFailedHashFile updated successfully%s\n", Yellow, Reset)
				}

				//revert
				no_err := revert(lastSuccessCommitHash)
				if no_err != nil {
					fmt.Printf("%sError Reverting to the last success commit%s%s\n", Red, no_err, Reset)
				} else {
					fmt.Printf("%sReverted the the last success commit successfully%s\n", Yellow, Reset)
				}

				//re-execute
				ne_err := runCommand(execCommand)
				if ne_err != nil {
					fmt.Printf("%s DANGER! Error re-running command: %s%s\n", Red, ne_err, Reset)
				} else {
					fmt.Printf("%sCommand re-executed successfully. %s\n", Yellow, Reset)
				}
			} else {
				fmt.Printf("%sCommand executed successfully. %s\n", Yellow, Reset)

				// Remember the success commit
				err := globalvar.SetVariableToFile(successHashFile, lastCommitOnline)
				if err != nil {
					fmt.Printf("%sError updating successHashFile%s%s\n", Red, err, Reset)
				} else {
					fmt.Printf("%sSuccessHashFile updated successfully%s\n", Yellow, Reset)
				}
			}

		} else {
			fmt.Println("No changes in the repository or the branch is up to date.")
		}
	}
}

func hasChanges() (bool, string, string, error) {
	lastSuccessCommitHash := globalvar.GetVariableFromFile(successHashFile)
	fmt.Println("Current Local Commit: ", lastSuccessCommitHash)
	
	lastFailedCommitHash := globalvar.GetVariableFromFile(failedHashFile)

	lastCommitOnline := getonlinecommit.Get(owner, repo, token)
	fmt.Println("Online Last Commit  : ", lastCommitOnline)

	if lastFailedCommitHash == lastCommitOnline {
		fmt.Println("The change is the failing one. Skip..")
		return false, lastSuccessCommitHash, lastCommitOnline, nil
	}

	return lastSuccessCommitHash != lastCommitOnline, lastSuccessCommitHash, lastCommitOnline, nil
}

func pull_repo() (error) {
	// Open the existing repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		fmt.Println("Error opening repository:", err)
		return err
	}

	// Create a new pull options struct
	pullOptions := &git.PullOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: owner, // Can be anything except an empty string
			Password: token,
		},
	}

	// Get the working directory for the repository
	worktree, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error getting worktree:", err)
		return err
	}

	// Perform a pull to update the local repository
	err = worktree.Pull(pullOptions)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("Repository is already up-to-date")
		} else {
			fmt.Println("Error pulling repository:", err)
		}
		return err
	} else{
		fmt.Println("Repository pulled successfully")
	}

	return nil
}

func clone_repo() {
	// Clone the repository
 
	url := fmt.Sprintf("https://github.com/%s/%s", owner, repo)
	_, err := git.PlainClone(localPath, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: owner, // Can be anything except an empty string
			Password: token,
		},
	})
	if err != nil {
		fmt.Println("Error cloning repository:", err)
		os.Exit(1)
	}
	fmt.Println("Repository cloned successfully to", localPath)
}

func runCommand(command string) error {
	cmd := exec.Command(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

//-----------------------
func cloneOrOpenRepo(directory, url, username, token string) (*git.Repository, error) {
    // Try to open the repo
    repo, err := git.PlainOpen(directory)
    if err == nil {
        return repo, nil
    }

    // Repo does not exist, so clone it
    repo, err = git.PlainClone(directory, false, &git.CloneOptions{
        URL: url,
        Auth: &http.BasicAuth{
            Username: username, // This can be anything when using personal access tokens
            Password: token,
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to clone the repository: %w", err)
    }
    return repo, nil
}

// checkoutCommit checks out the specified commit in the given repository.
func checkoutCommit(repo *git.Repository, commitSHA string) error {
    // Parse the commit hash
    hash := plumbing.NewHash(commitSHA)
    w, err := repo.Worktree()
    if err != nil {
        return fmt.Errorf("could not get worktree: %w", err)
    }

    // Checkout the commit
    err = w.Checkout(&git.CheckoutOptions{
        Hash: hash,
    })
    if err != nil {
        return fmt.Errorf("could not checkout: %w", err)
    }
    return nil
}

func revert(commitSHA string) error {

    repoURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)

    // Clone or open the repository
    repo, err := cloneOrOpenRepo(localPath, repoURL, "x-access-token", token)
    if err != nil {
        log.Fatalf("Error while cloning/opening repository: %v", err)
    }

    // Checkout the specific commit
    err = checkoutCommit(repo, commitSHA)
    if err != nil {
        log.Fatalf("Error while checking out commit: %v", err)
    }

    fmt.Println("Successfully checked out commit:", commitSHA)

	return nil
}

//api-------------------------------------------
// Handler to get the current value of supperFlag
func getFlagHandler(w nethttp.ResponseWriter, r *nethttp.Request) {
    mutex.Lock()
    defer mutex.Unlock()

    response := map[string]bool{"supperFlag": supperFlag}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        nethttp.Error(w, "Error generating response", nethttp.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

// Handler to update the value of supperFlag
func updateFlagHandler(w nethttp.ResponseWriter, r *nethttp.Request) {
    if r.Method != nethttp.MethodPost {
        nethttp.Error(w, "Invalid request method", nethttp.StatusMethodNotAllowed)
        return
    }

    var request struct {
        SupperFlag bool `json:"supperFlag"`
    }

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&request)
    if err != nil {
        nethttp.Error(w, "Invalid request body", nethttp.StatusBadRequest)
        return
    }

    mutex.Lock()
    supperFlag = request.SupperFlag
    mutex.Unlock()

    w.WriteHeader(nethttp.StatusOK)
    fmt.Fprintf(w, "SupperFlag updated to %v", supperFlag)
}

func main() {
    // Start the goroutine to print the time every 3 seconds
	supperFlag = true
    go run_circle()

    nethttp.HandleFunc("/get-supper-flag", getFlagHandler)
    nethttp.HandleFunc("/update-supper-flag", updateFlagHandler)

    fmt.Printf("%sStarting server at port 8080%s\n", Yellow, Reset)
    if err := nethttp.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("%sError starting server: %s%s\n", Red, err, Reset)
    }
}




