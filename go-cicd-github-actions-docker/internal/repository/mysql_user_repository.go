package repository

import (
	"errors"
	"time"

	"github.com/clin211/go-cicd-github-actions-docker/internal/model"
	"github.com/clin211/go-cicd-github-actions-docker/internal/pkg/db"
	"gorm.io/gorm"
)

// MySQLUserRepository 使用MySQL数据库实现UserRepository接口
type MySQLUserRepository struct{}

// NewMySQLUserRepository 创建一个新的基于MySQL的用户仓库
func NewMySQLUserRepository() *MySQLUserRepository {
	return &MySQLUserRepository{}
}

// GetByID 通过ID检索用户
func (r *MySQLUserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User

	if err := db.GetMasterDB().First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername 通过用户名检索用户
func (r *MySQLUserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User

	if err := db.GetMasterDB().Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail 通过邮箱检索用户
func (r *MySQLUserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User

	if err := db.GetMasterDB().Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// List 返回所有用户
func (r *MySQLUserRepository) List() ([]model.User, error) {
	var users []model.User

	if err := db.GetMasterDB().Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// Create 添加一个新用户
func (r *MySQLUserRepository) Create(user *model.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := db.GetMasterDB().Create(user).Error; err != nil {
		return err
	}

	return nil
}

// Update 修改现有用户
func (r *MySQLUserRepository) Update(user *model.User) error {
	user.UpdatedAt = time.Now()

	if err := db.GetMasterDB().Save(user).Error; err != nil {
		return err
	}

	return nil
}

// Delete 删除用户
func (r *MySQLUserRepository) Delete(id uint) error {
	if err := db.GetMasterDB().Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	return nil
}
