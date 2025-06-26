package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/clin211/go-cicd-github-actions-docker/internal/model"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrInvalidUserID = errors.New("invalid user ID")
)

// UserRepository 定义用户数据访问的操作
type UserRepository interface {
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	List() ([]model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
}

// InMemoryUserRepository 使用内存存储实现UserRepository
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[uint]*model.User
	nextID uint
}

// NewInMemoryUserRepository 创建一个新的InMemoryUserRepository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[uint]*model.User),
		nextID: 1,
	}
}

// GetByID 通过ID检索用户
func (r *InMemoryUserRepository) GetByID(id uint) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetByUsername 通过用户名检索用户
func (r *InMemoryUserRepository) GetByUsername(username string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// GetByEmail 通过邮箱检索用户
func (r *InMemoryUserRepository) GetByEmail(email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// List 返回所有用户
func (r *InMemoryUserRepository) List() ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, *user)
	}
	return users, nil
}

// Create 添加一个新用户
func (r *InMemoryUserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 检查用户名是否已存在
	for _, existingUser := range r.users {
		if existingUser.Username == user.Username {
			return ErrUserExists
		}
		if existingUser.Email == user.Email {
			return ErrUserExists
		}
	}

	// 设置ID和时间戳
	user.ID = r.nextID
	r.nextID++
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// 存储用户
	r.users[user.ID] = user
	return nil
}

// Update 修改现有用户
func (r *InMemoryUserRepository) Update(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}

	// 检查用户名或邮箱是否与现有用户冲突
	for id, existingUser := range r.users {
		if id == user.ID {
			continue
		}
		if existingUser.Username == user.Username {
			return ErrUserExists
		}
		if existingUser.Email == user.Email {
			return ErrUserExists
		}
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// Delete 删除用户
func (r *InMemoryUserRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}
