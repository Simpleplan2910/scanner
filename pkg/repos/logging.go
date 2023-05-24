package repos

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingService struct {
	logger *logrus.Entry
	next   Service
}

func NewLoggingService(logger *logrus.Entry, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) AddRepos(ctx context.Context, req *ReqAddRepos) (resp *RespAddRepos, err error) {
	defer func(begin time.Time) {
		took := time.Since(begin)
		s.logger.WithFields(logrus.Fields{
			"method":           "AddRepos",
			"request.ReposURL": req.ReposURL,
			"request.Name":     req.Name,
			"response.Code":    resp.Code,
			"response.Message": resp.Message,
			"took":             took,
			"atTime":           begin,
			"error":            fmt.Sprintf("%+v", err),
		}).Info("AddRepos")
	}(time.Now())
	return s.next.AddRepos(ctx, req)
}

func (s *loggingService) GetRepos(ctx context.Context, req *ReqGetRepos) (resp *RespGetRepos, err error) {
	defer func(begin time.Time) {
		took := time.Since(begin)
		s.logger.WithFields(logrus.Fields{
			"method":             "GetRepos",
			"request.PageSize":   req.PageSize,
			"request.PageNumber": req.PageNumber,
			"response.Code":      resp.Code,
			"response.Message":   resp.Message,
			"took":               took,
			"atTime":             begin,
			"error":              fmt.Sprintf("%+v", err),
		}).Info("GetRepos")
	}(time.Now())
	return s.next.GetRepos(ctx, req)
}

func (s *loggingService) UpdateRepos(ctx context.Context, req *ReqUpdateRepos) (resp *RespUpdateRepos, err error) {
	defer func(begin time.Time) {
		took := time.Since(begin)
		s.logger.WithFields(logrus.Fields{
			"method":           "UpdateRepos",
			"request.ReposURL": req.ReposURL,
			"request.Name":     req.Name,
			"request.ID":       req.ID,
			"response.Code":    resp.Code,
			"response.Message": resp.Message,
			"took":             took,
			"atTime":           begin,
			"error":            fmt.Sprintf("%+v", err),
		}).Info("UpdateRepos")
	}(time.Now())
	return s.next.UpdateRepos(ctx, req)
}

func (s *loggingService) ArchiveRepos(ctx context.Context, req *ReqDeleteRepos) (resp *RespDeleteRepos, err error) {
	defer func(begin time.Time) {
		took := time.Since(begin)
		s.logger.WithFields(logrus.Fields{
			"method":           "DeleteRepos",
			"request.ID":       req.ID,
			"response.Code":    resp.Code,
			"response.Message": resp.Message,
			"took":             took,
			"atTime":           begin,
			"error":            fmt.Sprintf("%+v", err),
		}).Info("DeleteRepos")
	}(time.Now())
	return s.next.ArchiveRepos(ctx, req)
}

func (s *loggingService) StartScanRepos(ctx context.Context, req *ReqScan) (resp *RespScan, err error) {
	defer func(begin time.Time) {
		took := time.Since(begin)
		s.logger.WithFields(logrus.Fields{
			"method":           "StartScanRepos",
			"request.ReposId":  req.ReposId,
			"response.Code":    resp.Code,
			"response.Message": resp.Message,
			"took":             took,
			"atTime":           begin,
			"error":            fmt.Sprintf("%+v", err),
		}).Info("StartScanRepos")
	}(time.Now())
	return s.next.StartScanRepos(ctx, req)
}
func (s *loggingService) GetScanResult(ctx context.Context, req *ReqGetResult) (resp *RespGetResult, err error) {
	defer func(begin time.Time) {
		took := time.Since(begin)
		s.logger.WithFields(logrus.Fields{
			"method":           "GetScanResult",
			"request.ReposId":  req.ScanId,
			"response.Code":    resp.Code,
			"response.Message": resp.Message,
			"took":             took,
			"atTime":           begin,
			"error":            fmt.Sprintf("%+v", err),
		}).Info("GetScanResult")
	}(time.Now())
	return s.next.GetScanResult(ctx, req)
}
