package main

import (
	"sort"
	"testing"

	"github.com/alphadose/logstreamer/grpc"
	"github.com/alphadose/logstreamer/mongo"
	"github.com/alphadose/logstreamer/types"
	"github.com/alphadose/logstreamer/utils"
)

// check if both payload are equal or not by checking each individual field
func isEqual(a, b *types.Payload) bool {
	if a.GetName() != b.GetName() || a.GetUserId() != b.GetUserId() || a.GetReviewCount() != b.GetReviewCount() || a.GetUrl() != b.GetUrl() {
		return false
	}
	return true
}

func printData(t *testing.T, grpcData, mongoData []*types.Payload) {
	t.Logf("MongoDB Data\n")
	for idx := range mongoData {
		t.Logf("%#v\n", mongoData[idx])
	}
	t.Logf("\nGRPC Data\n")
	for idx := range grpcData {
		t.Logf("%#v\n", grpcData[idx])
	}
}

// Description: Parse data.txt which contains 7 JSON payloads
// and store it in MongoDB as well as a GRPC server
// then retrieve data from both the MongoDB and GRPC servers separately in the form of golang slices
// sort both of them in a particular order then check if they are equal or not
func doTestRun(t *testing.T, collectionName string) {
	// process file and upload data to both MongoDB and GRPC service
	process(collectionName)

	grpcClient := grpc.NewClient(grpcURI)
	grpcData, err := grpcClient.Consume(1 << 61) // consume all objects present in the GRPC server
	if err != nil {
		t.Fatal(err)
	}

	mongoClient := mongo.NewClient[types.Payload](mongoURI, collectionName)
	mongoData, err := mongoClient.FetchDocs()
	if err != nil {
		t.Fatal(err)
	}

	// if mongoData and grpcData have all equal elements ignoring the order, then this test is successful
	if len(mongoData) != len(grpcData) {
		t.Logf("MongoSize %d : GRPC size %d\n", len(mongoData), len(grpcData))
		printData(t, grpcData, mongoData)
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
			printData(t, grpcData, mongoData)
			t.Fatal("Data retrieved from MongoDB and GRPC sources are inconsistent")
		}
	}
}

// TestSeuquentialFlow tests the application in single-goroutine mode
func TestSeuquentialFlow(t *testing.T) {
	file = "./data.txt"
	mongoURI = "mongodb://localhost:27017"
	grpcURI = "localhost:3002"
	parallel = false
	batchSize = 3
	doTestRun(t, "testseq"+utils.GetTimeStamp())
}

// TestParallelFlow tests the application in multi-goroutine mode
// this is the test for application run with the `-parallel` flag
func TestParallelFlow(t *testing.T) {
	file = "./data.txt"
	mongoURI = "mongodb://localhost:27017"
	grpcURI = "localhost:3002"
	parallel = true
	batchSize = 4
	doTestRun(t, "testpara"+utils.GetTimeStamp())
}
