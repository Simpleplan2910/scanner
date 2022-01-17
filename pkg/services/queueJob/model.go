package queuejob

import (
	"io"
	"scanner/pkg/services/git"

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

type findings struct {
	Findings []finding `json:"findings"`
}

type finding struct {
	Type     string   `json:"type"`
	RuleID   string   `json:"ruleId"`
	Location location `json:"location"`
	Metadata metadata `json:"metadata"`
}

type location struct {
	Path     string   `json:"path"`
	Position position `json:"positions"`
}

type position struct {
	Begin line `json:"begin"`
}

type line struct {
	Line int `json:"line"`
}

type metadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}
