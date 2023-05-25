package repos

import (
	"context"
	"fmt"
	apierror "scanner/pkg/apiError"
	"scanner/pkg/db"
	dbMocks "scanner/pkg/db/mocks"
	gitMocks "scanner/pkg/services/git/mocks"
	smocks "scanner/pkg/services/scanner/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddRepos(t *testing.T) {

	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)
	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqAddRepos
		response         *RespAddRepos
	}
	tt := testTable{
		{
			name: "success",
			request: &ReqAddRepos{
				Name:     "test",
				ReposURL: "https://github.com/vektra/mockery",
			},
			isExpectNilError: true,
			response: &RespAddRepos{
				Code: 0,
			},
		},
		{
			name: "empty name",
			request: &ReqAddRepos{
				Name:     "",
				ReposURL: "testURL",
			},
			isExpectNilError: false,
			response: &RespAddRepos{
				Code: apierror.InvalidRequest,
			},
		},
		{
			name: "empty url",
			request: &ReqAddRepos{
				Name:     "test",
				ReposURL: "",
			},
			isExpectNilError: false,
			response: &RespAddRepos{
				Code: apierror.InvalidRequest,
			},
		},
		{
			name: "wrong github url",
			request: &ReqAddRepos{
				Name:     "test",
				ReposURL: "https://google.com/vektra/mockery/",
			},
			isExpectNilError: false,
			response: &RespAddRepos{
				Code: apierror.InvalidRequest,
			},
		},
	}
	ctx := context.Background()
	reposStoreMock.On("Add", mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.AddRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestAddReposDBError(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)
	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqAddRepos
		response         *RespAddRepos
	}
	tt := testTable{
		{
			name: "success",
			request: &ReqAddRepos{
				Name:     "test",
				ReposURL: "testURL",
			},
			isExpectNilError: false,
			response: &RespAddRepos{
				Code: apierror.InternalServerError,
			},
		},
	}
	ctx := context.Background()
	reposStoreMock.On("Add", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test error"))
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.AddRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestStartScanRepos(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqScan
		response         *RespScan
	}

	tt := testTable{
		{
			name: "empty reposID",
			request: &ReqScan{
				Substr: "dd",
			},
			response: &RespScan{
				Code: apierror.InvalidRequest,
			},
			isExpectNilError: false,
		},
		{
			name: "empty sub string",
			request: &ReqScan{
				ReposId: primitive.NewObjectID(),
				Substr:  "",
			},
			response: &RespScan{
				Code: apierror.InvalidRequest,
			},
			isExpectNilError: false,
		},
		{
			name: "success scan",
			request: &ReqScan{
				ReposId: primitive.NewObjectID(),
				Substr:  "fff",
			},
			response: &RespScan{
				Code: 0,
			},
			isExpectNilError: true,
		},
	}

	ctx := context.Background()
	reposMock := &db.Repos{
		ID:        primitive.NewObjectID(),
		Name:      "test name",
		ReposURL:  "test url",
		IsArchive: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reposStoreMock.On("Get", mock.Anything, mock.Anything).Return(reposMock, nil)
	scanStoreM.On("Add", mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil)
	scannerM.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.StartScanRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestStartScanReposDBFailed(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqScan
		response         *RespScan
	}

	tt := testTable{
		{
			name: "get repos failed",
			request: &ReqScan{
				ReposId: primitive.NewObjectID(),
				Substr:  "dd",
			},
			response: &RespScan{
				Code: apierror.InternalServerError,
			},
			isExpectNilError: false,
		},
	}

	ctx := context.Background()
	reposStoreMock.On("Get", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test error"))
	scanStoreM.On("Add", mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil)
	scannerM.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.StartScanRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestStartScanReposAddDBFailed(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqScan
		response         *RespScan
	}

	tt := testTable{
		{
			name: "get repos failed",
			request: &ReqScan{
				ReposId: primitive.NewObjectID(),
				Substr:  "dd",
			},
			response: &RespScan{
				Code: apierror.InternalServerError,
			},
			isExpectNilError: false,
		},
	}

	ctx := context.Background()
	reposMock := &db.Repos{
		ID:        primitive.NewObjectID(),
		Name:      "test name",
		ReposURL:  "test url",
		IsArchive: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reposStoreMock.On("Get", mock.Anything, mock.Anything).Return(reposMock, nil)
	scanStoreM.On("Add", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("test error"))
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.StartScanRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestStartScanReposScanFailed(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqScan
		response         *RespScan
	}

	tt := testTable{
		{
			name: "get repos failed",
			request: &ReqScan{
				ReposId: primitive.NewObjectID(),
				Substr:  "dd",
			},
			response: &RespScan{
				Code: apierror.InternalServerError,
			},
			isExpectNilError: false,
		},
	}

	ctx := context.Background()
	reposMock := &db.Repos{
		ID:        primitive.NewObjectID(),
		Name:      "test name",
		ReposURL:  "test url",
		IsArchive: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	reposStoreMock.On("Get", mock.Anything, mock.Anything).Return(reposMock, nil)
	scanStoreM.On("Add", mock.Anything, mock.Anything).Return(primitive.NewObjectID(), nil)
	scannerM.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("test error"))
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.StartScanRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestGetRepos(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqGetRepos
		response         *RespGetRepos
	}

	tt := testTable{
		{
			name: "page size is less than 1",
			request: &ReqGetRepos{
				PageSize:   0,
				PageNumber: 1,
			},
			response: &RespGetRepos{
				Code: apierror.InvalidRequest,
			},
			isExpectNilError: false,
		},
		{
			name: "page number is less than 1",
			request: &ReqGetRepos{
				PageSize:   1,
				PageNumber: 0,
			},
			response: &RespGetRepos{
				Code: apierror.InvalidRequest,
			},
			isExpectNilError: false,
		},
		{
			name: "success get result",
			request: &ReqGetRepos{
				PageSize:   1,
				PageNumber: 1,
			},
			response: &RespGetRepos{
				Code: 0,
			},
			isExpectNilError: true,
		},
	}

	ctx := context.Background()
	reposMock := []db.Repos{
		{
			ID:        primitive.NewObjectID(),
			Name:      "test name",
			ReposURL:  "test url",
			IsArchive: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	reposStoreMock.On("Filter", mock.Anything, mock.Anything).Return(reposMock, int64(1), nil)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.GetRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestUpdateRepos(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqUpdateRepos
		response         *RespUpdateRepos
	}
	tt := testTable{
		{
			name:             "empty id",
			isExpectNilError: false,
			request: &ReqUpdateRepos{
				Name:     "dd",
				ReposURL: "https://github.com/vektra/mockery",
			},
			response: &RespUpdateRepos{
				Code: apierror.InvalidRequest,
			},
		},
		{
			name:             "empty name and url",
			isExpectNilError: false,
			request: &ReqUpdateRepos{
				ID:       primitive.NewObjectID(),
				Name:     "",
				ReposURL: "",
			},
			response: &RespUpdateRepos{
				Code: apierror.InvalidRequest,
			},
		},
		{
			name:             "invalid url",
			isExpectNilError: false,
			request: &ReqUpdateRepos{
				ID:       primitive.NewObjectID(),
				Name:     "dd",
				ReposURL: "ddasfasf",
			},
			response: &RespUpdateRepos{
				Code: apierror.InvalidRequest,
			},
		},
		{
			name:             "success",
			isExpectNilError: true,
			request: &ReqUpdateRepos{
				ID:       primitive.NewObjectID(),
				Name:     "dd",
				ReposURL: "https://github.com/vektra/mockery",
			},
			response: &RespUpdateRepos{
				Code: 0,
			},
		},
	}
	ctx := context.Background()
	reposStoreMock.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.UpdateRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, fmt.Sprintf("%s: error should be nil", tc.name))
			} else {
				assert.NotNil(err, fmt.Sprintf("%s: error shouldn't be nil", tc.name))
			}

			assert.Equal(tc.response.Code, resp.Code, fmt.Sprintf("%s: return the wrong code", tc.name))
		})
	}
}

func TestArchiveRepos(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqDeleteRepos
		response         *RespDeleteRepos
	}

	tt := testTable{
		{
			name:             "empty id",
			isExpectNilError: false,
			request: &ReqDeleteRepos{
				ID: primitive.NilObjectID,
			},
			response: &RespDeleteRepos{
				Code: apierror.InvalidRequest,
			},
		},
		{
			name:             "success",
			isExpectNilError: true,
			request: &ReqDeleteRepos{
				ID: primitive.NewObjectID(),
			},
			response: &RespDeleteRepos{
				Code: 0,
			},
		},
	}

	ctx := context.Background()
	reposStoreMock.On("Archive", mock.Anything, mock.Anything).Return(nil)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.ArchiveRepos(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}

func TestGetScanResult(t *testing.T) {
	assert := assert.New(t)
	gitMock := gitMocks.NewService(t)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	scanStoreM := new(dbMocks.ScanStore)
	scannerM := smocks.NewService(t)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, scannerM, scanStoreM)

	type testTable []struct {
		name             string
		isExpectNilError bool
		request          *ReqGetResult
		response         *RespGetResult
	}

	tt := testTable{
		{
			name: "empty id",
			request: &ReqGetResult{
				ScanId:     primitive.NilObjectID,
				PageSize:   1,
				PageNumber: 1,
			},
			response: &RespGetResult{
				Code: apierror.InvalidRequest,
			},

			isExpectNilError: false,
		},
		{
			name: "invalid page size",
			request: &ReqGetResult{
				ScanId:     primitive.NewObjectID(),
				PageSize:   0,
				PageNumber: 1,
			},
			response: &RespGetResult{
				Code: apierror.InvalidRequest,
			},
			isExpectNilError: false,
		},
		{
			name: "invalid page number",
			request: &ReqGetResult{
				ScanId:     primitive.NewObjectID(),
				PageSize:   1,
				PageNumber: 0,
			},
			response: &RespGetResult{
				Code: apierror.InvalidRequest,
			},
			isExpectNilError: false,
		},
		{
			name: "success",
			request: &ReqGetResult{
				ScanId:     primitive.NewObjectID(),
				PageSize:   1,
				PageNumber: 1,
			},
			response: &RespGetResult{
				Code: 0,
			},
			isExpectNilError: true,
		},
	}

	ctx := context.Background()
	resultStoreMock.On("Filter", mock.Anything, mock.Anything).Return(nil, int64(1), nil)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := service.GetScanResult(ctx, tc.request)
			if tc.isExpectNilError {
				assert.Nil(err, "error should be nil")
			} else {
				assert.NotNil(err, "error shouldn't be nil")
			}

			assert.Equal(tc.response.Code, resp.Code, "return the wrong code")
		})
	}
}
