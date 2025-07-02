package task

import (
	"github.com/MGajendra22/GoFr/model/task"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"strconv"
)

type handler struct {
	svc TaskServiceInterface
}

// NewHandler : Factory function to implement and return behaviour
func NewHandler(s TaskServiceInterface) *handler {
	return &handler{svc: s}
}

func (h *handler) Create(c *gofr.Context) (any, error) {

	var t task.Task

	if err := c.Bind(&t); err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"body"}}
	}

	if err := t.Validate(); err != nil {
		return nil, err
	}

	task1, err := h.svc.Create(c, t)
	if err != nil {
		return task1, err
	}

	return task1, nil
}

func (h *handler) GetTask(c *gofr.Context) (any, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"body"}}
	}

	task1, err := h.svc.GetTask(c, id)
	if err != nil {
		return task1, err
	}

	return task1, nil

}

func (h *handler) Complete(c *gofr.Context) (any, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	if err := h.svc.Complete(c, id); err != nil {
		return nil, err
	}

	return task.Task{}, nil

}

func (h *handler) GetTasksByUserID(c *gofr.Context) (any, error) {
	userid, err := strconv.Atoi(c.PathParam("userid"))
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	tasks, err := h.svc.GetTasksByUserID(c, userid)
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (h *handler) Delete(c *gofr.Context) (any, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"id"}}
	}

	if err := h.svc.Delete(c, id); err != nil {
		return nil, err
	}

	return task.Task{}, nil
}

func (h *handler) All(c *gofr.Context) (any, error) {

	tasks, err := h.svc.All(c)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
