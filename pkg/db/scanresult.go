package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const resultCollection = "results"

type Result struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ScanID    primitive.ObjectID `bson:"scanId"`
	Lines     []Line             `bson:"lines"`
	Filename  string             `bson:"fileName"`
	Error     string             `bson:"error"`
	CreatedAt time.Time          `bson:"createdAt"`
}

type resultStore struct {
	collection     *mongo.Collection
	firstWriteDone bool
}

type ResultStore interface {
	Add(ctx context.Context, v *Result) (id primitive.ObjectID, err error)
	Filter(ctx context.Context, filter *FilterResult) (results []Result, total int64, err error)
}

func newResultStore(db *mongo.Database) ResultStore {
	return &resultStore{
		collection: db.Collection(resultCollection),
	}
}

func (db *resultStore) Add(ctx context.Context, v *Result) (id primitive.ObjectID, err error) {
	if !db.firstWriteDone {
		if err = db.createIndexes(ctx); err != nil {
			return id, err
		}
	}
	result, err := db.collection.InsertOne(ctx, v)
	if err != nil {
		return id, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return id, fmt.Errorf("invalid object id return")
	}
	return oid, nil
}

func (db *resultStore) Filter(ctx context.Context, filter *FilterResult) (results []Result, total int64, err error) {
	findOptions := options.Find()
	results = []Result{}
	match := bson.M{}
	if filter.PageNumber >= 1 && filter.PageSize >= 1 {
		findOptions.SetSkip(int64((filter.PageNumber - 1) * filter.PageSize))
		findOptions.SetLimit(int64(filter.PageSize))
	}
	if !filter.ReposID.IsZero() {
		match["scanId"] = filter.ReposID
	}

	findOptions.SetSort(bson.M{"createdAt": -1})

	cursor, err := db.collection.Find(ctx, match, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		temp := Result{}
		err = cursor.Decode(&temp)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, temp)
	}
	total, err = db.collection.CountDocuments(ctx, match)
	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (db *resultStore) createIndexes(ctx context.Context) error {
	iv := db.collection.Indexes()
	model := mongo.IndexModel{
		Keys: bson.D{
			{"CreatedAt", int32(1)},
		},
	}
	if err := createIndexIfNotExists(ctx, iv, model); err != nil {
		return err
	}

	iv = db.collection.Indexes()
	model = mongo.IndexModel{
		Keys: bson.D{
			{"scanId", int32(1)},
		},
	}
	if err := createIndexIfNotExists(ctx, iv, model); err != nil {
		return err
	}
	db.firstWriteDone = true
	return nil
}
