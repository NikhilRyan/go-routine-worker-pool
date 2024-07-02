# Worker Pool with Batching

The Worker Pool with Batching is a Go library that provides a scalable and concurrent solution for processing large data sets efficiently. This library allows you to leverage goroutines and a worker pool to process data in smaller batches concurrently, improving overall performance and throughput.

## Features

- Concurrent processing using a worker pool
- Data processing in smaller batches
- Scalable and adjustable pool size
- Graceful shutdown
- Error handling

## Installation

To install the Worker Pool with Batching library, use the following command:

```
go get github.com/NikhilRyan/go-routine-worker-pool
```

## Dependencies

The Worker Pool with Batching library relies on the `github.com/panjf2000/ants/v2` library for managing goroutines. Please ensure you have it installed. You can install it using the following command:

```
go get github.com/panjf2000/ants/v2
```

## Usage

Import the Worker Pool with Batching library into your Go code:

```go
import (
    "github.com/NikhilRyan/go-routine-worker-pool"
)
```

### Initialization

Initialize the worker pool with the desired pool size:

```go
poolSize := 10 // Set the desired pool size
wp := workerpool.NewWorkerPool(poolSize)
```

### Submitting Tasks

Submit tasks to the worker pool for processing:

```go
// Define your task function
taskFunc := func(data interface{}) error {
    // Process the data
    return nil
}

// Submit a task to the worker pool
err := wp.SubmitTask(taskFunc, data)
if err != nil {
    // Handle error
}
```

### Processing Data in Batches

To process large data sets in smaller batches, use the `ProcessDataInBatches` function:

```go
// Define your batch processing function
batchFunc := func(batch []interface{}) {
    for _, data := range batch {
        // Process the data
    }
}

// Process data in batches
err := wp.ProcessDataInBatches(data, batchSize, batchFunc)
if err != nil {
    // Handle error
}
```

### Graceful Shutdown

To gracefully shut down the worker pool, ensuring all tasks are completed, use the `Close` function:

```go
wp.Close()
```

### Error Handling

The Worker Pool with Batching library provides error handling mechanisms. You can propagate errors from your task functions and handle them appropriately:

```go
taskFunc := func(data interface{}) error {
    // Process the data
    if err := doSomething(data); err != nil {
        return err // Propagate the error
    }
    return nil
}
```

### Monitoring Pool Statistics

You can monitor the running and available goroutines in the worker pool:

```go
running := wp.Running() // Get the number of running goroutines
available := wp.Available() // Get the number of available goroutines
```

### Handling Termination Signal

To handle termination signals and ensure graceful shutdown of the worker pool, use the `WaitForTerminationSignal` function:

```go
workerpool.WaitForTerminationSignal(wp)
```

## Contributing

Contributions are welcome! If you have any suggestions, feature requests, or bug reports, please open an issue or submit a pull request.
