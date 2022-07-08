package main

import (
	"flag"

	"github.com/alphadose/logstreamer/utils"
)

func main() {
	flag.StringVar(&file, "f", "./data.txt", "Absolute path to file")
	flag.StringVar(&mongoURI, "mongo", "mongodb://localhost:27017", "MongoDB URI")
	flag.StringVar(&grpcURI, "grpc", "localhost:3002", "URI of GRPC server")
	flag.Uint64Var(&batchSize, "bsize", 200, "Batch size of upload operations (restriction helpful in cases of file_size > 16 GB)")
	flag.BoolVar(&parallel, "parallel", false, "Should storage upload operations run in parallel?")
	flag.Parse()

	utils.LogInfo("Main-Start", "Starting Operation")
	process()
	utils.LogInfo("Main-End", "Successfully Completed")
}
