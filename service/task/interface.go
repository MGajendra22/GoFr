package task

import (
	"github.com/MGajendra22/GoFr/model/task"
	userModel "github.com/MGajendra22/GoFr/model/user"
	"gofr.dev/pkg/gofr"
)

type TaskStoreInterface interface {
	CreateTask(c *gofr.Context, task task.Task) (task.Task, error)
	GetByIDTask(c *gofr.Context, id int) (task.Task, error)
	GetAllTask(c *gofr.Context) ([]task.Task, error)
	CompleteTask(c *gofr.Context, id int) error
	DeleteTask(c *gofr.Context, id int) error
	GetTasksByUserIDTask(c *gofr.Context, userId int) ([]task.Task, error)
}

type UserServiceInterface interface {
	Get(c *gofr.Context, id int) (userModel.User, error)
}
