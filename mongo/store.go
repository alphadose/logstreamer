package mongo

import (
	"context"

	"github.com/alphadose/logstreamer/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Store represents MongoDB repository storage
type Store struct {
	client *mongo.Client
	coll   *mongo.Collection
}

// NewClient returns a new MongoDB storage
// url format mongodb://mongodb0.example.com:27017
func NewClient(url string) *Store {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		utils.GracefulExit("Mongo-Connection-1", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		utils.GracefulExit("Mongo-Connection-2", err)
	}
	utils.LogInfo("Mongo-Connection-3", "MongoDB Connection Established")
	return &Store{
		client: client,
		coll:   client.Database("tyk").Collection("users"),
	}
}

// Upload uploads all payloads to a MongoDB collection with a callback function attached
// This uses a MongoDB transaction https://www.mongodb.com/docs/manual/core/transactions/
// for ACID transactions in tandem with the callback function provided
// This ensures either both MongoDB and the callback succeed or they both fail
// the callback in this case is `upload to GRPC server` as per the problem statement
func (s *Store) Upload(callback func() error, payloads ...any) error {
	ctx := context.Background()
	session, err := s.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Majority()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	transaction := func(sc mongo.SessionContext) (any, error) {
		results, err := s.coll.InsertMany(sc, payloads)
		if err != nil {
			return nil, err
		}
		if err := callback(); err != nil {
			return nil, err
		}
		return results, nil
	}
	_, err = session.WithTransaction(ctx, transaction, txnOpts)
	return err
}
