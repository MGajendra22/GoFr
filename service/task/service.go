package task

import (
	"fmt"
	"github.com/MGajendra22/GoFr/model/task"
	"gofr.dev/pkg/gofr"
)

type TaskService struct {
	str            TaskStoreInterface
	userServiceref UserServiceInterface
}

func NewService(s TaskStoreInterface, us UserServiceInterface) *TaskService {
	return &TaskService{
		str:            s,
		userServiceref: us,
	}
}

func (s *TaskService) Create(c *gofr.Context, t task.Task) (task.Task, error) {
	if err := t.Validate(); err != nil {
		return t, err
	}

	_, err := s.userServiceref.Get(c, t.Userid)
	if err != nil {
		return t, fmt.Errorf("user with ID %d does not exist: %v", t.Userid, err)
	}

	return s.str.CreateTask(c, t)
}

func (s *TaskService) GetTask(c *gofr.Context, id int) (task.Task, error) {
	return s.str.GetByIDTask(c, id)
}

func (s *TaskService) Complete(c *gofr.Context, id int) error {
	return s.str.CompleteTask(c, id)
}

func (s *TaskService) Delete(c *gofr.Context, id int) error {
	return s.str.DeleteTask(c, id)
}

func (s *TaskService) All(c *gofr.Context) ([]task.Task, error) {
	return s.str.GetAllTask(c)
}

func (s *TaskService) GetTasksByUserID(c *gofr.Context, userid int) ([]task.Task, error) {
	_, err := s.userServiceref.Get(c, userid)

	if err != nil {
		return nil, err
	}

	return s.str.GetTasksByUserIDTask(c, userid)
}
