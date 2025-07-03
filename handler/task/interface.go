package task

import (
	"github.com/MGajendra22/GoFr/model/task"
	"gofr.dev/pkg/gofr"
)

type TaskServiceInterface interface {
	Create(c *gofr.Context, t task.Task) (task.Task, error)
	GetTask(c *gofr.Context, id int) (task.Task, error)
	Complete(c *gofr.Context, id int) error
	Delete(c *gofr.Context, id int) error
	All(c *gofr.Context) ([]task.Task, error)
	GetTasksByUserID(c *gofr.Context, userId int) ([]task.Task, error)
}
