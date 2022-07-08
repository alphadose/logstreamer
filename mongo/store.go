package mongo

import (
	"context"

	"github.com/alphadose/logstreamer/types"
	"github.com/alphadose/logstreamer/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Store represents MongoDB repository storage
type Store[T any] struct {
	client *mongo.Client
	coll   *mongo.Collection
}

// NewClient returns a new MongoDB storage
// url format mongodb://mongodb0.example.com:27017
func NewClient[T any](url string, collectionName string) *Store[T] {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		utils.GracefulExit("Mongo-Connection-1", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		utils.GracefulExit("Mongo-Connection-2", err)
	}
	utils.LogInfo("Mongo-Connection-3", "MongoDB Connection Established")
	return &Store[T]{
		client: client,
		coll:   client.Database("tyk").Collection(collectionName),
	}
}

// Upload uploads all payloads to a MongoDB collection with a callback function attached
// This uses a MongoDB transaction https://www.mongodb.com/docs/manual/core/transactions/
// for ACID transactions in tandem with the callback function provided
// This ensures either both MongoDB and the callback succeed or they both fail
// the callback in this case is `upload to GRPC server` as per the problem statement
// Note:- MongoDB transaction feature requires replica sets https://www.mongodb.com/docs/manual/replication/
func (s *Store[T]) Upload(callback func() error, payloads []*T) error {
	ctx := context.Background()
	session, err := s.client.StartSession()
	if err != nil {
		return err
	}

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Majority()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	// convert data to []any in order to make it compatible with MongoDB function arguments
	uploadData := make([]any, len(payloads), len(payloads))
	for idx := range payloads {
		uploadData[idx] = payloads[idx]
	}

	transaction := func(sc mongo.SessionContext) (any, error) {
		results, err := s.coll.InsertMany(sc, uploadData)
		if err != nil {
			return nil, err
		}
		// Ensure either both the insert many operation and the callback (GRPC upload) succeed
		// or else this entire transaction fails
		if err := callback(); err != nil {
			return nil, err
		}
		return results, nil
	}
	_, err = session.WithTransaction(ctx, transaction, txnOpts)
	return err
}

// FetchDocs is a function which returns all documents present the MongoDB store
// used for testing purposes in `e2e_test.go`
func (s *Store[T]) FetchDocs() ([]*T, error) {
	ctx := context.Background()
	cursor, err := s.coll.Find(ctx, types.M{})
	if err != nil {
		return nil, err
	}
	data := make([]*T, cursor.RemainingBatchLength(), cursor.RemainingBatchLength())
	for index := 0; cursor.Next(ctx) && index < cap(data); index++ {
		data[index] = new(T)
		if err := cursor.Decode(data[index]); err != nil {
			return nil, err
		}
	}
	return data, nil
}
