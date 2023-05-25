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

type DataChunk struct {
	ChunkID int   `json:"chunk_id"`
	Data    []int `json:"data"`
}

type BatchRequest struct {
	TotalData   int         `json:"total_data"`
	ChunkSize   int         `json:"chunk_size"`
	Concurrency int         `json:"concurrency"`
	Data        []DataChunk `json:"data"`
}

type BatchResponse struct {
	Message string `json:"message"`
}

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
	http.HandleFunc("/api/batch-process", handleBatchProcess)

	// Start the HTTP server
	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func divideDataIntoChunks(data []DataChunk, chunkSize int) []DataChunk {
	chunks := make([]DataChunk, 0)
	chunk := DataChunk{}

	for _, d := range data {
		chunk.Data = append(chunk.Data, d.Data...)

		if len(chunk.Data) >= chunkSize {
			chunks = append(chunks, chunk)
			chunk = DataChunk{}
		}
	}

	// Add the remaining elements as the last chunk
	if len(chunk.Data) > 0 {
		chunks = append(chunks, chunk)
	}

	return chunks
}

func handleBatchProcess(w http.ResponseWriter, r *http.Request) {

	wp, errIn := workerpool.NewWorkerPoolInstance()
	if errIn != nil {
		// Handle error
		log.Println("Error getting worker pool:", errIn)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var request BatchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Divide the data into smaller chunks
	chunks := divideDataIntoChunks(request.Data, request.ChunkSize)

	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for _, chunk := range chunks {
		go func(chunk DataChunk) {
			defer wg.Done()

			newTask := workerpool.CreateTask(func() error {
				return processDataChunk(chunk)
			})

			errNewTask := wp.SubmitNewTask(*newTask)
			if errNewTask != nil {
				log.Println("Error submitting task to worker pool:", errNewTask)
				return
			}

		}(chunk)
	}

	// Wait for all the chunks to be processed
	wg.Wait()

	response := BatchResponse{
		Message: "Data processed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
}

func processDataChunk(chunk DataChunk) error {
	// Process the data chunk
	// You can perform any operations on the data here
	log.Printf("Processing chunk %d: %v\n", chunk.ChunkID, chunk.Data)
	return nil
}

// Add context in the params
func initialiseServices() func() {
	workerpool.InitWorkerPool()

	return func() {
		workerpool.Close()
	}
}
