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

func CreateTask(fn TaskFunc) *Task {
	return &Task{
		Func: func() error {
			return fn()
		},
	}
}

func (wp *WorkerPool) SubmitNewTask(task Task) error {

	err := wp.SubmitTask(&task)
	if err != nil {
		fmt.Println("Error submitting task to worker pool from client:", err)
	}

	return err
}

func (wp *WorkerPool) GetStats() Statistic {
	return wp.getStatistics()
}
