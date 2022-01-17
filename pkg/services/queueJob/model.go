package queuejob

import (
	"io"
	"test_guard/pkg/services/git"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ReposId   primitive.ObjectID
	ReposURL  string
	ReposName string
}

type repos struct {
	ResultId primitive.ObjectID
	Job      *Job
	Repos    git.Repos
}

type scanResult struct {
	IsContainVulnerable bool
	Line                []int
	Filename            string
	Error               error
}

type singleFile struct {
	Reader   io.Reader
	FileName string
}
