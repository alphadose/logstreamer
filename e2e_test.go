package main

import (
	"testing"

	"github.com/alphadose/logstreamer/grpc"
	"github.com/alphadose/logstreamer/mongo"
	"github.com/alphadose/logstreamer/types"
)

func isEqual(a, b *types.Payload) bool {
	if a.GetName() != b.GetName() || a.GetUserId() != b.GetUserId() || a.GetReviewCount() != b.GetReviewCount() || a.GetUrl() != b.GetUrl() {
		return false
	}
	return true
}

func contains(arr []*types.Payload, elem *types.Payload) bool {
	for idx := range arr {
		if isEqual(arr[idx], elem) {
			return true
		}
	}
	return false
}

func TestEnd2End(t *testing.T) {
	file = "./data.txt"
	mongoURI = "mongodb://localhost:27017"
	grpcURI = "localhost" + grpc.Port
	parallel = false
	batchSize = 200

	// process file and upload data to both MongoDB and GRPC service
	process()

	grpcClient := grpc.NewClient(grpcURI)

	grpcData, err := grpcClient.Consume(5) // number of payloads in data.txt = 5
	if err != nil {
		t.Fatal(err)
	}

	mongoClient := mongo.NewClient(mongoURI)

	mongoData, err := mongoClient.FetchDocs()
	if err != nil {
		t.Fatal(err)
	}

	// if mongoData and grpcData have all equal elements ignoring the order, then this test is successful
	if len(mongoData) != len(grpcData) {
		t.Fatal("Unequal sizes of array received from MongoDB and GRPC service")
	}

	for idx := range mongoData {
		if !contains(grpcData, mongoData[idx]) {
			t.Fatal("Data retrieved from both MongoDB and GRPC sources are inconsistent")
		}
	}
}
