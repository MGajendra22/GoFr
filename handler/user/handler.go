package user

import (
	"fmt"
	"github.com/MGajendra22/GoFr/model/user"
	"gofr.dev/pkg/gofr"
	gofrHttp "gofr.dev/pkg/gofr/http"
	"strconv"
)

type UserHandler struct {
	Service UserServiceInterface
}

// NewUserHandler : Factory function to implement and return behaviour
func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) Create(c *gofr.Context) (any, error) {

	var u user.User

	if err := c.Bind(&u); err != nil {
		return nil, gofrHttp.ErrorInvalidParam{Params: []string{"User"}}
	}

	user1, err := h.Service.Create(c, u)
	if err != nil {
		return user.User{}, err
	}

	return user1, nil
}

func (h *UserHandler) Get(c *gofr.Context) (any, error) {

	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return user.User{}, gofrHttp.ErrorInvalidParam{Params: []string{"UserID"}}
	}

	user1, err := h.Service.Get(c, id)
	if err != nil {
		return user.User{}, err
	}

	return user1, nil

}

func (h *UserHandler) Delete(c *gofr.Context) (any, error) {
	id, err := strconv.Atoi(c.PathParam("id"))
	if err != nil {
		return user.User{}, gofrHttp.ErrorInvalidParam{Params: []string{"UserID"}}
	}
	err = h.Service.Delete(c, id)
	if err != nil {
		return user.User{}, err
	}

	return fmt.Sprintf("Successfully Deleted user with id %d", id), nil

}

func (h *UserHandler) All(c *gofr.Context) (any, error) {

	users, err := h.Service.All(c)
	if err != nil {
		return nil, err
	}

	return users, nil

}
