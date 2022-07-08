package main

import (
	"flag"

	"github.com/alphadose/logstreamer/grpc"
)

func main() {
	var port uint64
	flag.Uint64Var(&port, "port", 3002, "port for running the GRPC server")
	flag.Parse()
	if err := grpc.ListenAndServe(port); err != nil {
		println(err.Error())
	}
}
