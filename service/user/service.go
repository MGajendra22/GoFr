package user

import (
	"github.com/MGajendra22/GoFr/model/user"
	"gofr.dev/pkg/gofr"
)

type UserService struct {
	store UserStoreInterface
}

func NewUserService(store UserStoreInterface) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Create(c *gofr.Context, u user.User) (user.User, error) {
	if err := u.Validate(); err != nil {
		return u, err
	}

	return s.store.CreateUser(c, u)
}

func (s *UserService) Get(c *gofr.Context, id int) (user.User, error) {
	return s.store.GetByIDUser(c, id)
}

func (s *UserService) Delete(c *gofr.Context, id int) error {
	return s.store.DeleteUser(c, id)
}

func (s *UserService) All(c *gofr.Context) ([]user.User, error) {
	return s.store.GetAllUser(c)

}
