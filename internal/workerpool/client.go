package workerpool

import (
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"sync"
)

const poolSize = 10

var workerPool *WorkerPool
var onceInit sync.Once

func functionWithoutParams() error {
	fmt.Println("Executing functionWithoutParams")
	return nil
}

func functionWithParams(param int) error {
	fmt.Printf("Executing functionWithParams with param: %d\n", param)
	return nil
}

func InitWorkerPool() {
	onceInit.Do(func() {
		// Create a new worker pool
		wp, err := NewWorkerPool(poolSize)
		if err != nil {
			// Handle error
			fmt.Println("Error creating worker pool:", err)
			return
		}
		workerPool = wp
	})
}

func Close() {
	// Wait for all tasks to complete
	workerPool.WaitAll()

	// Release pool
	workerPool.Close()

	// Clean up the ants pool
	ants.Release()
}

func NewWorkerPoolInstance() (*WorkerPool, error) {
	if workerPool == nil {
		return nil, errors.New("worker pool not initiated")
	}

	return workerPool, nil
}

func (wp *WorkerPool) NewCall() error {

	// Create a slice of tasks
	tasks := []*Task{
		{Func: functionWithoutParams},
		{Func: func() error { return functionWithParams(1) }},
		{Func: functionWithoutParams},
		{Func: func() error { return functionWithParams(2) }},
		// Add more tasks as needed...
	}

	// Submit all tasks to the worker pool
	for _, task := range tasks {
		err := wp.SubmitTask(task)
		if err != nil {
			// Handle error
			fmt.Println("Error submitting task to worker pool from client:", err)
		}
	}

	return nil
}
