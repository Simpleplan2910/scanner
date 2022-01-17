package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	PageSize   int
	PageNumber int
	Name       string
}

type ResultStatus string

const (
	Queued     ResultStatus = "Queued"
	InProgress ResultStatus = "In Progress"
	Success    ResultStatus = "Success"
	Failure    ResultStatus = "Failure"
)
