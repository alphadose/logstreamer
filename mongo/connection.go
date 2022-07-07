package mongo

import (
	"context"
	"time"

	"github.com/alphadose/logstreamer/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client = getClient()
var db = client.Database(projectDatabase)
var url = ""

func getClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		utils.GracefulExit("Mongo-Connection-1", err)
	}
	return client
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		utils.GracefulExit("Mongo-Connection-4", err)
	} else {
		utils.LogInfo("Mongo-Connection-2", "MongoDB Connection Established")
	}
}
