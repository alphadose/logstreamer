package grpc

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/alphadose/logstreamer/types"
	"github.com/alphadose/logstreamer/utils"

	"google.golang.org/grpc"
)

type server struct {
	types.UnimplementedBrokerServer
}

// ListenAndServe binds the GRPC server to the port and starts listening
func ListenAndServe(port uint64) error {
	utils.LogInfo("GRPC-Serve-1", "Starting GRPC server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		utils.GracefulExit("GRPC-Serve-2", err)
	}
	utils.LogInfo("GRPC-Serve-3", "Listening on port %d", port)
	s := grpc.NewServer()
	types.RegisterBrokerServer(s, &server{})
	return s.Serve(lis)
}

// Publish receives data stream from multiple clients and stores them in-memory
func (*server) Publish(stream types.Broker_PublishServer) error {
	utils.LogInfo("GRPC-Publish-1", "New Connection Received")
	var (
		req         *types.Payload
		err         error
		dataStaging = make([]*types.Payload, 0) // staging storage which is only committed after sucessful operation
	)
	for {
		req, err = stream.Recv()
		// Check if the stream has finished
		if err == io.EOF {
			// reached end of stream, commit stating data into main storage and return success response
			// this ensures both GRPC upload and MongoDB InsertMany transaction is entirely
			// atomic in nature and ACID compliant
			for idx := range dataStaging {
				store.Enqueue(dataStaging[idx])
			}
			return stream.SendAndClose(&types.Response{Success: true})
		}
		if err != nil {
			utils.LogError("GRPC-Publish-2", err)
			return stream.SendAndClose(&types.Response{Success: false, Error: err.Error()})
		}
		// Store payload to staging storage
		dataStaging = append(dataStaging, req)
	}
}

var errInvalidCountParameter = errors.New("Invalid count parameter specified")

// Consume consumes data objects stored in the GRPC server via a streaming response
// NOTE:- each payload can be consumed only once
func (*server) Consume(req *types.ConsumeRequest, stream types.Broker_ConsumeServer) error {
	utils.LogInfo("GRPC-Consume-1", "New Connection Received")
	if req.GetCount() <= 0 {
		return errInvalidCountParameter
	}
	var data *types.Payload
	// try dequeue till given count, if count > len(queue) then len(queue) items will be returned
	for ctr := req.GetCount(); ctr > 0; ctr-- {
		data = store.Dequeue()
		if data != nil {
			if err := stream.Send(data); err != nil {
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}
