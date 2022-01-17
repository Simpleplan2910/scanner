package db

import (
	"bytes"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create an index if it doesn't already exist
func createIndexIfNotExists(ctx context.Context, iv mongo.IndexView, model mongo.IndexModel) error {
	c, err := iv.List(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = c.Close(ctx)
	}()

	var found bool
	for c.Next(ctx) {
		keyElem, err := c.Current.LookupErr("key")
		if err != nil {
			return err
		}

		keyElemDoc := keyElem.Document()
		modelKeysDoc, err := bson.Marshal(model.Keys)
		if err != nil {
			return err
		}

		if bytes.Equal(modelKeysDoc, keyElemDoc) {
			found = true
			break
		}
	}

	if !found {
		_, err = iv.CreateOne(ctx, model)
		if err != nil {
			return err
		}
	}

	return nil
}
