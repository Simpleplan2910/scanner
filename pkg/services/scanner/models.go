package scanner

import (
	"scanner/pkg/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repos struct {
	ReposId   primitive.ObjectID
	ScanId    primitive.ObjectID
	ReposURL  string
	ReposName string
}

type scanResult struct {
	Lines    []db.Line
	Filename string
	Error    error
}

type file struct {
	File     []byte
	FileName string
}
