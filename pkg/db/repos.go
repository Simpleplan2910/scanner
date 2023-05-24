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

const reposCollection = "repository"

type reposStore struct {
	collection     *mongo.Collection
	firstWriteDone bool
}

type Repos struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ReposURL  string             `bson:"reposURL"`
	IsArchive bool               `bson:"isArchive"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

type ReposStore interface {
	Add(ctx context.Context, v *Repos) (id primitive.ObjectID, err error)
	Filter(ctx context.Context, filter *FilterRepos) (repos []Repos, total int64, err error)
	Archive(ctx context.Context, id primitive.ObjectID) error
	Update(ctx context.Context, id primitive.ObjectID, update *UpdateRepos) error
	Get(ctx context.Context, id primitive.ObjectID) (repos *Repos, err error)
}

func newReposStore(db *mongo.Database) ReposStore {
	return &reposStore{
		collection: db.Collection(reposCollection),
	}
}

func (db *reposStore) Add(ctx context.Context, v *Repos) (id primitive.ObjectID, err error) {
	if !db.firstWriteDone {
		if err = db.createIndexes(ctx); err != nil {
			return id, err
		}
	}
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
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

func (db *reposStore) Filter(ctx context.Context, filter *FilterRepos) (repos []Repos, total int64, err error) {
	findOptions := options.Find()
	rep := []Repos{}
	match := bson.M{}
	if filter.PageNumber >= 1 && filter.PageSize >= 1 {
		findOptions.SetSkip(int64((filter.PageNumber - 1) * filter.PageSize))
		findOptions.SetLimit(int64(filter.PageSize))
	}
	// create time descending
	findOptions.SetSort(bson.M{"createdAt": -1})

	// filter docs that name contains filter text
	if filter.Name != "" {
		match["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}

	cursor, err := db.collection.Find(ctx, match, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		temp := Repos{}
		err = cursor.Decode(&temp)
		if err != nil {
			return nil, 0, err
		}
		rep = append(rep, temp)
	}
	total, err = db.collection.CountDocuments(ctx, match)
	if err != nil {
		return nil, 0, err
	}
	return rep, total, nil
}

func (db *reposStore) Get(ctx context.Context, id primitive.ObjectID) (repos *Repos, err error) {
	repos = &Repos{}
	filter := bson.M{
		"_id": id,
	}
	err = db.collection.FindOne(ctx, filter).Decode(repos)
	return repos, err
}

func (db *reposStore) Archive(ctx context.Context, id primitive.ObjectID) error {
	up := bson.M{
		"isArchive": true,
	}
	_, err := db.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": up})
	return err
}

func (db *reposStore) Update(ctx context.Context, id primitive.ObjectID, update *UpdateRepos) error {
	up := update.toBson()
	_, err := db.collection.UpdateOne(ctx, bson.M{"_id": id}, up)
	return err
}

func (db *reposStore) createIndexes(ctx context.Context) error {
	iv := db.collection.Indexes()
	model := mongo.IndexModel{
		Keys: bson.D{
			{"createdAt", int32(1)},
		},
	}
	if err := createIndexIfNotExists(ctx, iv, model); err != nil {
		return err
	}
	db.firstWriteDone = true
	return nil
}
