package repos

import (
	"scanner/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReqScan struct {
	ReposId primitive.ObjectID
	Substr  string
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

type ReqGetResult struct {
	ReposID              primitive.ObjectID
	PageSize, PageNumber int
}

type RespGetResult struct {
	Code    int
	Message string
	Total   int64
	Results []db.Result
}
