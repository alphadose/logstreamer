package mongo

import (
	"context"

	"github.com/alphadose/logstreamer/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create Operations

// insertOne inserts a document into a mongoDB collection
func insertOne(collection *mongo.Collection, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

// insertMany inserts multiple document into a mongoDB collection
func insertMany(collection *mongo.Collection, data []interface{}) ([]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}

// Read Operations

// fetchDocs is a generic function which takes a collection name and mongoDB filter as input and returns documents
func fetchDocs[T any](collection *mongo.Collection, filter types.M, opts ...*options.FindOptions) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		defer cursor.Close(ctx)
		return nil, err
	}
	data := make([]T, cursor.RemainingBatchLength(), cursor.RemainingBatchLength())
	for index := 0; cursor.Next(ctx) && index < cap(data); index++ {
		if err := cursor.Decode(&data[index]); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// fetchOne returns a document corresponding to the type given a filter
func fetchOne[T any](collection *mongo.Collection, filter interface{}) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	data := new(T)
	err := collection.FindOne(ctx, filter).Decode(data)
	return data, err
}

// countDocs returns the number of documents matching a filter
func countDocs(collection *mongo.Collection, filter types.M) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return collection.CountDocuments(ctx, filter)
}

// Update Operations

// updateOne updates a document in the mongoDB collection
func updateOne(collection *mongo.Collection, filter types.M, modifier interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return collection.FindOneAndUpdate(ctx, filter, modifier, opts...).Err()
}

// bulkUpsert upserts multiple documents using BulkWrite
func bulkUpsert(collection *mongo.Collection, data []mongo.WriteModel, opts ...*options.BulkWriteOptions) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return collection.BulkWrite(ctx, data, opts...)
}

// updateMany updates multiple documents in the mongoDB collection
func updateMany(collection *mongo.Collection, filter types.M, data interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return collection.UpdateMany(ctx, filter, types.M{operatorSet: data}, nil)
}

// Delete Operations

// deleteOne deletes a document from a mongoDB collection
func deleteOne(collection *mongo.Collection, filter types.M) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return collection.DeleteOne(ctx, filter)
}

// deleteMany deletes multiple documents from a mongoDB collection
func deleteMany(collection *mongo.Collection, filter types.M) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return collection.DeleteMany(ctx, filter)
}
