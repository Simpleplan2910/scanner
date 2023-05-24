package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Scan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	RepoID      primitive.ObjectID `bson:"repoId"`
	Status      ScanStatus         `bson:"status"`
	StartedAt   time.Time          `bson:"startedAt"`
	CompletedAt time.Time          `bson:"completedAt"`
}

type scanStore struct {
	collection     *mongo.Collection
	firstWriteDone bool
}

type ScanStore interface {
	Add(ctx context.Context, v *Scan) (id primitive.ObjectID, err error)
	UpdateCompleted(ctx context.Context, id primitive.ObjectID) error
	UpdateError(ctx context.Context, id primitive.ObjectID) error
}

func newScanStore(db *mongo.Database) ScanStore {
	return &scanStore{
		collection: db.Collection(reposCollection),
	}
}

func (s *scanStore) Add(ctx context.Context, v *Scan) (id primitive.ObjectID, err error) {
	if !s.firstWriteDone {
		if err = s.createIndexes(ctx); err != nil {
			return id, err
		}
	}
	result, err := s.collection.InsertOne(ctx, v)
	if err != nil {
		return id, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return id, fmt.Errorf("invalid object id return")
	}
	return oid, nil
}

func (s *scanStore) UpdateCompleted(ctx context.Context, id primitive.ObjectID) error {
	up := bson.M{
		"completedAt": time.Now(),
	}
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": up})
	return err
}

func (s *scanStore) UpdateError(ctx context.Context, id primitive.ObjectID) error {
	up := bson.M{
		"completedAt": time.Now(),
		"status":      Failure,
	}
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": up})
	return err
}

func (s *scanStore) createIndexes(ctx context.Context) error {
	iv := s.collection.Indexes()
	model := mongo.IndexModel{
		Keys: bson.D{
			{"repoId", int32(1)},
		},
	}
	if err := createIndexIfNotExists(ctx, iv, model); err != nil {
		return err
	}

	model = mongo.IndexModel{
		Keys: bson.D{
			{"startedAt", int32(1)},
		},
	}
	if err := createIndexIfNotExists(ctx, iv, model); err != nil {
		return err
	}
	s.firstWriteDone = true
	return nil
}
