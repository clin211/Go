package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/Forest-211/miniblog/internal/pkg/log"
	"github.com/Forest-211/miniblog/internal/pkg/model"
)

type PostStore interface {
	Create(ctx context.Context, user *model.PostM) error
	Get(ctx context.Context, username string) (*model.PostM, error)
	Update(ctx context.Context, user *model.PostM) error
}

// UserStore 接口的实现.
type posts struct {
	db *gorm.DB
}

// 确保 users 实现了 UserStore 接口.
var _ UserStore = (*users)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db}
}

func (p *posts) Create(ctx context.Context, post *model.PostM) error {
	return p.db.Create(&post).Error
}

func (p *posts) Get(ctx context.Context, username string) (*model.PostM, error) {
	var post model.PostM
	if err := p.db.Where("username = ?", username).First(&post).Error; err != nil {
		log.C(ctx).Errorw("get post <"+username+"> error", "error", err)
		return nil, err
	}
	return &post, nil
}

func (p *posts) Update(ctx context.Context, post *model.PostM) error {
	return p.db.Save(post).Error
}
