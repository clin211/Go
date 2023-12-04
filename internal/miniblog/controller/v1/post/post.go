package post

import (
	"github.com/Forest-211/miniblog/internal/miniblog/biz"
	"github.com/Forest-211/miniblog/internal/miniblog/store"
	"github.com/Forest-211/miniblog/pkg/auth"
)

type PostController struct {
	a *auth.Authz
	b biz.IBiz
}

func New(ds store.IStore, a *auth.Authz) *PostController {
	return &PostController{a: a, b: biz.NewBiz(ds)}
}
