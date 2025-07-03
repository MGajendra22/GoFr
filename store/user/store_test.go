package user

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	user "github.com/MGajendra22/GoFr/model/user"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"testing"
)

type badResult struct{}

func (badResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId failed")
}

func (badResult) RowsAffected() (int64, error) {
	return 1, nil
}

func Test_CreateUser(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewUserStore()

	u1 := user.User{Name: "John", Email: "john@nidevrrtvwn.com"}
	u2 := user.User{Name: "Johrvrn", Email: "john@nvrrvn.com"}
	u3 := user.User{Name: "Jvrohn", Email: "john@rvrvrvnidebiwn.com"}

	mock.SQL.ExpectExec("INSERT INTO users (name, email) VALUES (?, ?)").WithArgs(u2.Name, u2.Email).WillReturnError(errors.New("error"))

	_, err1 := str.CreateUser(ctx, u2)
	if err1 == nil {
		t.Error("expected an error, got nil")
	}

	mock.SQL.ExpectExec("INSERT INTO users (name, email) VALUES (?, ?)").WithArgs(u3.Name, u3.Email).WillReturnResult(badResult{})

	_, err3 := str.CreateUser(ctx, u3)
	if err3 == nil || err3.Error() != "LastInsertId failed" {
		t.Errorf("Expected LastInsertId error, got: %v", err3)
	}

	mock.SQL.ExpectExec("INSERT INTO users (name, email) VALUES (?, ?)").WithArgs(u1.Name, u1.Email).WillReturnResult(sqlmock.NewResult(1, 1))

	getUser, err := str.CreateUser(ctx, u1)
	if err != nil {
		t.Error(err)
	}

	if getUser.ID != 1 {
		t.Error("Expected 1, got ", getUser.ID)
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet SQL expectations: %v", err)
	}
}

func Test_GetByIDUser(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}
	str := NewUserStore()

	rows := mock.SQL.NewRows([]string{"id", "name", "email"}).AddRow(1, "John Doe", "john@example.com")
	mock.SQL.ExpectQuery("SELECT id, name, email FROM users WHERE id = ?").WithArgs(2).WillReturnError(errors.New("Id not found"))

	mock.SQL.ExpectQuery("SELECT id, name, email FROM users WHERE id = ?").WithArgs(1).WillReturnRows(rows)

	_, err1 := str.GetByIDUser(ctx, 2)
	if err1 == nil {
		t.Error("expected error, got nil")
	}

	u, err := str.GetByIDUser(ctx, 1)
	if err != nil {
		t.Error(err)
	}

	if u.ID != 1 {
		t.Error("Expected 1, got ", u.ID)
	}
}

func Test_DeleteUser(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewUserStore()

	u1 := user.User{ID: 1}
	u2 := user.User{ID: 1}
	mock.SQL.ExpectExec("DELETE FROM users WHERE id = ?").WithArgs(u2.ID).WillReturnError(errors.New("User with id not found"))

	mock.SQL.ExpectExec("DELETE FROM users WHERE id = ?").WithArgs(u1.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := str.DeleteUser(ctx, u2.ID)
	if err == nil {
		t.Error("expected error, got nil")
	}

	err = str.DeleteUser(ctx, u1.ID)
	if err != nil {
		t.Error(err)
	}

	if err := mock.SQL.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %s", err)
	}
}

func Test_GetAllUsers(t *testing.T) {
	mockContainer, mock := container.NewMockContainer(t)

	ctx := &gofr.Context{
		Container: mockContainer,
	}

	str := NewUserStore()

	rows := mock.SQL.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "John Doe", "john@example.com").
		AddRow(2, "John Doe", "john@example.com")
	mock.SQL.ExpectQuery("SELECT id, name, email FROM users").WillReturnError(errors.New("Unable to fetch all users"))

	_, err := str.GetAllUser(ctx)
	if err == nil {
		t.Error("expected error, got nil")
	}

	rowsWithScanErr := mock.SQL.NewRows([]string{"id", "name", "email"}).AddRow("invalid-id", "Jane", "jane@example.com")
	mock.SQL.ExpectQuery("SELECT id, name, email FROM users").WillReturnRows(rowsWithScanErr)

	_, err = str.GetAllUser(ctx)
	if err == nil || !errors.Is(err, ErrScanUser) {
		t.Error("expected error, got nil")

	}

	mock.SQL.ExpectQuery("SELECT id, name, email FROM users").WillReturnRows(rows)

	users, err := str.GetAllUser(ctx)
	if err != nil {
		t.Error(err)
	}

	if len(users) != 2 {
		t.Error("Expected 2 users, got ", len(users))
	}
}
