package workerpool

import (
	"log"
	"math"
	"time"
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

type StatsResponse struct {
	Message string    `json:"message"`
	Stats   Statistic `json:"stats"`
}

type PostBatchRequest struct {
	Data        []int `json:"data"`
	ChunkSize   int   `json:"chunkSize"`
	Concurrency int   `json:"concurrency"`
}

type PostBatchResult struct {
	Result int
	Error  error
}

type PostBatchResponse struct {
	Result int   `json:"result"`
	Error  error `json:"error,omitempty"`
}

func DivideIntDataIntoChunks(data []int, chunkSize int) [][]int {
	numChunks := int(math.Ceil(float64(len(data)) / float64(chunkSize)))
	chunks := make([][]int, numChunks)

	for i := 0; i < numChunks; i++ {
		startIndex := i * chunkSize
		endIndex := startIndex + chunkSize
		if endIndex > len(data) {
			endIndex = len(data)
		}
		chunks[i] = data[startIndex:endIndex]
	}

	return chunks
}

func DivideDataIntoChunks(data []DataChunk, chunkSize int) []DataChunk {
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

func ProcessPostDataChunk(chunkId int, chunk []int) (int, error) {
	// Perform some processing on the chunk
	result := 0
	for _, num := range chunk {
		result += num
	}

	// Simulate some delay
	time.Sleep(time.Second)

	log.Printf("Processing chunk: %v, result: %v", chunkId, result)
	// Return the result
	return result, nil
}

func ProcessDataChunk(chunk DataChunk) error {
	// Process the data chunk
	// You can perform any operations on the data here
	log.Printf("Processing chunk %d: %v\n", chunk.ChunkID, chunk.Data)
	return nil
}
