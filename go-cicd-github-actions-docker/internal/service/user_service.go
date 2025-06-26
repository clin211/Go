package service

import (
	"errors"

	"github.com/clin211/go-cicd-github-actions-docker/internal/model"
	"github.com/clin211/go-cicd-github-actions-docker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService 处理用户相关的业务逻辑
type UserService struct {
	repo repository.UserRepository
}

// NewUserService 创建一个新的UserService
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// GetByID 通过ID检索用户
func (s *UserService) GetByID(id uint) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

// GetByUsername 通过用户名检索用户
func (s *UserService) GetByUsername(username string) (*model.UserResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return toUserResponse(user), nil
}

// GetFullUserByUsername 通过用户名检索完整用户信息（包括密码）
func (s *UserService) GetFullUserByUsername(username string) (*model.User, error) {
	return s.repo.GetByUsername(username)
}

// List 检索所有用户
func (s *UserService) List() ([]model.UserResponse, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *toUserResponse(&user)
	}
	return userResponses, nil
}

// Create 创建一个新用户
func (s *UserService) Create(createData model.UserCreate) (*model.UserResponse, error) {
	// 检查用户是否已存在
	_, err := s.repo.GetByUsername(createData.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	_, err = s.repo.GetByEmail(createData.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// 密码哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username:  createData.Username,
		Email:     createData.Email,
		Password:  string(hashedPassword),
		FirstName: createData.FirstName,
		LastName:  createData.LastName,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// Update 更新现有用户
func (s *UserService) Update(id uint, updateData model.UserUpdate) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 如果提供了数据，则更新字段
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}

	err = s.repo.Update(user)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// toUserResponse 将User模型转换为UserResponse
func toUserResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
