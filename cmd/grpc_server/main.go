package main

import "github.com/alphadose/logstreamer/grpc"

func main() {
	if err := grpc.ListenAndServe(); err != nil {
		println(err.Error())
	}
}
