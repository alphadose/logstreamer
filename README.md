# logstreamer
[![Main Actions Status](https://github.com/alphadose/logstreamer/workflows/Go/badge.svg)](https://github.com/alphadose/logstreamer/actions)
> A CLI tool for streaming logs to different kinds of storage layers such as MongoDB and GRPC service

## Usage

You need Golang [1.18.x](https://go.dev/dl/) or above since this package uses generics

```bash
# Build the binary
$ go build -o lsr

# Get Help and checkout the flags
$ ./lsr --help
Usage of ./lsr:
  -bsize uint
    	Batch size of upload operations (restriction helpful in cases of file_size > 16 GB) (default 200)
  -f string
    	Absolute path to file (default "./data.txt")
  -grpc string
    	URI of GRPC server (default "localhost:3002")
  -mongo string
    	MongoDB URI (default "mongodb://localhost:27017")
  -parallel
    	Should storage upload operations run in parallel?

# Start the GRPC server for testing
$ go run cmd/grpc_server/main.go &
[INFO] 7-7-2022 21:45:37 -> Starting GRPC server
[INFO] 7-7-2022 21:45:37 -> Listening on port :3002

# Use the binary
$ ./lsr -f data_large.txt -grpc "localhost:3002" -mongo "mongodb://localhost:27017"
[INFO] 7-7-2022 21:47:54 -> MongoDB Connection Established
[INFO] 7-7-2022 21:47:54 -> Starting Operation
[INFO] 7-7-2022 21:47:54 -> Processing Batch: 1
[INFO] 7-7-2022 21:47:54 -> Processing Batch: 2
[INFO] 7-7-2022 21:47:54 -> Processing Batch: 3
[INFO] 7-7-2022 21:47:54 -> Processing Batch: 4
[INFO] 7-7-2022 21:47:54 -> Successfully Completed
```

## Testing

```bash
$ go test e2e_test.go core.go
ok  	command-line-arguments	0.290s
```

## Author

Anish Mukherjee
