package main

import (
	"fmt"

	"github.com/alphadose/logstreamer/grpc"
)

func main() {
	client := grpc.NewClient("localhost" + grpc.Port)
	data, err := client.Consume(5)
	if err != nil {
		println(err.Error())
	}
	for idx := range data {
		fmt.Printf("%#v\n", data[idx])
	}
}
