package main

import (
	"flag"
	"fmt"

	"github.com/alphadose/logstreamer/grpc"
)

func main() {
	var (
		count int64
		url   string
	)
	flag.Int64Var(&count, "c", 1<<61, "Number of objects to consume")
	flag.StringVar(&url, "url", "localhost:3002", "URL of GRPC servr")
	flag.Parse()
	client := grpc.NewClient(url)
	data, err := client.Consume(count)
	if err != nil {
		println(err.Error())
	}
	for idx := range data {
		fmt.Printf("%#v\n", data[idx])
	}
}
