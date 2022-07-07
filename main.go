/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"flag"

	"github.com/alphadose/logstreamer/grpc"
	"github.com/alphadose/logstreamer/utils"
)

func main() {
	utils.LogInfo("m", "m")
	var (
		file      string
		mongoURI  string
		grpcURI   string
		batchSize uint64
		parallel  bool
	)
	flag.StringVar(&file, "f", "./data.txt", "Absolute path to file")
	flag.StringVar(&mongoURI, "mongo", "mongodb://localhost:27017", "MongoDB URI")
	flag.StringVar(&grpcURI, "grpc", "localhost"+grpc.Port, "URI of GRPC server")
	flag.Uint64Var(&batchSize, "bsize", 200, "Batch size of upload operations (restriction helpful in cases of file_size > 16 GB)")
	flag.BoolVar(&parallel, "parallel", false, "Should storage upload operations run in parallel?")
	flag.Parse()

}
