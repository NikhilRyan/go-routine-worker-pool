package workerpool

import (
	"errors"
	"sync"

	"github.com/panjf2000/ants/v2"
)

type Task struct {
	Func func() error
}

type WorkerPool struct {
	pool       *ants.Pool
	waitGroup  *sync.WaitGroup
	completion chan struct{}
}

func NewWorkerPool(poolSize int) (*WorkerPool, error) {
	p, err := ants.NewPool(poolSize)
	if err != nil {
		return nil, err
	}

	return &WorkerPool{
		pool:       p,
		waitGroup:  &sync.WaitGroup{},
		completion: make(chan struct{}),
	}, nil
}

func (wp *WorkerPool) SubmitTask(task *Task) error {
	if wp == nil {
		return errors.New("worker pool is not initialized")
	}

	wp.waitGroup.Add(1)

	err := wp.pool.Submit(func() {
		defer wp.waitGroup.Done()
		err := task.Func()
		if err != nil {
			// Handle task error if needed
		}
	})

	if err != nil {
		wp.waitGroup.Done()
		return err
	}

	return nil
}

func (wp *WorkerPool) WaitAll() {
	if wp == nil {
		return
	}

	go func() {
		wp.waitGroup.Wait()
		close(wp.completion)
	}()

	<-wp.completion
}

func (wp *WorkerPool) Close() {
	if wp == nil {
		return
	}

	wp.pool.Release()
}

type WaitGroup struct {
	wg sync.WaitGroup
}

func (wg *WaitGroup) Add(num int) {
	if wg == nil {
		return
	}

	wg.wg.Add(num)
}

func (wg *WaitGroup) Done() {
	if wg == nil {
		return
	}

	wg.wg.Done()
}

func (wg *WaitGroup) Wait() {
	if wg == nil {
		return
	}

	wg.wg.Wait()
}
