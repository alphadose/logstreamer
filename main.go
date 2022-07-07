package main

import (
	"flag"
	"io"

	"github.com/alphadose/logstreamer/grpc"
	"github.com/alphadose/logstreamer/mongo"
	"github.com/alphadose/logstreamer/types"
	"github.com/alphadose/logstreamer/utils"
)

func main() {
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
	flag.StringVar(&file, "f", "./data.txt", "Absolute path to file")
	flag.StringVar(&mongoURI, "mongo", "mongodb://localhost:27017", "MongoDB URI")
	flag.StringVar(&grpcURI, "grpc", "localhost"+grpc.Port, "URI of GRPC server")
	flag.Uint64Var(&batchSize, "bsize", 200, "Batch size of upload operations (restriction helpful in cases of file_size > 16 GB)")
	flag.BoolVar(&parallel, "parallel", false, "Should storage upload operations run in parallel?")
	flag.Parse()

	// Initialize storage links
	mongoStore := mongo.NewClient(mongoURI)
	grpcStore := grpc.NewClient(grpcURI)

	// Initialize file reader
	reader := utils.NewFileReader[types.Payload](file)
	defer reader.Close()

	for {
		payloadArray, err := reader.ReadLines(batchSize)
		if err == io.EOF {
			// Reached end of file
			return
		}
		if err != nil {
			utils.GracefulExit("Main-1", err)
		}
		if err := mongoStore.Upload(func() error { return grpcStore.Publish(payloadArray) }, nil); err != nil {
			utils.GracefulExit("Main-2", err)
		}
	}

}
