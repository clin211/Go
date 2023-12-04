package user

import (
	"github.com/Forest-211/miniblog/internal/miniblog/biz"
	"github.com/Forest-211/miniblog/internal/miniblog/store"
	"github.com/Forest-211/miniblog/pkg/auth"
)

type UserController struct {
	a *auth.Authz
	b biz.IBiz
}

func New(ds store.IStore, a *auth.Authz) *UserController {
	return &UserController{a: a, b: biz.NewBiz(ds)}
}
