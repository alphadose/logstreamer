package main

import (
	"sort"
	"testing"

	"github.com/alphadose/logstreamer/grpc"
	"github.com/alphadose/logstreamer/mongo"
	"github.com/alphadose/logstreamer/types"
	"github.com/alphadose/logstreamer/utils"
)

func isEqual(a, b *types.Payload) bool {
	if a.GetName() != b.GetName() || a.GetUserId() != b.GetUserId() || a.GetReviewCount() != b.GetReviewCount() || a.GetUrl() != b.GetUrl() {
		return false
	}
	return true
}

// Description: Parse data.txt which contains 5 JSON payloads
// and store it in MongoDB as well as a GRPC server
// then retrieve data from both the MongoDB and GRPC servers separately in the form of golang slices
// sort both of them in a particular order then check if they are equal or not
func TestEnd2End(t *testing.T) {
	file = "./data.txt"
	mongoURI = "mongodb://localhost:27017"
	grpcURI = "localhost" + grpc.Port
	parallel = false
	batchSize = 200

	var collectionName = "test" + utils.GetTimeStamp()

	// process file and upload data to both MongoDB and GRPC service
	process(collectionName)

	grpcClient := grpc.NewClient(grpcURI)

	grpcData, err := grpcClient.Consume(5) // number of payloads in data.txt = 5
	if err != nil {
		t.Fatal(err)
	}

	mongoClient := mongo.NewClient(mongoURI, collectionName)

	mongoData, err := mongoClient.FetchDocs()
	if err != nil {
		t.Fatal(err)
	}

	// if mongoData and grpcData have all equal elements ignoring the order, then this test is successful
	if len(mongoData) != len(grpcData) {
		t.Fatal("Unequal sizes of array received from MongoDB and GRPC service")
	}

	// Sort both slices in a particular order for checking if both are equal or not
	// In this case they are sorted with respect to review_count
	// all review_count fields have unique values in the file `data.txt` so there are no edge cases in this test
	sort.Slice(grpcData, func(i, j int) bool {
		return grpcData[i].GetReviewCount() < grpcData[j].GetReviewCount()
	})
	sort.Slice(mongoData, func(i, j int) bool {
		return mongoData[i].GetReviewCount() < mongoData[j].GetReviewCount()
	})

	// check if both slices are equal
	for idx := range mongoData {
		if !isEqual(mongoData[idx], grpcData[idx]) {
			t.Fatal("Data retrieved from MongoDB and GRPC sources are inconsistent")
		}
	}
}
