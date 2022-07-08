package main

import (
	"io"
	"sync"
	"sync/atomic"

	"github.com/alphadose/logstreamer/grpc"
	"github.com/alphadose/logstreamer/mongo"
	"github.com/alphadose/logstreamer/types"
	"github.com/alphadose/logstreamer/utils"
)

var (
	file     string
	mongoURI string
	grpcURI  string
	// if the host system has enough power then network operations can be run asynchronously
	// each via a separate goroutine
	parallel bool
	// read file and upload in batches to reduce the network I/O pressure as well as host system memory
	// if file_size > 16 GB, it might not be loaded into main memory all at once due to hardware constraints
	// and even if its loaded it will put tremendous pressure for transmission over the wire in which case
	// significant network latency might be observable
	// More efficient to chunk data and then transmit over the wire
	batchSize uint64
)

// starts processing with the above populated global params
func process(mongoCollectionName ...string) {
	var (
		collName = "users"      // default
		wg       sync.WaitGroup // waitgroup for synchronization in case `-parallel` flag is specified
	)
	if len(mongoCollectionName) > 0 {
		collName = mongoCollectionName[0]
	}
	// Initialize storage links
	mongoStore := mongo.NewClient(mongoURI, collName)
	grpcStore := grpc.NewClient(grpcURI)

	// Initialize file reader
	reader := utils.NewFileReader[types.Payload](file)
	defer reader.Close()

	for {
		payloadBatch, err := reader.ReadLines(batchSize)
		if err == io.EOF {
			// Reached end of file
			if parallel {
				// synchronize all running upload goroutines and wait for them to either finish or log error before process exit
				wg.Wait()
			}
			if len(payloadBatch) > 0 {
				if err = processBatch(payloadBatch, mongoStore, grpcStore); err != nil {
					utils.LogError("Core-1", err)
				}
			}
			return
		}
		if err != nil {
			utils.GracefulExit("Core-2", err)
		}
		if parallel {
			wg.Add(1)
			go func() {
				if err := processBatch(payloadBatch, mongoStore, grpcStore); err != nil {
					utils.LogError("Core-Parallel-3", err)
				}
				wg.Done()
			}()
		} else if err = processBatch(payloadBatch, mongoStore, grpcStore); err != nil {
			utils.GracefulExit("Core-4", err)
		}
	}
}

// track the current batch being processed
var batchNumber uint64

// process a batch with both MongoDB and GRPC endpoints atomically
func processBatch(payloads []*types.Payload, m *mongo.Store, g *grpc.Client) error {
	utils.LogInfo("Core-Intermmediate", "Processing Batch: %d", atomic.AddUint64(&batchNumber, 1))
	return m.Upload(func() error { return g.Publish(payloads) }, payloads)
}
