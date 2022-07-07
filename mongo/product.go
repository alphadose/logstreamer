package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	// productCollectionKey is the collection for all products
	productCollectionKey = "products"
	// Key containing deleted products
	deletedProductsCollectionKey = "deleted_products"
	nameKey                      = "name"
	hubNameKey                   = "hub_name"
	hubIDKey                     = "hub_id"
	imageURLKey                  = "image_url"
	priceKey                     = "price"
	quantityKey                  = "quantity"
)

// The link to the product collection
var productCollection = db.Collection(productCollectionKey)

// TransferProductsToCustomer shifts the product inventories to the customer depending on his/her order
func TransferProductsToCustomer(productList map[string]float64) error {
	ctx := context.Background()
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	models := make([]mongo.WriteModel, 0)
	for productID, quantity := range productList {
		productDocID, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			return err
		}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(bson.D{
			bson.E{Key: primaryKey, Value: productDocID},
			bson.E{Key: quantityKey, Value: bson.D{bson.E{Key: operatorGreaterThanEquals, Value: quantity}}},
		}).SetUpdate(bson.D{bson.E{Key: operatorIncrement, Value: bson.D{bson.E{Key: quantityKey, Value: -1 * quantity}}}}))
	}
	opts := options.BulkWrite().SetOrdered(true)

	callback := func(sc mongo.SessionContext) (any, error) {
		results, err := productCollection.BulkWrite(sc, models, opts)
		if err != nil {
			return nil, err
		}
		if results.MatchedCount != int64(len(productList)) {
			// Happens when a product runs out of stock but the customer has paid for it
			// Inventory transfer only occurs after payments and depends on razorpay's webhook latency
			// Race condition occurs when another payment was processed faster than this payment
			// as a result there was a shortage of inventory with respect to the given product list
			// In this case we need to refund the customer and rollback the update
			// Razorpay refund order_id
			// Also SMS the customer that his order has been refunded
			return nil, nil
		}
		return results, nil
	}

	_, err = session.WithTransaction(ctx, callback, txnOpts)
	return err
}
