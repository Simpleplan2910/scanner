package repos

import (
	"context"
	"fmt"
	apierror "scanner/pkg/apiError"
	dbMocks "scanner/pkg/db/mocks"
	gitMocks "scanner/pkg/services/git/mocks"
	queueMocks "scanner/pkg/services/queueJob/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_AddRepos(t *testing.T) {
	assert := assert.New(t)
	gitMock := new(gitMocks.Service)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	queueMock := new(queueMocks.QueueJob)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, queueMock)
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

func Test_AddReposDBError(t *testing.T) {
	assert := assert.New(t)
	gitMock := new(gitMocks.Service)
	reposStoreMock := new(dbMocks.ReposStore)
	resultStoreMock := new(dbMocks.ResultStore)
	queueMock := new(queueMocks.QueueJob)

	service := NewService(gitMock, reposStoreMock, resultStoreMock, queueMock)
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
