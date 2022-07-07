package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// projectDatabase is the name of the database used for storing all of alphadose/logstreamer's information
	projectDatabase = "alphadose/logstreamer"

	// primaryKey is the primary key for mongoDB documents
	primaryKey = "_id"

	// timeout is the context timeout for generic operations
	timeout = 20 * time.Second
)

// query operators
const (
	operatorSet               = "$set"
	operatorExists            = "$exists"
	operatorIn                = "$in"
	operatorGreaterThan       = "$gt"
	operatorGreaterThanEquals = "$gte"
	operatorIncrement         = "$inc"
	operatorMatch             = "$match"
	operatorMerge             = "$merge"
	operatorAND               = "$and"
	operatorOR                = "$or"
	operatorNearSphere        = "$nearSphere"
	operatorGeometry          = "$geometry"
	operatorMinDistance       = "$minDistance"
	operatorMaxDistance       = "$maxDistance"
)

// ErrNoDocuments is the error when no matching documents are found
// for an update operation
var ErrNoDocuments = mongo.ErrNoDocuments
