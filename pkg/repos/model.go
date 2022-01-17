package repos

import (
	"test_guard/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReqScan struct {
	ReposId primitive.ObjectID
}

type RespScan struct {
	Code    int
	Message string
}

type ReqAddRepos struct {
	Name     string
	ReposURL string
}

type RespAddRepos struct {
	Code    int
	Message string
	ID      primitive.ObjectID
}

type ReqGetRepos struct {
	PageSize, PageNumber int
}

type RespGetRepos struct {
	Code    int
	Message string
	Repos   []db.Repos
	Total   int64
}

type ReqUpdateRepos struct {
	ID       primitive.ObjectID
	Name     string
	ReposURL string
}

type RespUpdateRepos struct {
	Code    int
	Message string
}

type ReqDeleteRepos struct {
	ID primitive.ObjectID
}

type RespDeleteRepos struct {
	Code    int
	Message string
}
