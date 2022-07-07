package grpc

import (
	"context"
	"errors"
	"io"

	"github.com/alphadose/logstreamer/types"
	"github.com/alphadose/logstreamer/utils"

	"google.golang.org/grpc"
)

// Client is the GRPC client struct
type Client struct {
	conn types.BrokerClient
}

// NewClient returns a new GRPC client give its URL IP:PORT or DNS:PORT
func NewClient(url string) *Client {
	opts := grpc.WithInsecure()
	c, err := grpc.Dial(url, opts)
	if err != nil {
		utils.GracefulExit("GRPC-Client-1", err)
	}
	return &Client{conn: types.NewBrokerClient(c)}
}

// Publish data to the GRPC server
func (c *Client) Publish(payloads []*types.Payload) error {
	stream, err := c.conn.Publish(context.Background())
	if err != nil {
		return err
	}

	// Iterate over the request message
	for _, p := range payloads {
		if err := stream.Send(p); err != nil {
			return err
		}
	}

	// Once the for loop finishes, the stream is closed
	// and get the response and a potential error
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	// Unsuccessful response
	if !res.GetSuccess() {
		return errors.New(res.GetError())
	}
	return nil
}

// Consume data from the GRPC server
func (c *Client) Consume(count int64) ([]*types.Payload, error) {
	if count <= 0 {
		return nil, errInvalidCountParameter
	}
	stream, err := c.conn.Consume(context.Background(), &types.ConsumeRequest{Count: count})
	if err != nil {
		return nil, err
	}
	var (
		data = make([]*types.Payload, 0)
		tmp  *types.Payload
	)
	for {
		tmp, err = stream.Recv()
		if err == io.EOF {
			return data, nil
		}
		if err != nil {
			return data, err
		}
		data = append(data, tmp)
	}
}
