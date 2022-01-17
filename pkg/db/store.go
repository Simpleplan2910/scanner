package db

import "go.mongodb.org/mongo-driver/mongo"

type Store struct {
	Repos  ReposStore
	Result ResultStore
}

func NewStore(db *mongo.Database) (s *Store) {
	s = &Store{
		Repos:  newReposStore(db),
		Result: newResultStore(db),
	}
	return s
}
