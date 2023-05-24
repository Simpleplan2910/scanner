package repos

import (
	"context"
	"fmt"
	apierror "scanner/pkg/apiError"
	"scanner/pkg/db"
	"scanner/pkg/services/git"
	"scanner/pkg/services/scanner"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	AddRepos(ctx context.Context, req *ReqAddRepos) (resp *RespAddRepos, err error)
	GetRepos(ctx context.Context, req *ReqGetRepos) (resp *RespGetRepos, err error)
	UpdateRepos(ctx context.Context, req *ReqUpdateRepos) (resp *RespUpdateRepos, err error)
	ArchiveRepos(ctx context.Context, req *ReqDeleteRepos) (resp *RespDeleteRepos, err error)
	StartScanRepos(ctx context.Context, req *ReqScan) (resp *RespScan, err error)
	GetResult(ctx context.Context, req *ReqGetResult) (resp *RespGetResult, err error)
}

type service struct {
	gitService     git.Service
	repoStore      db.ReposStore
	result         db.ResultStore
	scannerService scanner.Service
}

func NewService(gitService git.Service, repoStore db.ReposStore, result db.ResultStore, s scanner.Service) Service {
	return &service{
		gitService:     gitService,
		repoStore:      repoStore,
		result:         result,
		scannerService: s,
	}
}

func (s *service) StartScanRepos(ctx context.Context, req *ReqScan) (resp *RespScan, err error) {
	// get repos
	resp = &RespScan{}
	repos, err := s.repoStore.Get(ctx, req.ReposId)
	if err != nil {
		resp.Code = apierror.InternalServerError
		resp.Message = apierror.InternalServerErrorMess
		return resp, fmt.Errorf("can't get respos from db with error: %w", err)
	}
	scanId := primitive.NewObjectID()
	repo := &scanner.Repos{
		ScanId:    scanId,
		ReposId:   repos.ID,
		ReposURL:  repos.ReposURL,
		ReposName: repos.Name,
	}
	// add repos to job queue
	err = s.scannerService.Scan(ctx, repo, req.Substr)
	if err != nil {
		resp.Code = apierror.InternalServerError
		resp.Message = apierror.InternalServerErrorMess
		return resp, fmt.Errorf("filter error: %w", err)
	}

	return resp, nil
}

func (s *service) AddRepos(ctx context.Context, req *ReqAddRepos) (resp *RespAddRepos, err error) {
	resp = &RespAddRepos{}
	// TODO: should check if reposURL valid and exist
	if req.Name == "" || req.ReposURL == "" {
		resp.Code = apierror.InvalidRequest
		resp.Message = apierror.InvalidRequestMess
		return resp, fmt.Errorf("empty name or reposURL when add")
	}
	repos := &db.Repos{
		Name:      req.Name,
		ReposURL:  req.ReposURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := s.repoStore.Add(ctx, repos)
	if err != nil {
		resp.Code = apierror.InternalServerError
		resp.Message = apierror.InternalServerErrorMess
		return resp, fmt.Errorf("add to db error: %w", err)
	}
	resp.ID = id
	return resp, nil
}

func (s *service) GetRepos(ctx context.Context, req *ReqGetRepos) (resp *RespGetRepos, err error) {
	resp = &RespGetRepos{}
	if req.PageNumber < 1 || req.PageSize < 1 {
		resp.Code = apierror.InvalidRequest
		resp.Message = apierror.InvalidRequestMess
		return resp, fmt.Errorf("page number or page size less than 1")
	}
	filter := &db.FilterRepos{
		PageSize:   req.PageSize,
		PageNumber: req.PageNumber,
	}
	repos, total, err := s.repoStore.Filter(ctx, filter)
	if err != nil {
		resp.Code = apierror.InternalServerError
		resp.Message = apierror.InternalServerErrorMess
		return resp, fmt.Errorf("filter error: %w", err)
	}
	resp.Repos = repos
	resp.Total = total
	return resp, nil

}

func (s *service) UpdateRepos(ctx context.Context, req *ReqUpdateRepos) (resp *RespUpdateRepos, err error) {
	resp = &RespUpdateRepos{}
	// TODO: should check if reposURL valid and exist
	if req.ID.IsZero() {
		resp.Code = apierror.InvalidRequest
		resp.Message = apierror.InvalidRequestMess
		return resp, fmt.Errorf("id empty")
	}

	if req.Name == "" && req.ReposURL == "" {
		resp.Code = apierror.InvalidRequest
		resp.Message = apierror.InvalidRequestMess
		return resp, fmt.Errorf("name and reposURL empty")
	}
	update := &db.UpdateRepos{
		Name:     req.Name,
		ReposURL: req.ReposURL,
	}
	err = s.repoStore.Update(ctx, req.ID, update)
	if err != nil {
		resp.Code = apierror.InternalServerError
		resp.Message = apierror.InternalServerErrorMess
		return resp, fmt.Errorf("update to db error: %w", err)
	}

	return resp, nil
}

func (s *service) ArchiveRepos(ctx context.Context, req *ReqDeleteRepos) (resp *RespDeleteRepos, err error) {
	resp = &RespDeleteRepos{}
	// TODO: should check if reposURL valid and exist
	if req.ID.IsZero() {
		resp.Code = apierror.InvalidRequest
		resp.Message = apierror.InvalidRequestMess
		return resp, fmt.Errorf("id empty")
	}

	err = s.repoStore.Archive(ctx, req.ID)
	if err != nil {
		resp.Code = apierror.InternalServerError
		resp.Message = apierror.InternalServerErrorMess
		return resp, fmt.Errorf("delete error: %w", err)
	}

	return resp, nil
}

func (s *service) GetResult(ctx context.Context, req *ReqGetResult) (resp *RespGetResult, err error) {
	// resp = &RespGetResult{}
	// if req.ReposID.IsZero() {
	// 	resp.Code = apierror.InvalidRequest
	// 	resp.Message = apierror.InvalidRequestMess
	// 	return resp, fmt.Errorf("id empty")
	// }

	// if req.PageNumber < 1 || req.PageSize < 1 {
	// 	resp.Code = apierror.InvalidRequest
	// 	resp.Message = apierror.InvalidRequestMess
	// 	return resp, fmt.Errorf("page number or page size less than 1")
	// }

	// filter := &db.FilterResult{
	// 	ReposID:    req.ReposID,
	// 	PageSize:   req.PageSize,
	// 	PageNumber: req.PageNumber,
	// }
	// results, total, err := s.result.Filter(ctx, filter)
	// if err != nil {
	// 	resp.Code = apierror.InternalServerError
	// 	resp.Message = apierror.InternalServerErrorMess
	// 	return resp, fmt.Errorf("filter error: %w", err)
	// }
	// resp.Results = results
	// resp.Total = total
	return resp, nil
}
