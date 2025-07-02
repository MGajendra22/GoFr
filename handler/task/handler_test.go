package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MGajendra22/GoFr/model/task"
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

// Test_NewHandler : To test that interface is correctly implemented or not
func Test_NewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockTaskServiceInterface(ctrl)

	h := NewHandler(mockSvc)

	if h == nil {
		t.Fatal("Expected non-nil handler")
	}

	if h.svc != mockSvc {
		t.Error("Expected service to be assigned correctly")
	}
}

// Test_CreateTasK : Tests task is created or not
func Test_CreateTasK(t *testing.T) {
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
		{"Success Create", "application/json", task.Task{1, "Working", false, 1}, gofrResponse{result: task.Task{1, "Working", false, 1}, err: nil}, true},
		{"Binding Error", "application/json", 10, gofrResponse{
			result: nil,
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"body"}},
		}, false},
		{"User id not found", "application/json", task.Task{1, "", false, 100}, gofrResponse{
			result: nil,
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"task.desc"}},
		}, false},
		{name: "Creation Failure",
			contentType: "application/json",
			input:       task.Task{1, "Working", false, 1},
			expectedResponse: gofrResponse{
				result: task.Task{},
				err:    errors.New("simulated create user error"),
			},
			ifMock: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockTaskServiceInterface(ctrl)
			svc := NewHandler(mock)

			var body []byte
			if str, ok := tt.input.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(body))
			req.Header.Set("Content-Type", tt.contentType)
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.result, tt.expectedResponse.err)
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

// Test_GetTasK : Tests task is retrieved or not
func Test_GetTasK(t *testing.T) {
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
			result: task.Task{1, "Working", false, 1},
			err:    nil,
		}, true},
		{"Invalid user id", "abc", gofrResponse{
			result: nil,
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"body"}},
		}, false},
		{"Id not found", "99", gofrResponse{
			result: task.Task{},
			err:    errors.New("simulated get user error"),
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockTaskServiceInterface(ctrl)
			svc := NewHandler(mock)

			req := httptest.NewRequest(http.MethodGet, "/task/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().GetTask(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.result, tt.expectedResponse.err)
			}

			val, err := svc.GetTask(ctx)
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

// Test_GetTasksByUserID : Tests task is retrieved or not of User
func Test_GetTasksByUserID(t *testing.T) {
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
		userid           string
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{"Success Get", "1", gofrResponse{
			result: []task.Task{{1, "Working", false, 1}},
			err:    nil,
		}, true},
		{"Invalid user id", "abc", gofrResponse{
			result: nil,
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		}, false},
		{"Id not found", "99", gofrResponse{
			result: []task.Task{},
			err:    errors.New("simulated get user error"),
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockTaskServiceInterface(ctrl)
			svc := NewHandler(mock)

			req := httptest.NewRequest(http.MethodGet, "/task/user/"+tt.userid, nil)
			req = mux.SetURLVars(req, map[string]string{"userid": tt.userid})
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().GetTasksByUserID(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.result, tt.expectedResponse.err)
			}

			val, err := svc.GetTasksByUserID(ctx)
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

// Test_CompleteTask : Tests assigned id's task is completed or not
func Test_CompleteTask(t *testing.T) {
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
		{"Success Update", "1", gofrResponse{
			result: task.Task{Status: true},
			err:    nil,
		}, true},
		{"Invalid user id", "abc", gofrResponse{
			result: nil,
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		}, false},
		{"Id not found", "99", gofrResponse{
			result: nil,
			err:    errors.New("simulated get user error"),
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockTaskServiceInterface(ctrl)
			svc := NewHandler(mock)

			req := httptest.NewRequest(http.MethodPut, "/task/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().Complete(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.err)
			}

			_, err := svc.Complete(ctx)
			response := gofrResponse{nil, err}

			if tt.expectedResponse.err != nil {
				assert.Error(t, response.err)
				assert.Contains(t, response.err.Error(), tt.expectedResponse.err.Error())
			} else {
				assert.NoError(t, response.err)
			}
		})

	}
}

// Test_AllTasks : Tests all tasks are retrieved or not

func Test_AllTasks(t *testing.T) {
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
		input            []task.Task
		expectedResponse gofrResponse
		ifMock           bool
	}{
		{"Successfully Get", []task.Task{{1, "Working", false, 1}}, gofrResponse{result: []task.Task{{1, "Working", false, 1}}, err: nil}, true},
		{"Unable to fetch user data", []task.Task{{1, "Working", false, 1}}, gofrResponse{nil, errors.New("Failed to fetch user's data")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockTaskServiceInterface(ctrl)
			svc := NewHandler(mock)

			req := httptest.NewRequest(http.MethodGet, "/task/", nil)
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

// Test_DeleteTask : Tests Task with task-id is deleted or not
func Test_DeleteTask(t *testing.T) {
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
		{"Success Delete", "1", gofrResponse{
			result: task.Task{},
			err:    nil,
		}, true},
		{"Invalid user id", "abc", gofrResponse{
			result: nil,
			err:    gofrHttp.ErrorInvalidParam{Params: []string{"id"}},
		}, false},
		{"Id not found", "99", gofrResponse{
			result: nil,
			err:    errors.New("simulated get user error"),
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock := NewMockTaskServiceInterface(ctrl)
			svc := NewHandler(mock)

			req := httptest.NewRequest(http.MethodDelete, "/task/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			request := gofrHttp.NewRequest(req)
			ctx.Request = request

			if tt.ifMock {
				mock.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.expectedResponse.err)
			}

			_, err := svc.Delete(ctx)
			response := gofrResponse{nil, err}

			if tt.expectedResponse.err != nil {
				assert.Error(t, response.err)
				assert.Contains(t, response.err.Error(), tt.expectedResponse.err.Error())
			} else {
				assert.NoError(t, response.err)
			}
		})

	}
}
