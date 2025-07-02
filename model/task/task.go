package task

import (
	gofrHttp "gofr.dev/pkg/gofr/http"
)

type Task struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status bool   `json:"status"`
	Userid int    `json:"userid"`
}

func (t *Task) Validate() error {
	if t.Desc == "" {
		return gofrHttp.ErrorInvalidParam{Params: []string{"task.desc"}}
	}

	return nil
}
