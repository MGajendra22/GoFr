package user

import (
	"github.com/MGajendra22/GoFr/model/user"
	"gofr.dev/pkg/gofr"
)

type UserStoreInterface interface {
	CreateUser(c *gofr.Context, u user.User) (user.User, error)
	GetByIDUser(c *gofr.Context, id int) (user.User, error)
	DeleteUser(c *gofr.Context, id int) error
	GetAllUser(c *gofr.Context) ([]user.User, error)
}
