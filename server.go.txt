package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
    runningFlag bool
    mutex       sync.Mutex
)

// Handler to get the current value of runningFlag
func getFlagHandler(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    defer mutex.Unlock()

    response := map[string]bool{"runningFlag": runningFlag}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error generating response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

// Handler to update the value of runningFlag
func updateFlagHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        RunningFlag bool `json:"runningFlag"`
    }

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&request)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    mutex.Lock()
    runningFlag = request.RunningFlag
    mutex.Unlock()

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "runningFlag updated to %v", runningFlag)
}

// Function to print the current time every 3 seconds
func printCurrentTime() {
    for {
        time.Sleep(3 * time.Second)
        fmt.Println("Current time:", time.Now().Format(time.RFC1123))
    }
}

func main() {
    // Start the goroutine to print the time every 3 seconds
    go printCurrentTime()

    http.HandleFunc("/getFlag", getFlagHandler)
    http.HandleFunc("/updateFlag", updateFlagHandler)

    fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Error starting server: %s\n", err)
    }
}
