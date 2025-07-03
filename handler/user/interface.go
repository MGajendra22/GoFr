package user

import (
	"github.com/MGajendra22/GoFr/model/user"
	"gofr.dev/pkg/gofr"
)

type UserServiceInterface interface {
	Create(c *gofr.Context, u user.User) (user.User, error)
	Get(c *gofr.Context, id int) (user.User, error)
	Delete(c *gofr.Context, id int) error
	All(c *gofr.Context) ([]user.User, error)
}
