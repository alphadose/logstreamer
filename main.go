/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import "flag"

func main() {
	var (
		file      string
		batchSize uint64
		parallel  bool
	)
	flag.StringVar(&file, "f", "./data.txt", "Absolute path to the file")
	flag.Uint64Var(&batchSize, "f", 200, "Batch size of upload operations (restriction helpful in cases of file_size > 16 GB)")
	flag.BoolVar(&parallel, "p", false, "Should storage upload operations run in parallel?")
	flag.Parse()
}
