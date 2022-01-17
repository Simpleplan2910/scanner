package repos

import (
	"context"
	"fmt"
	apierror "scanner/pkg/apiError"
	"scanner/pkg/db"
	"scanner/pkg/services/git"
	queuejob "scanner/pkg/services/queueJob"
	"time"
)

type Service interface {
	AddRepos(ctx context.Context, req *ReqAddRepos) (resp *RespAddRepos, err error)
	GetRepos(ctx context.Context, req *ReqGetRepos) (resp *RespGetRepos, err error)
	UpdateRepos(ctx context.Context, req *ReqUpdateRepos) (resp *RespUpdateRepos, err error)
	DeleteRepos(ctx context.Context, req *ReqDeleteRepos) (resp *RespDeleteRepos, err error)
	StartScanRepos(ctx context.Context, req *ReqScan) (resp *RespScan, err error)
}

type service struct {
	gitService git.Service
	repoStore  db.ReposStore
	queue      queuejob.QueueJob
}

func NewService(gitService git.Service, repoStore db.ReposStore, queue queuejob.QueueJob) Service {
	return &service{
		gitService: gitService,
		repoStore:  repoStore,
		queue:      queue,
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
	j := &queuejob.Job{
		ReposId:   repos.ID,
		ReposURL:  repos.ReposURL,
		ReposName: repos.Name,
	}
	// add repos to job queue
	err = s.queue.AddJob(ctx, j)
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
	// TODO: should check if reposURL valid  and exist
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

func (s *service) DeleteRepos(ctx context.Context, req *ReqDeleteRepos) (resp *RespDeleteRepos, err error) {
	return
}
