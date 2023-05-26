package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
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

	// Define your API routes and handlers
	http.HandleFunc("/api/new-task", handleNewTask)
	http.HandleFunc("/api/pre-batch-process", handlePreBatchProcess)
	http.HandleFunc("/api/post-batch-process", handlePostBatchProcess)

	// Start the HTTP server
	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePreBatchProcess(w http.ResponseWriter, r *http.Request) {

	wp, errIn := workerpool.NewWorkerPoolInstance()
	if errIn != nil {
		// Handle error
		log.Println("Error getting worker pool:", errIn)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Pool Stats in Start
	fmt.Println("Workerpool stats in start: ", wp.GetStats())

	var request workerpool.BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Divide the data into smaller chunks
	chunks := workerpool.DivideDataIntoChunks(request.Data, request.ChunkSize)

	// Create a channel to receive the results of each chunk processing
	results := make(chan error, len(chunks))

	// Create a wait group to track the completion of the tasks for the current request
	var wg sync.WaitGroup
	wg.Add(len(chunks)) // Set the wait group count to the number of chunks

	// Submit a separate task to the worker pool for each chunk
	for _, chunk := range chunks {
		chunk := chunk // Create a new variable to avoid closure-related issues

		newTask := workerpool.CreateTask(func() error {
			defer wg.Done() // Mark the task as done when it completes
			return workerpool.ProcessDataChunk(chunk)
		})

		err := wp.SubmitNewTask(*newTask)

		if err != nil {
			log.Println("Error submitting task to worker pool:", err)
			results <- err // Report the error to the results channel
			wg.Done()      // Mark the task as done in case of error
		}
	}

	// Wait for all the tasks associated with the current request to be completed
	go func() {
		wg.Wait()
		close(results) // Close the results channel when all tasks are finished
	}()

	// Collect the results from the results channel
	var errors []error
	for err := range results {
		if err != nil {
			errors = append(errors, err)
		}
	}

	// Check if any errors occurred during processing
	if len(errors) > 0 {
		// Handle the errors accordingly
		log.Println("Errors occurred during processing:", errors)
		http.Error(w, "Errors occurred during processing", http.StatusInternalServerError)
		return
	}

	response := workerpool.BatchResponse{
		Message: "Data processed successfully",
	}

	// Pool Stats in End
	fmt.Println("Workerpool stats in end: ", wp.GetStats())

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlePostBatchProcess(w http.ResponseWriter, r *http.Request) {

	wp, errIn := workerpool.NewWorkerPoolInstance()
	if errIn != nil {
		// Handle error
		log.Println("Error getting worker pool:", errIn)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Pool Stats in Start
	fmt.Println("Workerpool stats in start: ", wp.GetStats())

	// Parse the request body to get the data
	var requestData workerpool.PostBatchRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Divide the data into chunks
	chunks := workerpool.DivideIntDataIntoChunks(requestData.Data, requestData.ChunkSize)

	// Create a channel to collect the results of each task
	results := make(chan workerpool.PostBatchResult, len(chunks))

	// Submit each chunk as a task to the worker pool
	for i, chunk := range chunks {
		fmt.Printf("chunk %v, %v\n", i, chunk)
		newTask := workerpool.CreateTask(func() error {
			result, err := workerpool.ProcessPostDataChunk(i, chunk)
			results <- workerpool.PostBatchResult{Result: result, Error: err}
			return err
		})

		// Pool Stats in between
		fmt.Println("Workerpool stats in between: ", wp.GetStats())

		err := wp.SubmitTask(newTask)
		if err != nil {
			log.Println("Error submitting task to worker pool:", err)
			results <- workerpool.PostBatchResult{Error: err}
		}
	}

	// Collect the results from the tasks
	var batchResults []workerpool.PostBatchResult
	for range chunks {
		result := <-results
		batchResults = append(batchResults, result)
	}

	// Close the results channel
	close(results)

	// Combine the results
	var combinedResult workerpool.PostBatchResult
	for _, result := range batchResults {
		combinedResult.Result += result.Result
		if result.Error != nil {
			combinedResult.Error = result.Error
		}
	}

	// Return the combined result
	response := workerpool.PostBatchResponse{
		Result: combinedResult.Result,
		Error:  combinedResult.Error,
	}

	// Pool Stats in End
	fmt.Println("Workerpool stats in end: ", wp.GetStats())

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func handleNewTask(w http.ResponseWriter, r *http.Request) {

	// Call the functionToCall asynchronously using the worker pool
	wp, errIn := workerpool.NewWorkerPoolInstance()
	if errIn != nil {
		// Handle error
		log.Println("Error getting worker pool:", errIn)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Pool Stats in Start
	fmt.Println("Workerpool stats in start: ", wp.GetStats())

	newTask := workerpool.CreateTask(func() error {
		return functionWithParams(1)
	})

	errNewTask := wp.SubmitNewTask(*newTask)
	if errNewTask != nil {
		log.Println("Error submitting task to worker pool:", errNewTask)
		return
	}

	// Pool Stats in end
	fmt.Println("Workerpool stats in end: ", wp.GetStats())

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

// Add context in the params
func initialiseServices() func() {
	workerpool.InitWorkerPool()

	return func() {
		workerpool.Close()
	}
}
