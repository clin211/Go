package post

import (
	"context"
	"fmt"

	"github.com/Forest-211/miniblog/internal/miniblog/store"
	"github.com/Forest-211/miniblog/internal/pkg/known"
	"github.com/Forest-211/miniblog/internal/pkg/model"
	v1 "github.com/Forest-211/miniblog/pkg/api/miniblog/v1"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type PostBiz interface {
	Create(ctx context.Context, r *v1.CreatePostRequest) error
}

type postBiz struct {
	ds store.IStore
}

var _ PostBiz = (*postBiz)(nil)

// New 创建一个实现了 PostBiz 接口的实例.
func New(ds store.IStore) PostBiz {
	return &postBiz{ds: ds}
}

func (p *postBiz) Create(ctx context.Context, r *v1.CreatePostRequest) error {
	var postM model.PostM
	_ = copier.Copy(&postM, r)
	postM.PostID = uuid.New().String()
	postM.Username = ctx.Value(known.XUsernameKey).(string)
	fmt.Println("postM: ", postM)
	if err := p.ds.Posts().Create(ctx, &postM); err != nil {
		return err
	}
	return nil
}
