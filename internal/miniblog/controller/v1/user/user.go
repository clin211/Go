package user

import (
	"github.com/Forest-211/miniblog/internal/miniblog/biz"
	"github.com/Forest-211/miniblog/internal/miniblog/store"
)

type UserController struct {
	b biz.IBiz
}

func NewUserController(ds store.IStore) *UserController {
	return &UserController{b: biz.NewBiz(ds)}
}
