package main

import (
	"log"
	"net/http"

	"worker-pool/internal/workerpool"
)

func main() {

	closeFn := initialiseServices()
	defer closeFn()

	// Define an HTTP handler for the API endpoint
	http.HandleFunc("/api/call-function", func(w http.ResponseWriter, r *http.Request) {
		// Call the functionToCall asynchronously using the worker pool
		wp, errIn := workerpool.NewWorkerPoolInstance()
		if errIn != nil {
			// Handle error
			log.Println("Error getting worker pool:", errIn)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		errNewCall := wp.NewCall()
		if errNewCall != nil {
			log.Println("Error submitting task to worker pool:", errNewCall)
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
