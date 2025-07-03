package task

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MGajendra22/GoFr/model/task"
	"gofr.dev/pkg/gofr"
)

type Store struct {
}

func NewStore() *Store {
	return &Store{}
}

var ErrScanUser = errors.New("scan user failed")

// CreateTask inserts a new task into the database
func (*Store) CreateTask(c *gofr.Context, t task.Task) (task.Task, error) {
	DB := c.SQL

	res, err := DB.Exec("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)", t.Desc, t.Status, t.Userid)
	if err != nil {
		return t, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return t, err
	}

	t.ID = int(id)

	return t, nil
}

// GetByIDTask fetches a task by its ID
func (*Store) GetByIDTask(c *gofr.Context, id int) (task.Task, error) {
	DB := c.SQL

	var t task.Task

	err := DB.QueryRow("SELECT * FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Desc, &t.Status, &t.Userid)
	if err != nil {
		return t, fmt.Errorf("%w: %v", ErrScanUser, err)
	}

	return t, err
}

// CompleteTask marks a task as completed
func (*Store) CompleteTask(c *gofr.Context, id int) error {
	DB := c.SQL

	res, err := DB.Exec("UPDATE tasks SET status = true WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteTask removes a task by ID
func (*Store) DeleteTask(c *gofr.Context, id int) error {
	DB := c.SQL

	res, err := DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetAllTask returns all tasks from the database
func (*Store) GetAllTask(c *gofr.Context) ([]task.Task, error) {
	DB := c.SQL

	rows, err := DB.Query("SELECT id, description, status , userid FROM tasks")
	if err != nil {
		return nil, err
	}

	var tasks []task.Task

	for rows.Next() {
		var t task.Task

		if err := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.Userid); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanUser, err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTasksByUserID it will send the tasks , which are assigned to user
func (*Store) GetTasksByUserIDTask(c *gofr.Context, userid int) ([]task.Task, error) {
	DB := c.SQL

	rows, err := DB.Query("SELECT id, description, status , userid FROM tasks where userid =?", userid)
	if err != nil {
		return nil, err
	}

	//defer func(rows *sql.Rows) {
	//	err := rows.Close()
	//	if err != nil {
	//		return
	//	}
	//}(rows)

	var tasks []task.Task

	for rows.Next() {
		var t task.Task

		if err := rows.Scan(&t.ID, &t.Desc, &t.Status, &t.Userid); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanUser, err)
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
