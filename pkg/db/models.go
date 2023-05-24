package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateRepos struct {
	Name     string
	ReposURL string
}

func (u *UpdateRepos) toBson() bson.M {
	update := bson.M{}
	if u.Name != "" {
		update["name"] = u.Name
	}
	if u.ReposURL != "" {
		update["reposURL"] = u.ReposURL
	}
	update["updatedAt"] = time.Now()
	return update
}

type FilterRepos struct {
	PageSize   int
	PageNumber int
	Name       string
}

type FilterResult struct {
	ReposID    primitive.ObjectID
	PageSize   int
	PageNumber int
	Name       string
}

type ScanStatus string

const (
	Queued     ScanStatus = "Queued"
	InProgress ScanStatus = "In Progress"
	Success    ScanStatus = "Success"
	Failure    ScanStatus = "Failure"
)

type Line struct {
	LineNum uint32
	Indexes []int
}
