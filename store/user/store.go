package user

import (
	"errors"
	"fmt"
	"github.com/MGajendra22/GoFr/model/user"
	"gofr.dev/pkg/gofr"
)

type UserStore struct {
}

func NewUserStore() *UserStore {
	return &UserStore{}
}

var ErrScanUser = errors.New("scan user failed")

func (*UserStore) CreateUser(c *gofr.Context, user user.User) (user.User, error) {
	DB := c.SQL

	query := "INSERT INTO users (name, email) VALUES (?, ?)"

	result, err := DB.Exec(query, user.Name, user.Email)
	if err != nil {
		return user, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = int(id)

	return user, nil
}

func (*UserStore) GetByIDUser(c *gofr.Context, id int) (user.User, error) {
	DB := c.SQL

	var user user.User

	query := "SELECT id, name, email FROM users WHERE id = ?"

	err := DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return user, err
	}

	return user, err
}

func (*UserStore) DeleteUser(c *gofr.Context, id int) error {
	DB := c.SQL

	_, err := DB.Exec("DELETE FROM users WHERE id = ?", id)

	return err
}

func (*UserStore) GetAllUser(c *gofr.Context) ([]user.User, error) {
	DB := c.SQL

	query := "SELECT id, name, email FROM users"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	var users []user.User

	for rows.Next() {
		var u user.User

		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanUser, err)
		}

		users = append(users, u)
	}

	return users, nil
}
