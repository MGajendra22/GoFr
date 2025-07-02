package task

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MGajendra22/GoFr/model/task"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"strings"
	"testing"
)

type badResultForLastInsertId struct{}

func (badResultForLastInsertId) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId failed")
}

func (badResultForLastInsertId) RowsAffected() (int64, error) {
	return 1, nil
}

type badResultForRowsAffected struct{}

func (badResultForRowsAffected) LastInsertId() (int64, error) {
	return 1, nil
}
func (badResultForRowsAffected) RowsAffected() (int64, error) {
	return 0, errors.New("RowsAffected failed")
}

func Test_CreateTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewStore()

	t1 := task.Task{1, "abc", false, 2}
	t2 := task.Task{1, "abc", false, 2}
	t3 := task.Task{1, "abc", false, 2}

	mock.SQL.ExpectExec("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)").WithArgs(t2.Desc, t2.Status, t2.Userid).WillReturnError(errors.New("Insert failed"))

	_, err := str.CreateTask(ctx, t2)
	if err == nil || !strings.Contains(err.Error(), "Insert failed") {
		t.Error("expected an error, got nil")
	}

	mock.SQL.ExpectExec("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)").WithArgs(t3.Desc, t3.Status, t3.Userid).WillReturnResult(badResultForLastInsertId{})

	_, err3 := str.CreateTask(ctx, t3)
	if err3 == nil || err3.Error() != "LastInsertId failed" {
		t.Errorf("Expected LastInsertId error, got: %v", err3)
	}

	mock.SQL.ExpectExec("INSERT INTO tasks (description, status,userid) VALUES (?, ?,?)").WithArgs(t1.Desc, t1.Status, t1.Userid).WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := str.CreateTask(ctx, t1)
	if err != nil {
		t.Error("create task fail")
	}

	if res != t1 {
		t.Error("create task fail")
	}
}

func Test_GetByIDTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewStore()

	mock.SQL.ExpectQuery("SELECT * FROM tasks WHERE id = ?").WithArgs(2).WillReturnError(errors.New("Invalid Id"))

	_, err := str.GetByIDTask(ctx, 2)
	if err == nil {
		t.Error("expected an error, got nil")
	}

	rowWithScanErr := mock.SQL.NewRows([]string{"id", "desc", "status", "userid"}).AddRow("as", "abc", false, "a")

	mock.SQL.ExpectQuery("SELECT * FROM tasks WHERE id = ?").WithArgs(1).WillReturnRows(rowWithScanErr)

	_, err1 := str.GetByIDTask(ctx, 1)
	if err1 == nil || !errors.Is(err, ErrScanUser) {
		t.Error("Got scan error")
	}

	row := mock.SQL.NewRows([]string{"id", "desc", "status", "userid"}).AddRow(1, "abc", false, 1)

	mock.SQL.ExpectQuery("SELECT * FROM tasks WHERE id = ?").WithArgs(1).WillReturnRows(row)

	res, err := str.GetByIDTask(ctx, 1)
	if err != nil {
		t.Error("get task fail")
	}

	if res.Desc != "abc" || res.Status || res.Userid != 1 || res.ID != 1 {
		t.Error("get task fail")
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func Test_CompleteTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewStore()

	t1 := task.Task{1, "abc", false, 2}
	t2 := task.Task{1, "abc", false, 2}

	mock.SQL.ExpectExec("UPDATE tasks SET status = true WHERE id = ?").WithArgs(t2.ID).WillReturnError(errors.New("Not found"))

	err := str.CompleteTask(ctx, t2.ID)
	if err == nil {
		t.Error("expected an error, got nil")
	}

	mock.SQL.ExpectExec("UPDATE tasks SET status = true WHERE id = ?").WithArgs(t1.ID).WillReturnResult(badResultForRowsAffected{})

	err = str.CompleteTask(ctx, t1.ID)
	if err == nil || err.Error() != "RowsAffected failed" {
		t.Error("Rows affected fail")
	}

	mock.SQL.ExpectExec("UPDATE tasks SET status = true WHERE id = ?").WithArgs(t1.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = str.CompleteTask(ctx, t1.ID)
	if err != nil {
		t.Error("complete task fail")
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func Test_DeleteTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewStore()

	t1 := task.Task{1, "abc", false, 2}
	t2 := task.Task{1, "abc", false, 2}

	mock.SQL.ExpectExec("DELETE FROM tasks WHERE id = ?").WithArgs(t2.ID).WillReturnError(errors.New("Invalid Id"))

	err := str.DeleteTask(ctx, t2.ID)
	if err == nil {
		t.Error("expected an error, got nil")
	}

	mock.SQL.ExpectExec("DELETE FROM tasks WHERE id = ?").WithArgs(t1.ID).WillReturnResult(badResultForRowsAffected{})

	err = str.DeleteTask(ctx, t1.ID)
	if err == nil || err.Error() != "RowsAffected failed" {
		t.Error("Rows affected fail")
	}

	mock.SQL.ExpectExec("DELETE FROM tasks WHERE id = ?").WithArgs(t1.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = str.DeleteTask(ctx, t1.ID)
	if err != nil {
		t.Error("delete task fail")
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func Test_GetAllTasks(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewStore()

	mock.SQL.ExpectQuery("SELECT id, description, status , userid FROM tasks").WillReturnError(errors.New("Unable to fetch all tasks"))

	_, err := str.GetAllTask(ctx)
	if err == nil {
		t.Error("expected an error, got nil")
	}

	rowWithScanErr := mock.SQL.NewRows([]string{"id", "description", "status", "userid"}).AddRow("av", "abc", false, "as").AddRow("asd", "def", true, "as")

	mock.SQL.ExpectQuery("SELECT id, description, status , userid FROM tasks").WillReturnRows(rowWithScanErr)

	_, err = str.GetAllTask(ctx)
	if err == nil || !errors.Is(err, ErrScanUser) {
		t.Error("Got Scan error")
	}

	rows := mock.SQL.NewRows([]string{"id", "description", "status", "userid"}).AddRow(1, "abc", false, 1).AddRow(2, "def", true, 2)

	mock.SQL.ExpectQuery("SELECT id, description, status , userid FROM tasks").WillReturnRows(rows)

	tasks, err := str.GetAllTask(ctx)
	if err != nil {
		t.Error("get all tasks fail")
	}

	if len(tasks) != 2 {
		t.Error("get all tasks fail")
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}

}

func Test_GetTasksByUserIDTask(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewStore()

	t1 := task.Task{1, "abc", false, 2}
	t2 := task.Task{1, "abc", false, 2}

	mock.SQL.ExpectQuery("SELECT id, description, status , userid FROM tasks where userid =?").WithArgs(t1.Userid).WillReturnError(errors.New("Not found"))

	_, err := str.GetTasksByUserIDTask(ctx, t2.Userid)
	if err == nil {
		t.Error("expected an error, got nil")
	}

	rowWithScanErr := mock.SQL.NewRows([]string{"id", "description", "status", "userid"}).AddRow(1, "abc", false, 1).AddRow("dwa", "def", true, "dad")

	mock.SQL.ExpectQuery("SELECT id, description, status , userid FROM tasks where userid =?").
		WithArgs(t2.Userid).WillReturnRows(rowWithScanErr)

	_, err = str.GetTasksByUserIDTask(ctx, t2.Userid)
	if err == nil || !errors.Is(err, ErrScanUser) {
		t.Error("Got Scan error")
	}

	rows := mock.SQL.NewRows([]string{"id", "description", "status", "userid"}).AddRow(1, "abc", false, 1).AddRow(2, "def", true, 1)

	mock.SQL.ExpectQuery("SELECT id, description, status , userid FROM tasks where userid =?").WithArgs(t2.Userid).WillReturnRows(rows)

	tasks, err := str.GetTasksByUserIDTask(ctx, t2.Userid)
	if err != nil {
		t.Error("get tasks by user id fail")
	}

	if len(tasks) != 2 {
		t.Error("get tasks by user id fail")
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}
