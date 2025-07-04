package user

import (
	"errors"
	_ "github.com/MGajendra22/GoFr/model/task"
	"github.com/MGajendra22/GoFr/model/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"testing"
)

func Test_CreateUser(t *testing.T) {
	tests := []struct {
		name       string
		input      user.User
		mockOutput user.User
		mockError  error
		expErr     bool
	}{
		{
			name:       "Valid user",
			input:      user.User{ID: 1, Name: "Alice", Email: "alice@example.com"},
			mockOutput: user.User{ID: 1, Name: "Alice", Email: "alice@example.com"},
			mockError:  nil,
			expErr:     false,
		},
		{
			name:   "Validation error",
			input:  user.User{ID: 2, Name: "", Email: ""},
			expErr: true,
		},

		{
			name:      "Store error",
			input:     user.User{ID: 3, Name: "Bob", Email: "bob@example.com"},
			mockError: errors.New("db error"),
			expErr:    true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			mockContainer, _ := container.NewMockContainer(t)

			ctx := &gofr.Context{
				Container: mockContainer,
			}

			mockstore := NewMockUserStoreInterface(ctrl)

			service := NewUserService(mockstore)

			mockstore.EXPECT().CreateUser(ctx, gomock.Any()).Return(tt.mockOutput, tt.mockError).AnyTimes()

			result, err := service.Create(ctx, tt.input)

			if tt.expErr {
				assert.Error(t, err, tt.name)
			} else {
				assert.NoError(t, err, tt.name)
				assert.Equal(t, tt.mockOutput, result)
			}
		})
	}
}

func Test_GetUser(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		mockOutput user.User
		mockErr    error
		expErr     bool
	}{
		{"Valid Id", 1, user.User{1, "John", "mail"}, nil, false},
		{"User not found", 2, user.User{}, errors.New("task not found"), true},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		mockstore := NewMockUserStoreInterface(ctrl)

		service := NewUserService(mockstore)

		mockContainer, _ := container.NewMockContainer(t)

		ctx := &gofr.Context{
			Container: mockContainer,
		}

		mockstore.EXPECT().GetByIDUser(ctx, tt.id).Return(tt.mockOutput, tt.mockErr).AnyTimes()

		result, err := service.Get(ctx, tt.id)

		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.Equal(t, tt.mockOutput, result)
		}
	}
}

func Test_DeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		taskErr error
		expErr  bool
	}{
		{
			name:   "Valid user Deletion",
			input:  1,
			expErr: false,
		},
		{
			name:    "User Not Found",
			input:   1,
			taskErr: errors.New("user not found"),
			expErr:  true,
		},
	}
	for _, tt := range tests {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		mockstore := NewMockUserStoreInterface(ctrl)

		service := NewUserService(mockstore)

		mockContainer, _ := container.NewMockContainer(t)

		ctx := &gofr.Context{
			Container: mockContainer,
		}

		mockstore.EXPECT().DeleteUser(ctx, tt.input).Return(tt.taskErr).AnyTimes()

		err := service.Delete(ctx, tt.input)

		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
		}
	}
}

func Test_GetAllUsers(t *testing.T) {
	tests := []struct {
		name       string
		mockOutput []user.User
		mockErr    error
		expErr     bool
	}{
		{"Data fetched", []user.User{{1, "John", "mail"}}, nil, false},
		{"Unable to fetch", []user.User{}, errors.New("task not found"), true},
	}

	for _, tt := range tests {
		ctrl := gomock.NewController(t)

		defer ctrl.Finish()

		mockstore := NewMockUserStoreInterface(ctrl)

		service := NewUserService(mockstore)

		mockContainer, _ := container.NewMockContainer(t)

		ctx := &gofr.Context{
			Container: mockContainer,
		}

		mockstore.EXPECT().GetAllUser(ctx).Return(tt.mockOutput, tt.mockErr).AnyTimes()

		res, err := service.All(ctx)

		if tt.expErr {
			assert.Error(t, err, tt.name)
		} else {
			assert.NoError(t, err, tt.name)
			assert.Equal(t, tt.mockOutput, res)
		}
	}
}
