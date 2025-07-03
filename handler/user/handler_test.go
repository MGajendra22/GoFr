package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MGajendra22/GoFr/model/user"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test_UserHandler : To test that interface is correctly implemented or not
func Test_UserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockUserServiceInterface(ctrl)
	h := NewUserHandler(mockSvc)

	if h == nil {
		t.Fatal("Expected non-nil handler")
	}
	if h.Service != mockSvc {
		t.Error("Expected service to be assigned correctly")
	}
}

// Test_CreateUseR : Tests user is task is created or not
func Test_CreateUseR(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Container: mockContainer,
		Request:   nil,
	}

	tests := []struct {
		name             string
		contentType      string
		input            interface{}
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{
			name:        "Success Post",
			contentType: "application/json",
			input:       user.User{Name: "John", Email: "john@gmail.com"},
			expectedResponse: gofrResponse{
				result: user.User{ID: 1, Name: "John", Email: "john@gmail.com"},
				err:    nil,
			},
			ifMock: true,
		},
		{
			name:        "Binding Failure",
			contentType: "application/json",
			input:       2,
			expectedResponse: gofrResponse{
				result: nil,
				err:    gofrHttp.ErrorInvalidParam{Params: []string{"User"}},
			},
			ifMock: false,
		},
		{
			name:        "Creation Failure",
			contentType: "application/json",
			input:       user.User{Name: "John", Email: "john@gmail.com"},
			expectedResponse: gofrResponse{
				result: user.User{},
				err:    errors.New("simulated create user error"),
			},
			ifMock: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockUserServiceInterface(ctrl)
			svc := NewUserHandler(mockService)

			var body []byte
			if str, ok := tt.input.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mockService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.result, tt.expectedResponse.err)
			}

			val, err := svc.Create(ctx)

			response := gofrResponse{val, err}

			assert.Equal(t, tt.expectedResponse.result, response.result)

			if tt.expectedResponse.err != nil {
				assert.Error(t, response.err)
				assert.Contains(t, response.err.Error(), tt.expectedResponse.err.Error())
			} else {
				assert.NoError(t, response.err)
			}
		})
	}
}

// Test_GetUseR : Tests user details are retrieved or not
func Test_GetUseR(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Container: mockContainer,
		Request:   nil,
	}

	tests := []struct {
		name             string
		id               string
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{"Success Get", "1", gofrResponse{
			result: user.User{ID: 1},
			err:    nil,
		}, true},
		{"Invalid user id", "abc", gofrResponse{
			result: nil,
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"UserID"}},
		}, false},
		{"Id not found", "99", gofrResponse{
			result: user.User{},
			err:    errors.New("simulated get user error"),
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockUserServiceInterface(ctrl)
			svc := NewUserHandler(mock)

			req := httptest.NewRequest(http.MethodGet, "/user/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.result, tt.expectedResponse.err)
			}
			val, err := svc.Get(ctx)
			response := gofrResponse{val, err}

			if tt.expectedResponse.err != nil {
				assert.Error(t, response.err)
				assert.Contains(t, response.err.Error(), tt.expectedResponse.err.Error())
			} else {
				assert.NoError(t, response.err)
			}

		})
	}
}

// Test_DeleteUseR : Tests user with user-id is deleted or not
func Test_DeleteUseR(t *testing.T) {
	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Container: mockContainer,
		Request:   nil,
	}
	tests := []struct {
		name             string
		id               string
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{"Success delete", "1", gofrResponse{
			result: fmt.Sprintf("Successfully Deleted user with id %d", 1),
			err:    nil,
		}, true},
		{"InValid user id", "abc", gofrResponse{
			result: user.User{},
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"UserID"}},
		}, false},
		{"Delete error", "99", gofrResponse{
			result: nil,
			err:    errors.New("simulated get user error"),
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mock := NewMockUserServiceInterface(ctrl)
			svc := NewUserHandler(mock)

			req := httptest.NewRequest(http.MethodDelete, "/user/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.err)
			}
			val, err := svc.Delete(ctx)
			response := gofrResponse{val, err}

			if tt.expectedResponse.err != nil {
				assert.Error(t, response.err)
				assert.Contains(t, response.err.Error(), tt.expectedResponse.err.Error())
			} else {
				assert.NoError(t, response.err)
			}

		})
	}
}

// Test_GetAllUsers : Tests all user details are retrieved or not
func Test_GetAllUsers(t *testing.T) {

	type gofrResponse struct {
		result any
		err    error
	}

	mockContainer, _ := container.NewMockContainer(t)
	ctx := &gofr.Context{
		Container: mockContainer,
		Request:   nil,
	}

	tests := []struct {
		name             string
		input            []user.User
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{"Successfully Get", []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, gofrResponse{result: []user.User{{1, "John", "John@gmail.com"}, {2, "John", "John@gmail.com"}}, err: nil}, true},
		{"Unable to fetch user data", []user.User{{1, "John", "10"}}, gofrResponse{nil, errors.New("Failed to fetch user's data")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockUserServiceInterface(ctrl)
			svc := NewUserHandler(mock)

			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().All(gomock.Any()).Return(tt.expectedResponse.result, tt.expectedResponse.err)
			}

			val, err := svc.All(ctx)
			response := gofrResponse{val, err}

			if tt.expectedResponse.err != nil {
				assert.Error(t, response.err)
				assert.Contains(t, response.err.Error(), tt.expectedResponse.err.Error())
			} else {
				assert.NoError(t, response.err)
			}

		})
	}
}
