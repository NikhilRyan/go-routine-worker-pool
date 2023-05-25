package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"worker-pool/internal/workerpool"
)

func functionWithoutParams() error {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Executing functionWithoutParams")
	return nil
}

func functionWithParams(param int) error {
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Executing functionWithParams with param: %d\n", param)
	return nil
}

func main() {

	closeFn := initialiseServices()
	defer closeFn()

	// Define an HTTP handler for the API endpoint
	http.HandleFunc("/api/new-task", func(w http.ResponseWriter, r *http.Request) {
		// Call the functionToCall asynchronously using the worker pool
		wp, errIn := workerpool.NewWorkerPoolInstance()
		if errIn != nil {
			// Handle error
			log.Println("Error getting worker pool:", errIn)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		newTask := workerpool.CreateTask(func() error {
			return functionWithParams(1)
		})

		errNewTask := wp.SubmitNewTask(*newTask)
		if errNewTask != nil {
			log.Println("Error submitting task to worker pool:", errNewTask)
			return
		}

		// Return a success response
		w.WriteHeader(http.StatusOK)
	})

	// Start the HTTP server
	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Add context in the params
func initialiseServices() func() {
	workerpool.InitWorkerPool()

	return func() {
		workerpool.Close()
	}
}
