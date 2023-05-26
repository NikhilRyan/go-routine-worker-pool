package workerpool

import "log"

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

func ProcessDataChunk(chunk DataChunk) error {
	// Process the data chunk
	// You can perform any operations on the data here
	log.Printf("Processing chunk %d: %v\n", chunk.ChunkID, chunk.Data)
	return nil
}
